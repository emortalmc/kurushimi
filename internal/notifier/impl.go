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
	Stream         pb.Frontend_WatchTicketCountdownServer
	FinishNotifier chan struct{}
}

type assignmentStreamContainer struct {
	Stream         pb.Frontend_WatchTicketAssignmentServer
	FinishNotifier chan struct{}
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
		Stream:         stream,
		FinishNotifier: finishNotifier,
	})
}

func (n *notifierImpl) RemoveCountdownListener(ticketId string) {
	for _, container := range n.countdownStreams[ticketId] {
		container.FinishNotifier <- struct{}{}
	}
	n.countdownStreams[ticketId] = nil
}

func (n *notifierImpl) NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp) {
	for _, ticket := range tickets {
		for _, container := range n.countdownStreams[ticket.Id] {
			err := container.Stream.Send(&pb.WatchCountdownResponse{
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
			err := container.Stream.Send(&pb.WatchCountdownResponse{
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
		Stream:         stream,
		FinishNotifier: finishNotifier,
	})
}

func (n *notifierImpl) RemoveAssignmentListener(ticketId string) {
	for _, container := range n.assignmentStreams[ticketId] {
		container.FinishNotifier <- struct{}{}
	}
	n.assignmentStreams[ticketId] = nil
}

func (n *notifierImpl) RemoveAssignmentStream(ticketId string) {
	n.logger.Info("Removing assignment stream")
	n.assignmentStreams[ticketId] = nil
}

func (n *notifierImpl) notifyAssignment(match *pb.Match) {
	if match.Assignment == nil {
		n.logger.Errorw("Assignment is nil", "match", match)
		return
	}
	for _, ticket := range match.Tickets {
		for _, container := range n.assignmentStreams[ticket.Id] {
			err := container.Stream.Send(&pb.WatchAssignmentResponse{
				Assignment: match.Assignment,
			})

			if err != nil {
				n.logger.Errorw("Failed to send notification", err)
			}

			// An assignment is only sent once, so we can remove the container/container
			// even if there is an error.
			container.FinishNotifier <- struct{}{}
		}
	}
}

func (n *notifierImpl) NotifyMatchTeleport(ctx context.Context, match *pb.Match) error {
	log.Printf("Match %s is ready to teleport", match.Id)
	err := n.notifyTransport(ctx, match)
	if err != nil {
		return err
	}
	n.notifyAssignment(match)

	// remove from countdowns
	for _, ticket := range match.Tickets {
		n.RemoveCountdownListener(ticket.Id)
		n.RemoveAssignmentListener(ticket.Id)
	}
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
