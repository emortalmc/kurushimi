package lobbycontroller

import (
	v13 "agones.dev/agones/pkg/apis/allocation/v1"
	v1 "agones.dev/agones/pkg/client/clientset/versioned/typed/allocation/v1"
	"context"
	"github.com/emortalmc/kurushimi/internal/config"
	"github.com/emortalmc/kurushimi/internal/gsallocation"
	"github.com/emortalmc/kurushimi/internal/gsallocation/selector"
	"github.com/emortalmc/kurushimi/internal/kafka"
	pb "github.com/emortalmc/proto-specs/gen/go/model/matchmaker"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type LobbyController interface {
	QueuePlayer(playerId uuid.UUID, autoTeleport bool)
}

type lobbyControllerImpl struct {
	fleetName       string
	matchmakingRate time.Duration
	playersPerMatch int

	logger          *zap.SugaredLogger
	notifier        kafka.Notifier
	allocatorClient v1.GameServerAllocationInterface

	// queuedPlayers map[playerId]autoTeleport
	queuedPlayers     map[uuid.UUID]bool
	queuedPlayersLock sync.Mutex
}

func NewLobbyController(ctx context.Context, wg *sync.WaitGroup, logger *zap.SugaredLogger, cfg *config.Config, notifier kafka.Notifier,
	allocatorClient v1.GameServerAllocationInterface) LobbyController {

	c := &lobbyControllerImpl{
		fleetName:       cfg.LobbyFleetName,
		matchmakingRate: cfg.LobbyMatchRate,
		playersPerMatch: cfg.LobbyMatchSize,

		logger:          logger,
		notifier:        notifier,
		allocatorClient: allocatorClient,

		queuedPlayers:     make(map[uuid.UUID]bool),
		queuedPlayersLock: sync.Mutex{},
	}

	c.run(wg, ctx)

	return c
}

func (l *lobbyControllerImpl) QueuePlayer(playerId uuid.UUID, autoTeleport bool) {
	l.queuedPlayersLock.Lock()
	defer l.queuedPlayersLock.Unlock()

	l.queuedPlayers[playerId] = autoTeleport
}

func (l *lobbyControllerImpl) run(wg *sync.WaitGroup, ctx context.Context) {
	go func() {
		for {
			if ctx.Err() != nil {
				wg.Done()
				return
			}

			lastRunTime := time.Now()

			queuedPlayers := l.resetQueuedPlayers()

			matchAllocationReqMap := l.createMatchesFromPlayers(queuedPlayers)

			allocationErrors := gsallocation.AllocateServers(ctx, l.allocatorClient, matchAllocationReqMap)
			for match, err := range allocationErrors {
				l.logger.Errorw("failed to allocate server for match", "error", err, "match", match)
			}

			l.logger.Infow("created matches", "matchCount", len(matchAllocationReqMap))
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
		}
	}()
}

func (l *lobbyControllerImpl) resetQueuedPlayers() map[uuid.UUID]bool {
	l.queuedPlayersLock.Lock()
	defer l.queuedPlayersLock.Unlock()

	queuedPlayers := make(map[uuid.UUID]bool, len(l.queuedPlayers))
	for playerId, autoTeleport := range l.queuedPlayers {
		queuedPlayers[playerId] = autoTeleport
	}

	l.queuedPlayers = make(map[uuid.UUID]bool)

	return queuedPlayers
}

func (l *lobbyControllerImpl) createMatchesFromPlayers(playerMap map[uuid.UUID]bool) map[*pb.Match]*v13.GameServerAllocation {
	allocationReqs := make(map[*pb.Match]*v13.GameServerAllocation)

	l.logger.Infow("creating matches from players", "playerCount", len(playerMap))

	currentMatch := &pb.Match{
		Id:         primitive.NewObjectID().String(),
		GameModeId: "lobby",
		MapId:      nil,
		Tickets:    make([]*pb.Ticket, 0),
		Assignment: nil,
	}

	currentCount := 0
	for playerId, autoTeleport := range playerMap {
		currentMatch.Tickets = append(currentMatch.Tickets, &pb.Ticket{
			PlayerIds:           []string{playerId.String()},
			CreatedAt:           timestamppb.Now(),
			GameModeId:          "lobby",
			AutoTeleport:        autoTeleport,
			DequeueOnDisconnect: false,
			InPendingMatch:      false,
		})
		currentCount++

		if currentCount >= l.playersPerMatch {
			allocationReqs[currentMatch] = selector.CreatePlayerBasedSelector(l.fleetName, currentMatch, int64(currentCount))
			currentMatch = &pb.Match{
				Id:         primitive.NewObjectID().String(),
				GameModeId: "lobby",
				MapId:      nil,
				Tickets:    make([]*pb.Ticket, 0),
				Assignment: nil,
			}
			currentCount = 0
		}
	}

	if currentCount > 0 {
		allocationReqs[currentMatch] = selector.CreatePlayerBasedSelector(l.fleetName, currentMatch, int64(len(currentMatch.Tickets)))
	}

	l.logger.Infow("created matches from players", "matchCount", len(allocationReqs), "playerCount", len(playerMap))

	return allocationReqs
}
