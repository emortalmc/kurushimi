package model

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/pkg/pb"
	"time"
)

type QueuedPlayer struct {
	PlayerId uuid.UUID          `bson:"_id"`
	TicketId primitive.ObjectID `bson:"ticketId"`

	// MapId the map the player has voted for, nil if not voted
	MapId *string `bson:"mapId,omitempty"`
}

type Ticket struct {
	Id primitive.ObjectID `bson:"_id"`

	Removals *TicketRemovals `bson:"removals,omitempty"`

	// InPendingMatch is true if the ticket is currently in a pending match
	InPendingMatch bool `bson:"inPendingMatch"`

	// party fields are only present if the ticket is for a party
	PartyId       primitive.ObjectID    `bson:"partyId"`
	PartySettings *ReducedPartySettings `bson:"partySettings"`

	PlayerIds  []uuid.UUID `bson:"playerIds"`
	GameModeId string      `bson:"gameModeId"`

	AutoTeleport bool `bson:"autoTeleport"`

	InternalUpdates *TicketInternalUpdates `bson:"-"`
}

func NewTicket(partyId primitive.ObjectID, partySettings *ReducedPartySettings, playerIds []uuid.UUID,
	gameModeId string, autoTeleport bool) *Ticket {

	return &Ticket{
		Id:             primitive.NewObjectID(),
		Removals:       nil,
		InPendingMatch: false,
		PartyId:        partyId,
		PartySettings:  partySettings,
		PlayerIds:      playerIds,
		GameModeId:     gameModeId,
		AutoTeleport:   autoTeleport,
	}
}

func (t *Ticket) ToProto() *pb.Ticket {
	pbPlayerIds := make([]string, len(t.PlayerIds))
	for i, playerId := range t.PlayerIds {
		pbPlayerIds[i] = playerId.String()
	}

	return &pb.Ticket{
		Id:             t.Id.Hex(),
		PartyId:        t.PartyId.Hex(),
		CreatedAt:      timestamppb.New(t.Id.Timestamp()),
		PlayerIds:      pbPlayerIds,
		GameModeId:     t.GameModeId,
		AutoTeleport:   t.AutoTeleport,
		InPendingMatch: t.InPendingMatch,
	}
}

type TicketRemovals struct {
	MarkedForRemoval  bool        `bson:"markedForRemoval"`
	PlayersForRemoval []uuid.UUID `bson:"playersForRemoval"`
}

type ReducedPartySettings struct {
	LeaderId uuid.UUID `bson:"leaderId"`

	DequeueOnDisconnect bool `bson:"dequeueOnDisconnect"`
	AllowMemberDequeue  bool `bson:"allowMemberDequeue"`
}

// TicketInternalUpdates are updates not stored in the database but used internally once extracted.
// e.g. marking when a ticket's InPendingMatchUpdated field is updated
// the presence of this struct indicates that the ticket has been updated
type TicketInternalUpdates struct {
	InPendingMatchUpdated bool
}

func (t *Ticket) UpdateInPendingMach(value bool) {
	t.InPendingMatch = value
	if t.InternalUpdates == nil {
		t.InternalUpdates = &TicketInternalUpdates{}
	}

	t.InternalUpdates.InPendingMatchUpdated = true
}

// TODO use this method
// MarkUpdated doesn't perform any updates but marks the ticket as updated.
// This is currently only used by the director for Kafka events.
func (t *Ticket) MarkUpdated() {
	if t.InternalUpdates == nil {
		t.InternalUpdates = &TicketInternalUpdates{}
	}
}

type PendingMatch struct {
	Id primitive.ObjectID `bson:"_id"`

	GameModeId string `bson:"gameModeId"`

	TicketIds []primitive.ObjectID `bson:"ticketIds"`

	TeleportTime *time.Time `bson:"teleportTime"`
}

func (m *PendingMatch) ToProto() *pb.PendingMatch {
	pbTicketIds := make([]string, len(m.TicketIds))
	for i, ticketId := range m.TicketIds {
		pbTicketIds[i] = ticketId.Hex()
	}

	return &pb.PendingMatch{
		Id:           m.Id.Hex(),
		GameModeId:   m.GameModeId,
		TicketIds:    pbTicketIds,
		TeleportTime: timestamppb.New(*m.TeleportTime),
	}
}

// Backfill represents a backfill that is currently available.
// backfills are not yet finished or enabled :)
type Backfill struct {
	Id primitive.ObjectID `bson:"_id"`

	ServerInfo *ServerInfo `bson:"serverInfo"`
}

type ServerInfo struct {
	Address string `bson:"address"`
	Port    uint32 `bson:"port"`
	Name    string `bson:"name"`
}
