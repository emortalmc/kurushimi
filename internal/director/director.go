package main

import (
	"context"
	"go.uber.org/zap"
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

	// Only run every x milliseconds, but if that has already passed, run immediately.
	for {
		lastRunTime := time.Now()
		run(ctx, modeProfiles)
		timeSinceLastRun := time.Since(lastRunTime)
		if timeSinceLastRun < minTimeBetweenRuns {
			select {
			case <-ctx.Done():
				return
			case <-time.After(minTimeBetweenRuns - timeSinceLastRun):

			}
		}
	}
}

func run(ctx context.Context, profiles map[string]profile.ModeProfile) {
	for _, p := range profiles {
		go func(p profile.ModeProfile) {
			matches, pendingMatches, err := matchfunction.Run(ctx, p)
			if err != nil {
				logger.Info("Failed to fetch matches", zap.String("profileName", p.Name), zap.Error(err))
				return
			}

			logger.Info("Generated matches", zap.Int("generated", len(matches)),
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

			if err := assign(ctx, p, matches); err != nil {
				logger.Error("Failed to assign servers to matches", zap.Error(err))
				return
			}
		}(p)
	}
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

func assign(ctx context.Context, p profile.ModeProfile, matches []*pb.Match) error {
	// mock assign for now
	for _, match := range matches {
		match.Assignment = &pb.Assignment{
			ServerId:      "mock-gameserver-xxxx",
			ServerAddress: "0.0.0.0",
			ServerPort:    rand.Uint32(),
		}
		err := notifier.NotifyTransport(match)
		if err != nil {
			logger.Error("Failed to notify transport", zap.Error(err))
			continue
		}
	}
	return nil
}
