package statestore

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kurushimi/pkg/pb"
)

const (
	ticketPrefix = "ticket:"
	allTickets   = ticketPrefix + "all"
)

func (rs *redisStore) CreateTicket(ctx context.Context, ticket *pb.Ticket) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	value, err := proto.Marshal(ticket)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to marshal ticket %v", err)
	}

	_, err = redisConn.Do("SET", ticketPrefix+ticket.Id, value)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to set ticket %v", err)
	}

	return nil
}

func (rs *redisStore) GetTicket(ctx context.Context, ticketId string) (*pb.Ticket, error) {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	value, err := redis.Bytes(redisConn.Do("GET", ticketPrefix+ticketId))
	if err != nil {
		if err == redis.ErrNil {
			return nil, status.Errorf(codes.NotFound, "Ticket %s not found", ticketId)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get ticket %v", err)
	}

	if value == nil {
		return nil, status.Errorf(codes.NotFound, "Ticket %s not found", ticketId)
	}

	ticket := &pb.Ticket{}
	err = proto.Unmarshal(value, ticket)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unmarshal ticket %v", err)
	}

	return ticket, nil
}

func (rs *redisStore) GetAllTickets(ctx context.Context, game string) ([]*pb.Ticket, error) {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	ticketIds, err := redis.Strings(redisConn.Do("SMEMBERS", allTickets))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get all tickets %v", err)
	}

	tickets := make([]*pb.Ticket, 0, len(ticketIds))
	for _, ticketId := range ticketIds {
		ticket, err := rs.GetTicket(ctx, ticketId)
		if err != nil {
			return nil, err
		}
		// filter by game
		for _, tag := range ticket.SearchFields.Tags {
			if tag == game {
				tickets = append(tickets, ticket)
				continue
			}
		}
	}

	return tickets, nil
}

func (rs *redisStore) DeleteTicket(ctx context.Context, ticketId string) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("DEL", ticketPrefix+ticketId)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to delete ticket %v", err)
	}
	return nil
}

func (rs *redisStore) IndexTicket(ctx context.Context, ticket *pb.Ticket) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("SADD", allTickets, ticket.Id)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to index ticket %v", err)
	}
	return nil
}

func (rs *redisStore) UnIndexTicket(ctx context.Context, ticketId string) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("SREM", allTickets, ticketId)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to deindex ticket %v", err)
	}
	return nil
}
