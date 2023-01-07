package notifier

import (
	"context"
	"kurushimi/pkg/pb"
)

func NotifyMatchTeleport(ctx context.Context, match *pb.Match) error {
	err := notifyTransport(ctx, match)
	if err != nil {
		return err
	}
	notifyAssignment(match)
	return nil
}
