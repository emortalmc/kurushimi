package director

import (
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/internal/config"
	"kurushimi/internal/config/dynamic"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/matchfunction"
	"kurushimi/internal/messaging"
	"kurushimi/internal/notifier"
	"kurushimi/internal/statestore"
	"kurushimi/internal/utils/kubernetes"
	"kurushimi/pkg/pb"
	"math/rand"
	"time"
)

type KurushimiApplication struct {
	StateStore statestore.StateStore
	Config     dynamic.Config
	Messaging  messaging.Messenger
	Logger     *zap.SugaredLogger
}

func Init(ctx context.Context, n *notifier.Notifier, app KurushimiApplication) {
	modeProfiles := config.ModeProfiles

	for _, p := range modeProfiles {
		go func(p profile.ModeProfile) {
			// Only run every x milliseconds, but if that has already passed, run immediately.
			for {
				lastRunTime := time.Now()
				app.run(ctx, n, p)
				timeSinceLastRun := time.Since(lastRunTime)
				if timeSinceLastRun < p.MatchmakingRate {
					select {
					case <-ctx.Done():
						return
					case <-time.After(p.MatchmakingRate - timeSinceLastRun):

					}
				}
			}
		}(p)
	}
}

// This is blocking so another run can't happen until this one is done.
func (app *KurushimiApplication) run(ctx context.Context, n *notifier.Notifier, p profile.ModeProfile) {
	matches, pendingMatches, err := matchfunction.Run(ctx, app.StateStore, p)
	if err != nil {
		app.Logger.Errorw("Failed to fetch matches", "profileName", p.Name, err)
		return
	}

	// convert and pending matches to matches if they are ready to teleport
	convMatches := app.handlePendingMatches(ctx, pendingMatches)
	matches = append(matches, convMatches...)

	// NOTE: We don't add pendingMatch tickets here as they are depended upon
	// to detect if someone leaves the queue before teleporting.
	finishedTickets := make([]*pb.Ticket, 0)
	for _, match := range matches {
		finishedTickets = append(finishedTickets, match.Tickets...)
	}
	app.removeTickets(ctx, finishedTickets)

	app.assign(ctx, n, p, matches)
}

func (app *KurushimiApplication) removeTickets(ctx context.Context, tickets []*pb.Ticket) {
	for _, ticket := range tickets {
		err := app.StateStore.DeleteTicket(ctx, ticket.Id)
		if err != nil {
			app.Logger.Error("Failed to delete ticket", zap.Error(err))
			continue
		}
		err = app.StateStore.UnIndexTicket(ctx, ticket.Id)
		if err != nil {
			app.Logger.Error("Failed to unindex ticket", zap.Error(err))
			continue
		}
	}
}

func (app *KurushimiApplication) handlePendingMatches(ctx context.Context, pendingMatches []*pb.PendingMatch) []*pb.Match {
	matches := make([]*pb.Match, 0)

	for _, pendingMatch := range pendingMatches {
		if pendingMatch.TeleportTime.AsTime().Before(time.Now()) {
			// ready to teleport
			err := app.StateStore.DeletePendingMatch(ctx, pendingMatch.Id)
			if err != nil {
				app.Logger.Error("Failed to delete pending match", zap.Error(err))
				continue
			}
			err = app.StateStore.UnIndexPendingMatch(ctx, pendingMatch)
			if err != nil {
				app.Logger.Error("Failed to unindex pending match", zap.Error(err))
				continue
			}

			matches = append(matches, &pb.Match{
				Id:          pendingMatch.Id,
				ProfileName: pendingMatch.ProfileName,
				Tickets:     pendingMatch.Tickets,
			})
		} else {
			// not ready to teleport yet

			err := app.StateStore.CreatePendingMatch(ctx, pendingMatch)
			if err != nil {
				app.Logger.Error("Failed to save pending match", zap.Error(err))
				continue
			}
			err = app.StateStore.IndexPendingMatch(ctx, pendingMatch)
			if err != nil {
				app.Logger.Error("Failed to index pending match", zap.Error(err))
				continue
			}
		}
	}
	return matches
}

func (app *KurushimiApplication) assign(ctx context.Context, n *notifier.Notifier, p profile.ModeProfile, matches []*pb.Match) {
	logger := zap.S().With("profileName", p.Name)
	for _, match := range matches {
		if app.Config.Namespace == "" {
			match.Assignment = &pb.Assignment{
				ServerId:      "mock-gameserver-xxxx",
				ServerAddress: "0.0.0.0",
				ServerPort:    rand.Uint32(),
			}
		} else {
			alloc, err := kubernetes.AgonesClient.AllocationV1().GameServerAllocations(app.Config.Namespace).Create(ctx, p.Selector(p, match), v1.CreateOptions{})
			if err != nil {
				logger.Errorw("Failed to allocate gameserver", err)
				continue
			}

			status := alloc.Status
			if status.State != allocatorv1.GameServerAllocationAllocated {
				logger.Errorw("Failed to allocate server", "matchId", match.Id, "status", status)
				continue
			}
			match.Assignment = &pb.Assignment{
				ServerId:      status.GameServerName,
				ServerAddress: status.Address,
				ServerPort:    uint32(status.Ports[0].Port),
			}
		}
		logger.Infow("Assigned server to match", "matchId", match.Id, "assignment", match.Assignment)
		err := n.NotifyMatchTeleport(ctx, match)
		if err != nil {
			logger.Error("Failed to notify transport", err)
			continue
		}
	}
}
