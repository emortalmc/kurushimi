package kafka

import (
	"context"
	"fmt"
	"github.com/emortalmc/proto-specs/gen/go/message/party"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"kurushimi/internal/config"
	"kurushimi/internal/repository"
	"log"
)

// NOTE: We don't listen to player connections as player disconnect = party leave or disband - so it's handled by that.
const partyTopic = "party-manager"

type consumer struct {
	logger *zap.SugaredLogger

	partyReader *kafka.Reader

	repo repository.Repository
}

// NewConsumer is blocking
func NewConsumer(config *config.KafkaConfig, logger *zap.SugaredLogger, repo repository.Repository) {
	partyReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", config.Host, config.Port)},
		GroupID: "matchmaker",
		Topic:   partyTopic,
	})

	c := &consumer{
		logger: logger,

		partyReader: partyReader,

		repo: repo,
	}

	logger.Infow("listening for messages on topic", "topic", partyTopic)
	c.consume()
}

func (c *consumer) consume() {
	for {
		m, err := c.partyReader.ReadMessage(context.Background())
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
			c.logger.Warnw("no proto type found in message headers")
			continue
		}

		switch protoType {
		case string((&party.PartyDisbandedMessage{}).ProtoReflect().Descriptor().FullName()):
			var msg *party.PartyDisbandedMessage
			if err := proto.Unmarshal(m.Value, msg); err != nil {
				c.logger.Errorw("failed to unmarshal message", err)
				continue
			}

			log.Printf("received message: %s", msg.String())
		case string((&party.PartyPlayerJoinedMessage{}).ProtoReflect().Descriptor().FullName()):
			var msg *party.PartyPlayerJoinedMessage
			if err := proto.Unmarshal(m.Value, msg); err != nil {
				c.logger.Errorw("failed to unmarshal message", err)
				continue
			}
		case string((&party.PartyPlayerLeftMessage{}).ProtoReflect().Descriptor().FullName()):

		}
	}
}

func (c *consumer) handlePartyDisband(msg *party.PartyDisbandedMessage) {
	// Dequeue if possible
	partyId, err := primitive.ObjectIDFromHex(msg.Party.Id)
	if err != nil {
		c.logger.Errorw("failed to parse party id", err)
		return
	}

	_, err = c.repo.AddTicketDequeueRequestByPartyId(context.TODO(), partyId)
	if err != nil {
		c.logger.Errorw("failed to add ticket dequeue request", err)
		return
	}
}

// TODO we don't account for if a ticket size increases and there isn't enough space in their existing PendingMatch
func (c *consumer) handlePartyPlayerJoined(msg *party.PartyPlayerJoinedMessage) {
	// TODO
}

func (c *consumer) handlePartyPlayerLeft(msg *party.PartyPlayerLeftMessage) {
	partyId, err := primitive.ObjectIDFromHex(msg.PartyId)
	if err != nil {
		c.logger.Errorw("failed to parse party id", err)
		return
	}

	playerId, err := uuid.Parse(msg.Member.Id)
	if err != nil {
		c.logger.Errorw("failed to parse player id", err)
		return
	}

	err = c.repo.AddPlayerDequeueRequestByPartyId(context.TODO(), partyId, playerId)
	if err != nil {
		c.logger.Errorw("failed to add player dequeue request", err)
		return
	}
}
