package notifier

import (
	"context"
	"kurushimi/pkg/pb"
	"log"
)

func NotifyMatchTeleport(ctx context.Context, match *pb.Match) error {
	log.Printf("Match %s is ready to teleport", match.Id)
	err := notifyTransport(ctx, match)
	if err != nil {
		return err
	}
	notifyAssignment(match)
	return nil
}
