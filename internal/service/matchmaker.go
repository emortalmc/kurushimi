package service

import (
	"context"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	pbparty "github.com/emortalmc/proto-specs/gen/go/grpc/party"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kurushimi/internal/kafka"
	"kurushimi/internal/repository"
	"kurushimi/internal/repository/model"
	"kurushimi/pkg/pb"
)

type matchmakerService struct {
	pb.MatchmakerServer

	repo          repository.Repository
	notifier      kafka.Notifier
	cfgController liveconfig.GameModeConfigController

	partyService         pbparty.PartyServiceClient
	partySettingsService pbparty.PartySettingsServiceClient
}

func NewMatchmakerService(repository repository.Repository, notifier kafka.Notifier,
	cfgController liveconfig.GameModeConfigController, partyService pbparty.PartyServiceClient,
	partySettingsService pbparty.PartySettingsServiceClient) pb.MatchmakerServer {

	return &matchmakerService{
		repo:          repository,
		notifier:      notifier,
		cfgController: cfgController,

		partyService:         partyService,
		partySettingsService: partySettingsService,
	}
}

// QueueByPlayer requests a player is queued for a game.
// NOTE: A player is always in a party and we clean up when a player changes party.
// Therefore, we only need to check if the player's party is in a queue, not the player themselves.
// TODO go over and put in custom error responses
// TODO not tested
func (m *matchmakerService) QueueByPlayer(ctx context.Context, request *pb.QueueByPlayerRequest) (*pb.QueueByPlayerResponse, error) {
	playerId, err := uuid.Parse(request.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player_id")
	}

	// check if game_mode is valid
	modeConfig := m.cfgController.GetCurrentConfig(request.GameModeId)
	if modeConfig == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid game_mode_id")
	}

	// check if map is present
	if request.MapId != nil {
		_, ok := modeConfig.Maps[*request.MapId]
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid map_id")
		}
		// todo check if player has enough credits
	}

	res, err := m.partyService.GetParty(ctx, &pbparty.GetPartyRequest{Id: &pbparty.GetPartyRequest_PlayerId{PlayerId: playerId.String()}})
	if err != nil {
		return nil, err
	}

	party := res.Party
	partyId, err := primitive.ObjectIDFromHex(party.Id)
	if err != nil {
		return nil, err
	}

	// check if the party is already in a queue
	// if yes, return error
	isInQueue, err := m.repo.IsPartyQueued(ctx, partyId)
	if err != nil {
		return nil, err
	}

	if isInQueue {
		return nil, status.Error(codes.AlreadyExists, "party is already in queue")
	}

	// check if player is leader of party
	if playerId.String() != party.GetLeaderId() {
		return nil, status.Error(codes.PermissionDenied, "player is not leader of party")
	}

	// private game checks
	var privateGame bool
	if request.PrivateGame != nil {
		privateGame = *request.PrivateGame
	} else {
		privateGame = false
	}

	if privateGame {
		// check if there are enough players
		if len(party.Members) < modeConfig.MinPlayers {
			return nil, status.Error(codes.InvalidArgument, "not enough players to create private game")
		}

		// TODO check if player has enough credits
	}

	// Get the party settings
	settingsRes, err := m.partySettingsService.GetPartySettings(ctx, &pbparty.GetPartySettingsRequest{Id: &pbparty.GetPartySettingsRequest_PlayerId{PlayerId: party.LeaderId}})
	if err != nil {
		return nil, err
	}

	settings := settingsRes.GetSettings()

	autoTeleport := true
	if request.AutoTeleport != nil {
		autoTeleport = *request.AutoTeleport
	}

	memberIds := make([]uuid.UUID, 0)
	for _, member := range party.Members {
		id, err := uuid.Parse(member.Id)
		if err != nil {
			return nil, err
		}

		memberIds = append(memberIds, id)
	}

	partyLeaderId, err := uuid.Parse(party.LeaderId)
	if err != nil {
		return nil, err
	}

	// create ticket
	ticket := model.NewTicket(partyId, &model.ReducedPartySettings{
		LeaderId:            partyLeaderId,
		DequeueOnDisconnect: settings.DequeueOnDisconnect,
		AllowMemberDequeue:  settings.AllowMemberDequeue,
	}, memberIds, request.GameModeId, autoTeleport)

	err = m.repo.ExecuteTransaction(ctx, func(ctx mongo.SessionContext) error {
		err = m.repo.CreateTicket(ctx, ticket)
		if err != nil {
			return err
		}

		// save player map choice
		queuedPlayers := make([]*model.QueuedPlayer, 0)
		for _, playerId := range memberIds {
			var mapId *string
			if request.MapId != nil && playerId.String() == party.LeaderId {
				mapId = request.MapId
			}

			queuedPlayers = append(queuedPlayers, &model.QueuedPlayer{
				PlayerId: playerId,
				TicketId: ticket.Id,
				MapId:    mapId,
			})
		}

		err = m.repo.CreateQueuedPlayers(ctx, queuedPlayers)
		if err != nil {
			return err
		}

		err = m.notifier.TicketCreated(ctx, ticket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.QueueByPlayerResponse{}, nil
}

var (
	dequeueNotInQueueErr = panicIfErr(status.New(codes.NotFound, "player is not in queue").
		WithDetails(&pb.DequeueByPlayerErrorResponse{Reason: pb.DequeueByPlayerErrorResponse_NOT_IN_QUEUE})).Err()

	dequeueNoPermissionErr = panicIfErr(status.New(codes.PermissionDenied, "player does not have permission to dequeue").
		WithDetails(&pb.DequeueByPlayerErrorResponse{Reason: pb.DequeueByPlayerErrorResponse_NO_PERMISSION})).Err()

	dequeueAlreadyDequeuedErr = panicIfErr(status.New(codes.AlreadyExists, "party is already requested for dequeue").
		WithDetails(&pb.DequeueByPlayerErrorResponse{Reason: pb.DequeueByPlayerErrorResponse_ALREADY_MARKED_FOR_DEQUEUE})).Err()
)

// DequeueByPlayer requests a player is dequeued from a game. Note that this player acts on behalf of the party, not themselves.
// NOTE: This DOES NOT actually dequeue the player, it just requests it.
// Not every error will be fed back as it happens lazily.
func (m *matchmakerService) DequeueByPlayer(ctx context.Context, request *pb.DequeueByPlayerRequest) (*pb.DequeueByPlayerResponse, error) {
	playerId, err := uuid.Parse(request.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player_id")
	}

	ticket, err := m.repo.GetTicketByPlayerId(ctx, playerId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, dequeueNotInQueueErr
		}

		return nil, err
	}

	// check if the player has permission to dequeue
	if playerId != ticket.PartySettings.LeaderId && !ticket.PartySettings.AllowMemberDequeue {
		return nil, dequeueNoPermissionErr
	}

	// TODO we could do a preemptive check to predict if a PendingMatch is being processed and will become a Match
	// if so, we could return an error here to prevent the player from being dequeued

	modCount, err := m.repo.AddTicketDequeueRequest(ctx, ticket.Id)
	if err != nil {
		return nil, err
	}

	if modCount == 0 {
		return nil, dequeueAlreadyDequeuedErr
	}

	return &pb.DequeueByPlayerResponse{}, nil
}

var (
	changeMapInvalidMapErr = panicIfErr(status.New(codes.InvalidArgument, "invalid map_id").
		WithDetails(&pb.ChangePlayerMapVoteErrorResponse{Reason: pb.ChangePlayerMapVoteErrorResponse_INVALID_MAP})).Err()

	changeMapNotInQueueErr = panicIfErr(status.New(codes.NotFound, "player is not in queue").
		WithDetails(&pb.ChangePlayerMapVoteErrorResponse{Reason: pb.ChangePlayerMapVoteErrorResponse_NOT_IN_QUEUE})).Err()
)

func (m *matchmakerService) ChangePlayerMapVote(ctx context.Context, request *pb.ChangePlayerMapVoteRequest) (*pb.ChangePlayerMapVoteResponse, error) {
	playerId, err := uuid.Parse(request.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player_id")
	}

	ok := m.isMapIdValid(request.MapId)
	if !ok {
		return nil, changeMapInvalidMapErr
	}

	err = m.repo.SetMapIdOfQueuedPlayer(ctx, playerId, request.MapId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, changeMapNotInQueueErr
		}

		return nil, err
	}

	return &pb.ChangePlayerMapVoteResponse{}, nil
}

func (m *matchmakerService) isMapIdValid(mapId string) bool {
	cfg := m.cfgController.GetCurrentConfig(mapId)
	return cfg != nil
}

func panicIfErr[T any](thing T, err error) T {
	if err != nil {
		panic(err)
	}
	return thing
}
