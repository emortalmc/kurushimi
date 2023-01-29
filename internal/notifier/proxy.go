package notifier

import (
	"context"
	"github.com/emortalmc/grpc-api-specs/gen/go/messaging/general"
	"github.com/emortalmc/grpc-api-specs/gen/go/service/server_discovery"
	"github.com/golang/protobuf/proto"
	"github.com/rabbitmq/amqp091-go"
	"kurushimi/internal/messaging"
	"kurushimi/internal/utils/kubernetes"
	"kurushimi/pkg/pb"
	"log"
	"time"
)

func New(messenger *messaging.Messenger) *Notifier {
	return &Notifier{
		kubeClient: kubernetes.KubeClient,
		messenger:  messenger,
	}
}

func (n *Notifier) notifyTransport(ctx context.Context, match *pb.Match) error {
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

	assignment := match.Assignment
	msg := general.ProxyServerSwitchMessage{
		Server: &server_discovery.ConnectableServer{
			Id:      assignment.ServerId,
			Address: assignment.ServerAddress,
			Port:    assignment.ServerPort,
		},
		PlayerIds: playerIds,
	}
	msgBytes, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	log.Printf("Publishing message: " + msg.String())
	err = n.messenger.Channel.PublishWithContext(ctx, "mc:proxy:all", "", false, false, amqp091.Publishing{
		Timestamp: time.Now(),
		Type:      string(msg.ProtoReflect().Descriptor().FullName()),
		Body:      msgBytes,
	})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	return nil
}
