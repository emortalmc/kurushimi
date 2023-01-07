package notifier

import (
	"go.uber.org/zap"
	"kurushimi/pkg/pb"
)

// map [ticketId][]pb.Frontend_WatchTicketAssignmentServer
var assignment = make(map[string][]pb.Frontend_WatchTicketAssignmentServer)
var logger = zap.S()

func AddAssignmentListener(ticketId string, listener pb.Frontend_WatchTicketAssignmentServer) {
	assignment[ticketId] = append(assignment[ticketId], listener)
}

func notifyAssignment(match *pb.Match) {
	if match.Assignment == nil {
		logger.Error("Assignment is nil", zap.Any("match", match))
		return
	}
	for _, ticket := range match.Tickets {
		for _, listener := range assignment[ticket.Id] {
			err := listener.Send(&pb.WatchAssignmentResponse{
				Assignment: match.Assignment,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}
		}
	}
}

func RemoveAssignmentListener(ticketId string) {
	countdownListeners[ticketId] = nil
}
