package matchfunction

import (
	"context"
	"go.uber.org/zap"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/statestore"
	"kurushimi/pkg/pb"
)

var logger = zap.S()

func Run(ctx context.Context, profile profile.ModeProfile) ([]*pb.Match, []*pb.PendingMatch, error) {
	tickets, err := statestore.GetAllTickets(ctx, profile.GameName)
	pendingMatches, err := statestore.GetAllPendingMatches(ctx, profile)
	if err != nil {
		return nil, nil, err
	}

	logger.Info("Running match function", zap.String("profileName", profile.Name), zap.Int("ticketCount", len(tickets)), zap.Int("pendingMatchCount", len(pendingMatches)))

	return profile.MatchFunction(profile, pendingMatches, tickets)
}
