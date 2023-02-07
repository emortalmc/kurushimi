package notifier

import (
	"context"
	"go.uber.org/zap"
	"kurushimi/internal/rabbitmq/messenger"
	"kurushimi/pkg/pb"
)

type proxyNotifierImpl struct {
	ProxyNotifier

	messenger messenger.Messenger
	logger    *zap.SugaredLogger
}

func (n *proxyNotifierImpl) notifyTransport(ctx context.Context, match *pb.Match) error {
	// Get the player IDs ignoring tickets marked to not notify the proxy
	var playerIds []string
	for _, ticket := range match.Tickets {
		if ticket.NotifyProxy != nil && ticket.GetNotifyProxy() == false {
			continue
		}
		playerIds = append(playerIds, ticket.GetPlayerId())
	}

	// Check if len = 0 because NotifyProxy may have been set to false for all tickets
	if len(playerIds) == 0 {
		return nil
	}

	err := n.messenger.SwitchPlayerServer(ctx, match.Assignment, playerIds)

	if err != nil {
		n.logger.Errorw("failed to publish message", err)
	}

	return nil
}
