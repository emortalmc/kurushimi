package statestore

import (
	"context"
	"kurushimi/internal/config/profile"
	"kurushimi/pkg/pb"
)

type StateStore interface {
	CreatePendingMatch(ctx context.Context, pMatch *pb.PendingMatch) error

	GetPendingMatch(ctx context.Context, pMatchId string) (*pb.PendingMatch, error)

	GetAllPendingMatches(ctx context.Context, profile profile.ModeProfile) ([]*pb.PendingMatch, error)

	DeletePendingMatch(ctx context.Context, pMatchId string) error

	IndexPendingMatch(ctx context.Context, pMatch *pb.PendingMatch) error

	UnIndexPendingMatch(ctx context.Context, pMatch *pb.PendingMatch) error

	CreateTicket(ctx context.Context, ticket *pb.Ticket) error

	GetTicket(ctx context.Context, ticketId string) (*pb.Ticket, error)

	GetAllTickets(ctx context.Context, game string) ([]*pb.Ticket, error)

	DeleteTicket(ctx context.Context, ticketId string) error

	IndexTicket(ctx context.Context, ticket *pb.Ticket) error

	UnIndexTicket(ctx context.Context, ticketId string) error
}
