package lobbycontroller

import (
	v13 "agones.dev/agones/pkg/apis/allocation/v1"
	v1 "agones.dev/agones/pkg/client/clientset/versioned/typed/allocation/v1"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/internal/config"
	"kurushimi/internal/gsallocation/selector"
	"kurushimi/internal/kafka"
	"kurushimi/internal/utils/protoutils"
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

			matches := l.createMatchesFromPlayers(queuedPlayers)
			// Allocate server
			// TODO make matches get allocated in parallel
			for _, match := range matches {
				if err := l.allocateServer(ctx, match); err != nil {
					l.logger.Errorw("failed to allocate server", "error", err)
					continue
				}
			}

			for _, match := range matches {
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

func (l *lobbyControllerImpl) createMatchesFromPlayers(playerIds []uuid.UUID) []*pb.Match {
	matches := make([]*pb.Match, 0)

	currentMatch := &pb.Match{
		Id:         primitive.NewObjectID().String(),
		GameModeId: "lobby",
		MapId:      nil,
		Tickets:    make([]*pb.Ticket, 0),
		Assignment: nil,
	}
	for _, playerId := range playerIds {
		currentMatch.Tickets = append(currentMatch.Tickets, &pb.Ticket{
			PlayerIds: []string{playerId.String()},
		})

		if len(currentMatch.Tickets) >= l.playersPerMatch {
			matches = append(matches, currentMatch)
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
		matches = append(matches, currentMatch)
	}

	return matches
}

// TODO can we PLEASE make this common code?
// allocateServer allocates a server for the given match.
// Sets the assigned server in the pb.Match provided.
func (l *lobbyControllerImpl) allocateServer(ctx context.Context, match *pb.Match) error {
	playerCount := protoutils.GetMatchPlayerCount(match)

	resp, err := l.allocatorClient.Create(ctx,
		selector.CreatePlayerBasedSelector(l.fleetName, match, playerCount),
		v12.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to allocate server: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("allocation response was nil")
	}

	allocation := resp.Status
	if allocation.State != v13.GameServerAllocationAllocated {
		return fmt.Errorf("allocation state was not allocated: %s", allocation.State)
	}

	match.Assignment = &pb.Assignment{
		ServerId:      allocation.GameServerName,
		ServerAddress: allocation.Address,
		ServerPort:    uint32(allocation.Ports[0].Port),
	}

	return nil
}
