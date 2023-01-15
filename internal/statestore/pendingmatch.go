package statestore

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kurushimi/internal/config/profile"
	"kurushimi/pkg/pb"
)

const (
	pendingMatchPrefix = "pendingMatch:"
	allPendingMatches  = pendingMatchPrefix + "all"
)

func (rs *redisStore) CreatePendingMatch(ctx context.Context, pMatch *pb.PendingMatch) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	value, err := proto.Marshal(pMatch)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to marshal pending match %v", err)
	}

	_, err = redisConn.Do("SET", pendingMatchPrefix+pMatch.Id, value)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to set pending match %v", err)
	}
	return nil
}

func (rs *redisStore) GetPendingMatch(ctx context.Context, pMatchId string) (*pb.PendingMatch, error) {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	value, err := redis.Bytes(redisConn.Do("GET", pendingMatchPrefix+pMatchId))
	if err != nil {
		if err == redis.ErrNil {
			return nil, status.Errorf(codes.NotFound, "Pending match %s not found", pMatchId)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get pending match %v", err)
	}

	if value == nil {
		return nil, status.Errorf(codes.NotFound, "Pending match %s not found", pMatchId)
	}

	pMatch := &pb.PendingMatch{}
	err = proto.Unmarshal(value, pMatch)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unmarshal pending match %v", err)
	}

	return pMatch, nil
}

func (rs *redisStore) GetAllPendingMatches(ctx context.Context, profile profile.ModeProfile) ([]*pb.PendingMatch, error) {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	value, err := redis.Strings(redisConn.Do("SMEMBERS", allPendingMatches))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get all pending matches %v", err)
	}

	pMatches := make([]*pb.PendingMatch, 0)
	for _, pMatchId := range value {
		pMatch, err := rs.GetPendingMatch(ctx, pMatchId)
		if err != nil {
			return nil, err
		}
		if pMatch.ProfileName == profile.Name {
			pMatches = append(pMatches, pMatch)
		}
	}

	return pMatches, nil
}

func (rs *redisStore) DeletePendingMatch(ctx context.Context, pMatchId string) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("DEL", pendingMatchPrefix+pMatchId)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to delete pending match %v", err)
	}

	return nil
}

func (rs *redisStore) IndexPendingMatch(ctx context.Context, pMatch *pb.PendingMatch) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("SADD", allPendingMatches, pMatch.Id)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to index pending match %v", err)
	}

	return nil
}

func (rs *redisStore) UnIndexPendingMatch(ctx context.Context, pMatch *pb.PendingMatch) error {
	redisConn, err := rs.redisPool.GetContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to connect to redis: %v", err)
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("SREM", allPendingMatches, pMatch.Id)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to unindex pending match %v", err)
	}

	return nil
}
