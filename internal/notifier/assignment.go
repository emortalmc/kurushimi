package notifier

import (
	"go.uber.org/zap"
	"kurushimi/pkg/pb"
	"log"
)

// map [ticketId][]StreamContainer
var assignmentStreams = make(map[string][]AssignmentStreamContainer)

type AssignmentStreamContainer struct {
	Stream         pb.Frontend_WatchTicketAssignmentServer
	FinishNotifier chan struct{}
}

func AddAssignmentListener(ticketId string, stream pb.Frontend_WatchTicketAssignmentServer, finishNotifier chan struct{}) {
	assignmentStreams[ticketId] = append(assignmentStreams[ticketId], AssignmentStreamContainer{
		Stream:         stream,
		FinishNotifier: finishNotifier,
	})
}

func notifyAssignment(match *pb.Match) {
	logger := zap.S()
	if match.Assignment == nil {
		logger.Errorw("Assignment is nil", "match", match)
		return
	}
	for _, ticket := range match.Tickets {
		for _, container := range assignmentStreams[ticket.Id] {
			err := container.Stream.Send(&pb.WatchAssignmentResponse{
				Assignment: match.Assignment,
			})

			if err != nil {
				logger.Errorw("Failed to send notification", err)
			}

			log.Printf("Marking context as done")
			// An assignment is only sent once, so we can remove the container/container
			// even if there is an error.
			container.FinishNotifier <- struct{}{}
		}
	}
}

func RemoveAssignmentStream(ticketId string) {
	log.Printf("Removing assignment stream")
	assignmentStreams[ticketId] = nil
}
