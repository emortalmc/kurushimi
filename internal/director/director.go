package director

import (
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/internal/config"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/matchfunction"
	"kurushimi/internal/notifier"
	"kurushimi/internal/statestore"
	"kurushimi/internal/utils/kubernetes"
	"kurushimi/pkg/pb"
	"log"
	"math/rand"
	"time"
)

type director interface {
	run(ctx context.Context, p profile.ModeProfile)
	removeTickets(ctx context.Context, tickets []*pb.Ticket)
	handlePendingMatches(ctx context.Context, pendingMatches []*pb.PendingMatch) []*pb.Match
	assign(ctx context.Context, p profile.ModeProfile, matches []*pb.Match)
}

type directorImpl struct {
	director

	logger     *zap.SugaredLogger
	notifier   notifier.Notifier
	stateStore statestore.StateStore

	namespace string
}

func Init(ctx context.Context, logger *zap.SugaredLogger, notifier notifier.Notifier, stateStore statestore.StateStore, namespace string) {
	modeProfiles := config.ModeProfiles

	d := directorImpl{
		logger:     logger,
		notifier:   notifier,
		stateStore: stateStore,

		namespace: namespace,
	}

	for _, p := range modeProfiles {
		go func(p profile.ModeProfile) {
			// Only run every x milliseconds, but if that has already passed, run immediately.
			for {
				lastRunTime := time.Now()
				d.run(ctx, p)
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
func (d *directorImpl) run(ctx context.Context, p profile.ModeProfile) {
	matches, pendingMatches, err := matchfunction.Run(ctx, d.notifier, d.stateStore, p)
	if err != nil {
		d.logger.Errorw("Failed to fetch matches", "profileName", p.Name, err)
		return
	}

	// convert and pending matches to matches if they are ready to teleport
	if len(pendingMatches) > 0 {
		convMatches := d.handlePendingMatches(ctx, pendingMatches)
		matches = append(matches, convMatches...)
	}
	if len(matches) > 0 {
		log.Printf("Found %d matches: %+v", len(matches), matches)
	}

	// NOTE: We don't add pendingMatch tickets here as they are depended upon
	// to detect if someone leaves the queue before teleporting.
	finishedTickets := make([]*pb.Ticket, 0)
	for _, match := range matches {
		finishedTickets = append(finishedTickets, match.Tickets...)
	}
	d.removeTickets(ctx, finishedTickets)

	d.assign(ctx, p, matches)
}

func (d *directorImpl) removeTickets(ctx context.Context, tickets []*pb.Ticket) {
	for _, ticket := range tickets {
		err := d.stateStore.DeleteTicket(ctx, ticket.Id)
		if err != nil {
			d.logger.Error("Failed to delete ticket", zap.Error(err))
			continue
		}
		err = d.stateStore.UnIndexTicket(ctx, ticket.Id)
		if err != nil {
			d.logger.Error("Failed to unindex ticket", zap.Error(err))
			continue
		}
	}
}

func (d *directorImpl) handlePendingMatches(ctx context.Context, pendingMatches []*pb.PendingMatch) []*pb.Match {
	matches := make([]*pb.Match, 0)

	for _, pendingMatch := range pendingMatches {
		if pendingMatch.TeleportTime.AsTime().Before(time.Now()) {
			// ready to teleport
			err := d.stateStore.DeletePendingMatch(ctx, pendingMatch.Id)
			if err != nil {
				d.logger.Error("Failed to delete pending match", zap.Error(err))
				continue
			}
			err = d.stateStore.UnIndexPendingMatch(ctx, pendingMatch)
			if err != nil {
				d.logger.Error("Failed to unindex pending match", zap.Error(err))
				continue
			}

			matches = append(matches, &pb.Match{
				Id:          pendingMatch.Id,
				ProfileName: pendingMatch.ProfileName,
				Tickets:     pendingMatch.Tickets,
			})
		} else {
			// not ready to teleport yet

			err := d.stateStore.CreatePendingMatch(ctx, pendingMatch)
			if err != nil {
				d.logger.Error("Failed to save pending match", zap.Error(err))
				continue
			}
			err = d.stateStore.IndexPendingMatch(ctx, pendingMatch)
			if err != nil {
				d.logger.Error("Failed to index pending match", zap.Error(err))
				continue
			}
		}
	}
	return matches
}

func (d *directorImpl) assign(ctx context.Context, p profile.ModeProfile, matches []*pb.Match) {
	logger := zap.S().With("profileName", p.Name)
	for _, match := range matches {
		if d.namespace == "" {
			match.Assignment = &pb.Assignment{
				ServerId:      "mock-gameserver-xxxx",
				ServerAddress: "0.0.0.0",
				ServerPort:    rand.Uint32(),
			}
		} else {
			alloc, err := kubernetes.AgonesClient.AllocationV1().GameServerAllocations(d.namespace).Create(ctx, p.Selector(p, match), v1.CreateOptions{})
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
		err := d.notifier.NotifyMatchTeleport(ctx, match)
		if err != nil {
			logger.Error("Failed to notify transport", err)
			continue
		}
	}
}
