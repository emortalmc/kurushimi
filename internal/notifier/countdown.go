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

func AddCountdownListener(ticketId string, stream pb.Frontend_WatchTicketCountdownServer, finishNotifier chan struct{}) {
	countdownStreams[ticketId] = append(countdownStreams[ticketId], CountdownStreamContainer{
		Stream:         stream,
		FinishNotifier: finishNotifier,
	})
}

func NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp) {
	logger := zap.S()
	for _, ticket := range tickets {
		for _, container := range countdownStreams[ticket.Id] {
			err := container.Stream.Send(&pb.WatchCountdownResponse{
				TeleportTime: teleportTime,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func NotifyCountdownCancellation(tickets []*pb.Ticket) {
	logger := zap.S()
	for _, ticket := range tickets {
		for _, container := range countdownStreams[ticket.Id] {
			cancelled := true
			err := container.Stream.Send(&pb.WatchCountdownResponse{
				Cancelled: &cancelled,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func RemoveCountdownListener(ticketId string) {
	for _, container := range countdownStreams[ticketId] {
		container.FinishNotifier <- struct{}{}
	}
	countdownStreams[ticketId] = nil
}
