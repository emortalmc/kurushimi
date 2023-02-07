package messenger

import (
	"context"
	"github.com/emortalmc/proto-specs/gen/go/message/common"
	commonmodel "github.com/emortalmc/proto-specs/gen/go/model/common"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"kurushimi/pkg/pb"
	"time"
)

type rabbitMqMessenger struct {
	Messenger

	conn   *amqp091.Connection
	chann  *amqp091.Channel
	logger *zap.SugaredLogger
}

func NewRabbitMQMessenger(logger *zap.SugaredLogger, conn *amqp091.Connection) (Messenger, error) {
	chann, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &rabbitMqMessenger{
		conn:   conn,
		chann:  chann,
		logger: logger,
	}, nil
}

func (m *rabbitMqMessenger) SwitchPlayerServer(ctx context.Context, assignment *pb.Assignment, playerIds []string) error {
	msg := common.SwitchPlayersServerMessage{
		Server: &commonmodel.ConnectableServer{
			Id:      assignment.ServerId,
			Address: assignment.ServerAddress,
			Port:    assignment.ServerPort,
		},
		PlayerIds: playerIds,
	}
	bytes, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	return m.chann.PublishWithContext(ctx, "mc:proxy:all", "", false, false, amqp091.Publishing{
		ContentType: "application/x-protobuf",
		Type:        string(msg.ProtoReflect().Descriptor().FullName()),
		Timestamp:   time.Now(),
		Body:        bytes,
	})
}
