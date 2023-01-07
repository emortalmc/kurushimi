package notifier

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/pkg/pb"
)

// map [ticketId][]pb.Frontend_WatchTicketCountdownServer
var listeners = make(map[string][]pb.Frontend_WatchTicketCountdownServer)
var logger, _ = zap.NewDevelopment()

func AddCountdownListener(ticketId string, listener pb.Frontend_WatchTicketCountdownServer) {
	listeners[ticketId] = append(listeners[ticketId], listener)
}

func NotifyCountdown(tickets []*pb.Ticket, teleportTime *timestamppb.Timestamp) {
	for _, ticket := range tickets {
		for _, listener := range listeners[ticket.Id] {
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
		for _, listener := range listeners[ticket.Id] {
			err := listener.Send(&pb.WatchCountdownResponse{
				Cancelled: true,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func RemoveCountdownListener(ticketId string) {
	listeners[ticketId] = nil
}
