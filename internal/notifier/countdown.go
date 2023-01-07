package notifier

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/pkg/pb"
)

// map [ticketId][]pb.Frontend_WatchTicketCountdownServer
var countdownListeners = make(map[string][]pb.Frontend_WatchTicketCountdownServer)

func AddCountdownListener(ticketId string, listener pb.Frontend_WatchTicketCountdownServer) {
	countdownListeners[ticketId] = append(countdownListeners[ticketId], listener)
}

func NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp) {
	for _, ticket := range tickets {
		for _, listener := range countdownListeners[ticket.Id] {
			err := listener.Send(&pb.WatchCountdownResponse{
				TeleportTime: teleportTime,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func NotifyCountdownCancellation(tickets []*pb.Ticket) {
	for _, ticket := range tickets {
		for _, listener := range countdownListeners[ticket.Id] {
			cancelled := true
			err := listener.Send(&pb.WatchCountdownResponse{
				Cancelled: &cancelled,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func RemoveCountdownListener(ticketId string) {
	countdownListeners[ticketId] = nil
}
