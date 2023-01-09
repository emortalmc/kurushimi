package main

import (
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/internal/config"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/frontend"
	"kurushimi/internal/matchfunction"
	"kurushimi/internal/notifier"
	"kurushimi/internal/statestore"
	"kurushimi/internal/utils/kubernetes"
	"kurushimi/pkg/pb"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const (
	minTimeBetweenRuns = 4 * time.Second
)

var (
	logger, _ = zap.NewProduction()
	namespace = os.Getenv("NAMESPACE")
)

func main() {
	kubernetes.Init()

	modeProfiles := config.ModeProfiles

	logger.Info("Starting Kurushimi",
		zap.Int("profileCount", len(modeProfiles)),
		zap.Any("profiles", modeProfiles),
	)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go frontend.Run(ctx)

	for _, p := range modeProfiles {
		go func(p profile.ModeProfile) {
			// Only run every x milliseconds, but if that has already passed, run immediately.
			for {
				lastRunTime := time.Now()
				run(ctx, p)
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
	select {
	case <-ctx.Done():
		return
	}
}

// This is blocking so another run can't happen until this one is done.
func run(ctx context.Context, p profile.ModeProfile) {
	matches, pendingMatches, err := matchfunction.Run(ctx, p)
	if err != nil {
		logger.Error("Failed to fetch matches", zap.String("profileName", p.Name), zap.Error(err))
		return
	}

	logger.Debug("Generated matches", zap.Int("generated", len(matches)),
		zap.Int("generatedPending", len(pendingMatches)),
		zap.String("profileName", p.Name),
	)

	usedTickets := make([]*pb.Ticket, 0)
	for _, match := range matches {
		usedTickets = append(usedTickets, match.Tickets...)
	}
	for _, pendingMatch := range pendingMatches {
		usedTickets = append(usedTickets, pendingMatch.Tickets...)
	}
	removeTickets(ctx, usedTickets)

	// convert and pending matches to matches if they are ready to teleport
	convMatches := handlePendingMatches(ctx, pendingMatches)
	matches = append(matches, convMatches...)

	assign(ctx, p, matches)
}

func removeTickets(ctx context.Context, tickets []*pb.Ticket) {
	for _, ticket := range tickets {
		err := statestore.DeleteTicket(ctx, ticket.Id)
		if err != nil {
			logger.Error("Failed to delete ticket", zap.Error(err))
			continue
		}
		err = statestore.UnIndexTicket(ctx, ticket.Id)
		if err != nil {
			logger.Error("Failed to unindex ticket", zap.Error(err))
			continue
		}
	}
}

func handlePendingMatches(ctx context.Context, pendingMatches []*pb.PendingMatch) []*pb.Match {
	matches := make([]*pb.Match, 0)

	for _, pendingMatch := range pendingMatches {
		if pendingMatch.TeleportTime.AsTime().Before(time.Now()) {
			// ready to teleport
			err := statestore.DeletePendingMatch(ctx, pendingMatch.Id)
			if err != nil {
				logger.Error("Failed to delete pending match", zap.Error(err))
				continue
			}
			err = statestore.UnIndexPendingMatch(ctx, pendingMatch)
			if err != nil {
				logger.Error("Failed to unindex pending match", zap.Error(err))
				continue
			}

			matches = append(matches, &pb.Match{
				Id:          pendingMatch.Id,
				ProfileName: pendingMatch.ProfileName,
				Tickets:     pendingMatch.Tickets,
			})
		} else {
			// not ready to teleport yet

			err := statestore.CreatePendingMatch(ctx, pendingMatch)
			if err != nil {
				logger.Error("Failed to save pending match", zap.Error(err))
				continue
			}
			err = statestore.IndexPendingMatch(ctx, pendingMatch)
			if err != nil {
				logger.Error("Failed to index pending match", zap.Error(err))
				continue
			}
		}
	}
	return matches
}

func assign(ctx context.Context, p profile.ModeProfile, matches []*pb.Match) {
	for _, match := range matches {
		if namespace == "" {
			match.Assignment = &pb.Assignment{
				ServerId:      "mock-gameserver-xxxx",
				ServerAddress: "0.0.0.0",
				ServerPort:    rand.Uint32(),
			}
		} else {
			alloc, err := kubernetes.AgonesClient.AllocationV1().GameServerAllocations(namespace).Create(ctx, p.Selector(p, match), v1.CreateOptions{})
			if err != nil {
				logger.Error("Failed to allocate gameserver", zap.Error(err))
				continue
			}

			status := alloc.Status
			if status.State != allocatorv1.GameServerAllocationAllocated {
				logger.Error("Failed to allocate server", zap.String("matchId", match.Id), zap.Any("status", status))
				continue
			}
			match.Assignment = &pb.Assignment{
				ServerId:      status.GameServerName,
				ServerAddress: status.Address,
				ServerPort:    uint32(status.Ports[0].Port),
			}
		}
		logger.Info("Assigned server to match", zap.String("matchId", match.Id), zap.Any("assignment", match.Assignment))
		err := notifier.NotifyMatchTeleport(ctx, match)
		if err != nil {
			logger.Error("Failed to notify transport", zap.Error(err))
			continue
		}
	}
}
