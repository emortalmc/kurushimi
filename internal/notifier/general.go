package notifier

import (
	"context"
	"github.com/emortalmc/grpc-api-specs/gen/go/service/player_tracker"
	kubernetes2 "k8s.io/client-go/kubernetes"
	"kurushimi/internal/messaging"
	"kurushimi/pkg/pb"
	"log"
)

type Notifier struct {
	kubeClient          kubernetes2.Interface
	messenger           *messaging.Messenger
	playerTrackerClient player_tracker.PlayerTrackerClient
}

func (n *Notifier) NotifyMatchTeleport(ctx context.Context, match *pb.Match) error {
	log.Printf("Match %s is ready to teleport", match.Id)
	err := n.notifyTransport(ctx, match)
	if err != nil {
		return err
	}
	notifyAssignment(match)

	// remove from countdowns
	for _, ticket := range match.Tickets {
		RemoveCountdownListener(ticket.Id)
	}
	return nil
}
