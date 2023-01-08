package notifier

import (
	"go.uber.org/zap"
	"kurushimi/pkg/pb"
	"log"
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
	log.Printf("Notifying assignment for match %s (%s)", match.Id, assignment)
	for _, ticket := range match.Tickets {
		for _, listener := range assignment[ticket.Id] {
			log.Printf("Notifying assignment for ticket %s", ticket.Id)
			err := listener.Send(&pb.WatchAssignmentResponse{
				Assignment: match.Assignment,
			})

			if err != nil {
				logger.Error("Failed to send notification", zap.Error(err))
			}

			// An assignment is only sent once, so we can remove the listener/stream
			// even if there is an error.
			listener.Context().Done()
		}
	}
}

func RemoveAssignmentListener(ticketId string) {
	logger.Info("Removing assignment listener", zap.String("ticketId", ticketId))
	countdownListeners[ticketId] = nil
}
