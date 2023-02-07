package notifier

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/pkg/pb"
)

// todo: merge some of these methods (e.g notifyTransport, NotifyMatchTeleport)
type Notifier interface {
	AddCountdownListener(ticketId string, stream pb.Frontend_WatchTicketCountdownServer, finishNotifier chan struct{})
	RemoveCountdownListener(ticketId string)

	NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp)
	NotifyCountdownCancellation(tickets []*pb.Ticket)

	AddAssignmentListener(ticketId string, stream pb.Frontend_WatchTicketAssignmentServer, finishNotifier chan struct{})
	ClearAssignmentListeners(ticketId string)

	notifyAssignment(match *pb.Match)

	notifyTransport(ctx context.Context, match *pb.Match) error
	NotifyMatchTeleport(ctx context.Context, match *pb.Match) error
}
