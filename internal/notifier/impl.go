package notifier

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/internal/rabbitmq/messenger"
	"kurushimi/pkg/pb"
	"log"
)

type countdownStreamContainer struct {
	stream         pb.Frontend_WatchTicketCountdownServer
	finishNotifier chan struct{}
}

type assignmentStreamContainer struct {
	stream         pb.Frontend_WatchTicketAssignmentServer
	finishNotifier chan struct{}
}

type notifierImpl struct {
	Notifier

	logger    *zap.SugaredLogger
	messenger messenger.Messenger

	countdownStreams map[string][]countdownStreamContainer

	assignmentStreams map[string][]assignmentStreamContainer
}

func NewNotifier(logger *zap.SugaredLogger, messenger messenger.Messenger) Notifier {
	return &notifierImpl{
		logger:            logger,
		messenger:         messenger,
		countdownStreams:  make(map[string][]countdownStreamContainer),
		assignmentStreams: make(map[string][]assignmentStreamContainer),
	}
}

func (n *notifierImpl) AddCountdownListener(ticketId string, stream pb.Frontend_WatchTicketCountdownServer, finishNotifier chan struct{}) {
	n.countdownStreams[ticketId] = append(n.countdownStreams[ticketId], countdownStreamContainer{
		stream:         stream,
		finishNotifier: finishNotifier,
	})
}

func (n *notifierImpl) RemoveCountdownListener(ticketId string) {
	n.countdownStreams[ticketId] = nil
}

func (n *notifierImpl) NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp) {
	for _, ticket := range tickets {
		for _, container := range n.countdownStreams[ticket.Id] {
			err := container.stream.Send(&pb.WatchCountdownResponse{
				TeleportTime: teleportTime,
			})

			if err != nil {
				n.logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func (n *notifierImpl) NotifyCountdownCancellation(tickets []*pb.Ticket) {
	for _, ticket := range tickets {
		for _, container := range n.countdownStreams[ticket.Id] {
			cancelled := true
			err := container.stream.Send(&pb.WatchCountdownResponse{
				Cancelled: &cancelled,
			})

			if err != nil {
				n.logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func (n *notifierImpl) AddAssignmentListener(ticketId string, stream pb.Frontend_WatchTicketAssignmentServer, finishNotifier chan struct{}) {
	n.assignmentStreams[ticketId] = append(n.assignmentStreams[ticketId], assignmentStreamContainer{
		stream:         stream,
		finishNotifier: finishNotifier,
	})
}

func (n *notifierImpl) RemoveAssignmentListener(ticketId string) {
	n.assignmentStreams[ticketId] = nil
}

func (n *notifierImpl) notifyAssignment(match *pb.Match) {
	if match.Assignment == nil {
		n.logger.Errorw("Assignment is nil", "match", match)
		return
	}
	for _, ticket := range match.Tickets {
		for _, container := range n.assignmentStreams[ticket.Id] {
			err := container.stream.Send(&pb.WatchAssignmentResponse{
				Assignment: match.Assignment,
			})

			if err != nil {
				n.logger.Errorw("Failed to send notification", err)
			}

			// An assignment is only sent once, so we can remove the container/container
			// even if there is an error.
			container.finishNotifier <- struct{}{}
		}
	}
}

func (n *notifierImpl) NotifyMatchTeleport(ctx context.Context, match *pb.Match) error {
	log.Printf("Match %s is ready to teleport", match.Id)
	err := n.notifyTransport(ctx, match)
	if err != nil {
		return err
	}
	log.Printf("Match %s has been notified", match.Id)
	n.notifyAssignment(match)
	log.Printf("Match %s has been notified of assignment", match.Id)

	// remove from countdowns
	for _, ticket := range match.Tickets {
		log.Printf("Removing countdown for ticket %s", ticket.Id)
		n.RemoveCountdownListener(ticket.Id)
		log.Printf("Removing assignment for ticket %s", ticket.Id)
		n.RemoveAssignmentListener(ticket.Id)
	}
	log.Printf("Match %s has been removed from countdowns", match.Id)
	return nil
}

func (n *notifierImpl) notifyTransport(ctx context.Context, match *pb.Match) error {
	// Get the player IDs ignoring tickets marked to not notify the proxy
	var playerIds []string
	for _, ticket := range match.Tickets {
		if ticket.NotifyProxy != nil && ticket.GetNotifyProxy() == false {
			continue
		}
		playerIds = append(playerIds, ticket.GetPlayerId())
	}

	// Check if len = 0 because NotifyProxy may have been set to false for all tickets
	if len(playerIds) == 0 {
		return nil
	}

	err := n.messenger.SwitchPlayerServer(ctx, match.Assignment, playerIds)

	if err != nil {
		n.logger.Errorw("failed to publish message", err)
	}

	return nil
}
