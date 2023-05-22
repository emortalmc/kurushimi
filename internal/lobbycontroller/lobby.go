package lobbycontroller

import (
	v13 "agones.dev/agones/pkg/apis/allocation/v1"
	v1 "agones.dev/agones/pkg/client/clientset/versioned/typed/allocation/v1"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/internal/config"
	"kurushimi/internal/gsallocation"
	"kurushimi/internal/gsallocation/selector"
	"kurushimi/internal/kafka"
	"kurushimi/pkg/pb"
	"sync"
	"time"
)

type LobbyController interface {
	QueuePlayer(playerId uuid.UUID)

	// Run runs the lobby controller.
	// NOTE: This is blocking.
	Run(ctx context.Context)
}

type lobbyControllerImpl struct {
	fleetName       string
	matchmakingRate time.Duration
	playersPerMatch int

	logger          *zap.SugaredLogger
	notifier        kafka.Notifier
	allocatorClient v1.GameServerAllocationInterface

	queuedPlayers     []uuid.UUID
	queuedPlayersLock sync.Mutex
}

func NewLobbyController(logger *zap.SugaredLogger, cfg *config.Config, notifier kafka.Notifier,
	allocatorClient v1.GameServerAllocationInterface) LobbyController {

	return &lobbyControllerImpl{
		fleetName:       cfg.LobbyFleetName,
		matchmakingRate: cfg.LobbyMatchRate,
		playersPerMatch: cfg.LobbyMatchSize,

		logger:          logger,
		notifier:        notifier,
		allocatorClient: allocatorClient,

		queuedPlayers:     make([]uuid.UUID, 0),
		queuedPlayersLock: sync.Mutex{},
	}

}

func (l *lobbyControllerImpl) QueuePlayer(playerId uuid.UUID) {
	l.queuedPlayersLock.Lock()
	defer l.queuedPlayersLock.Unlock()

	l.queuedPlayers = append(l.queuedPlayers, playerId)
}

func (l *lobbyControllerImpl) Run(ctx context.Context) {
	for {
		func() {
			lastRunTime := time.Now()

			queuedPlayers := l.resetQueuedPlayers()
			queuedPlayers = removeSliceDuplicates(queuedPlayers)

			matchAllocationReqMap := l.createMatchesFromPlayers(queuedPlayers)

			gsallocation.AllocateServers(ctx, l.allocatorClient, matchAllocationReqMap)

			for match := range matchAllocationReqMap {
				if err := l.notifier.MatchCreated(ctx, match); err != nil {
					l.logger.Errorw("failed to send match created message", "error", err)
				}
			}

			// Wait for the next run
			timeSinceLastRun := time.Since(lastRunTime)
			if timeSinceLastRun < l.matchmakingRate {
				time.Sleep(l.matchmakingRate - timeSinceLastRun)
			}
		}()
	}
}

func (l *lobbyControllerImpl) resetQueuedPlayers() []uuid.UUID {
	l.queuedPlayersLock.Lock()
	defer l.queuedPlayersLock.Unlock()

	queuedPlayers := make([]uuid.UUID, len(l.queuedPlayers))
	copy(queuedPlayers, l.queuedPlayers)

	l.queuedPlayers = make([]uuid.UUID, 0)

	return queuedPlayers
}

func (l *lobbyControllerImpl) createMatchesFromPlayers(playerIds []uuid.UUID) map[*pb.Match]*v13.GameServerAllocation {
	allocationReqs := make(map[*pb.Match]*v13.GameServerAllocation)

	currentMatch := &pb.Match{
		Id:         primitive.NewObjectID().String(),
		GameModeId: "lobby",
		MapId:      nil,
		Tickets:    make([]*pb.Ticket, 0),
		Assignment: nil,
	}
	for _, playerId := range playerIds {
		currentMatch.Tickets = append(currentMatch.Tickets, &pb.Ticket{
			PlayerIds:           []string{playerId.String()},
			CreatedAt:           timestamppb.Now(),
			GameModeId:          "lobby",
			AutoTeleport:        true,
			DequeueOnDisconnect: false,
			InPendingMatch:      false,
		})

		if len(currentMatch.Tickets) >= l.playersPerMatch {
			allocationReqs[currentMatch] = selector.CreatePlayerBasedSelector(l.fleetName, currentMatch, int64(len(currentMatch.Tickets)))
			currentMatch = &pb.Match{
				Id:         primitive.NewObjectID().String(),
				GameModeId: "lobby",
				MapId:      nil,
				Tickets:    make([]*pb.Ticket, 0),
				Assignment: nil,
			}
		}
	}

	if len(currentMatch.Tickets) > 0 {
		allocationReqs[currentMatch] = selector.CreatePlayerBasedSelector(l.fleetName, currentMatch, int64(len(currentMatch.Tickets)))
	}

	return allocationReqs
}

func removeSliceDuplicates(strSlice []uuid.UUID) []uuid.UUID {
	allKeys := make(map[uuid.UUID]bool)
	var list []uuid.UUID
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
