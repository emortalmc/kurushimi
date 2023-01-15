package matchfunction

import (
	"context"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/statestore"
	"kurushimi/pkg/pb"
)

func Run(ctx context.Context, stateStore statestore.StateStore, profile profile.ModeProfile) ([]*pb.Match, []*pb.PendingMatch, error) {
	tickets, err := stateStore.GetAllTickets(ctx, profile.GameName)
	pendingMatches, err := stateStore.GetAllPendingMatches(ctx, profile)
	if err != nil {
		return nil, nil, err
	}

	return profile.MatchFunction(profile, pendingMatches, tickets)
}
