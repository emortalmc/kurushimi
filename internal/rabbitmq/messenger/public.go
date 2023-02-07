package messenger

import (
	"context"
	"kurushimi/pkg/pb"
)

type Messenger interface {
	SwitchPlayerServer(ctx context.Context, assignment *pb.Assignment, playerIds []string) error
}
