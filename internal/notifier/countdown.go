package notifier

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/pkg/pb"
)

// map [ticketId][]CountdownStreamContainer
var countdownStreams = make(map[string][]CountdownStreamContainer)

type CountdownStreamContainer struct {
	Stream         pb.Frontend_WatchTicketCountdownServer
	FinishNotifier chan struct{}
}

type countdownNotifierImpl struct {
	CountdownNotifier
	logger *zap.SugaredLogger
}

func (n *countdownNotifierImpl) AddCountdownListener(ticketId string, stream pb.Frontend_WatchTicketCountdownServer, finishNotifier chan struct{}) {
	countdownStreams[ticketId] = append(countdownStreams[ticketId], CountdownStreamContainer{
		Stream:         stream,
		FinishNotifier: finishNotifier,
	})
}

func (n *countdownNotifierImpl) NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp) {
	for _, ticket := range tickets {
		for _, container := range countdownStreams[ticket.Id] {
			err := container.Stream.Send(&pb.WatchCountdownResponse{
				TeleportTime: teleportTime,
			})

			if err != nil {
				n.logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func (n *countdownNotifierImpl) NotifyCountdownCancellation(tickets []*pb.Ticket) {
	for _, ticket := range tickets {
		for _, container := range countdownStreams[ticket.Id] {
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

func (n *countdownNotifierImpl) RemoveCountdownListener(ticketId string) {
	for _, container := range countdownStreams[ticketId] {
		container.FinishNotifier <- struct{}{}
	}
	countdownStreams[ticketId] = nil
}
