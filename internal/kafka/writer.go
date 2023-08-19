package kafka

import (
	"context"
	"fmt"
	"github.com/emortalmc/kurushimi/internal/config"
	"github.com/emortalmc/kurushimi/internal/repository/model"
	"github.com/emortalmc/kurushimi/pkg/pb"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"time"
)

const writeTopic = "matchmaker"

// TODO fire these methods
type Notifier interface {
	TicketCreated(ctx context.Context, ticket *model.Ticket) error
	TicketUpdated(ctx context.Context, ticket *model.Ticket) error
	TicketDeleted(ctx context.Context, ticket *pb.Ticket, reason pb.TicketDeletedMessage_Reason) error

	PendingMatchCreated(ctx context.Context, match *model.PendingMatch) error
	PendingMatchUpdated(ctx context.Context, match *model.PendingMatch) error
	PendingMatchDeleted(ctx context.Context, match *model.PendingMatch, reason pb.PendingMatchDeletedMessage_Reason) error

	MatchCreated(ctx context.Context, match *pb.Match) error
}

type kafkaNotifier struct {
	w *kafka.Writer
}

func NewKafkaNotifier(config *config.KafkaConfig, logger *zap.SugaredLogger) Notifier {
	w := &kafka.Writer{
		Addr:         kafka.TCP(fmt.Sprintf("%s:%d", config.Host, config.Port)),
		Topic:        writeTopic,
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchTimeout: 100 * time.Millisecond,
		ErrorLogger:  kafka.LoggerFunc(logger.Errorw),
	}

	return &kafkaNotifier{w: w}
}

func (k *kafkaNotifier) TicketCreated(ctx context.Context, ticket *model.Ticket) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.TicketCreatedMessage{Ticket: ticket.ToProto()}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}

func (k *kafkaNotifier) TicketUpdated(ctx context.Context, ticket *model.Ticket) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.TicketUpdatedMessage{NewTicket: ticket.ToProto()}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}

func (k *kafkaNotifier) TicketDeleted(ctx context.Context, ticket *pb.Ticket, reason pb.TicketDeletedMessage_Reason) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.TicketDeletedMessage{Ticket: ticket, Reason: reason}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}

func (k *kafkaNotifier) PendingMatchCreated(ctx context.Context, match *model.PendingMatch) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.PendingMatchCreatedMessage{PendingMatch: match.ToProto()}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}

func (k *kafkaNotifier) PendingMatchUpdated(ctx context.Context, match *model.PendingMatch) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.PendingMatchUpdatedMessage{PendingMatch: match.ToProto()}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}

func (k *kafkaNotifier) PendingMatchDeleted(ctx context.Context, match *model.PendingMatch, reason pb.PendingMatchDeletedMessage_Reason) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.PendingMatchDeletedMessage{PendingMatch: match.ToProto(), Reason: reason}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}

func (k *kafkaNotifier) MatchCreated(ctx context.Context, match *pb.Match) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pMsg := &pb.MatchCreatedMessage{Match: match}
	bytes, err := proto.Marshal(pMsg)
	if err != nil {
		return err
	}

	err = k.w.WriteMessages(ctx, kafka.Message{
		Headers: []kafka.Header{{Key: "X-Proto-Type", Value: []byte(pMsg.ProtoReflect().Descriptor().FullName())}},
		Value:   bytes,
	})

	return err
}
