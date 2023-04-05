package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"kurushimi/internal/config"
	"kurushimi/pkg/pb"
	"log"
	"time"
)

func main() {
	cfg := &config.KafkaConfig{
		Host: "localhost",
		Port: 9092,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Topic:   "matchmaker",
	})

	// Only read new messages. Since this is for testing, we don't care about old messages.
	err := reader.SetOffsetAt(context.Background(), time.Now())
	if err != nil {
		panic(err)
	}

	consume(reader)
}

func consume(reader *kafka.Reader) {
	for range time.Tick(time.Millisecond * 250) {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		var protoType string
		for _, header := range m.Headers {
			if header.Key == "X-Proto-Type" {
				protoType = string(header.Value)
			}
		}
		if protoType == "" {
			log.Printf("no proto type found in message headers")
			continue
		}

		var message proto.Message
		switch protoType {
		case string((&pb.TicketCreatedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.TicketCreatedMessage{}
		case string((&pb.TicketUpdatedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.TicketUpdatedMessage{}
		case string((&pb.TicketDeletedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.TicketDeletedMessage{}
		case string((&pb.PendingMatchCreatedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.PendingMatchCreatedMessage{}
		case string((&pb.PendingMatchUpdatedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.PendingMatchUpdatedMessage{}
		case string((&pb.PendingMatchDeletedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.PendingMatchDeletedMessage{}
		case string((&pb.MatchCreatedMessage{}).ProtoReflect().Descriptor().FullName()):
			message = &pb.MatchCreatedMessage{}
		}

		typeName := message.ProtoReflect().Descriptor().Name()

		if err := proto.Unmarshal(m.Value, message); err != nil {
			log.Printf("error unmarshalling %v message: %v (original: %v)", typeName, err, string(m.Value))
			continue
		}
		log.Printf("received message %s: %+v", typeName, message)
	}
}
