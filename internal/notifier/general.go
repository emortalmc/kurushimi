package notifier

import (
	"context"
	"kurushimi/pkg/pb"
	"log"
)

//type Notifier struct {
//	kubeClient          kubernetes2.Interface
//	messenger           *messaging.Messenger
//	playerTrackerClient player_tracker.PlayerTrackerClient
//}

type notifierImpl struct {
	Notifier
}

func NewNotifier() Notifier {
	return &notifierImpl{
		CountdownNotifier: NewCountdownNotifier(),
	}
}

func (n *notifierImpl) NotifyMatchTeleport(ctx context.Context, match *pb.Match) error {
	log.Printf("Match %s is ready to teleport", match.Id)
	err := n.notifyTransport(ctx, match)
	if err != nil {
		return err
	}
	n.notifyAssignment(match)

	// remove from countdowns
	for _, ticket := range match.Tickets {
		n.RemoveCountdownListener(ticket.Id)
	}
	return nil
}
