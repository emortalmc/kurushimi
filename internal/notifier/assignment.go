package notifier

import (
	"go.uber.org/zap"
	"kurushimi/pkg/pb"
)

// map [ticketId][]StreamContainer
var assignmentStreams = make(map[string][]AssignmentStreamContainer)

type AssignmentStreamContainer struct {
	Stream         pb.Frontend_WatchTicketAssignmentServer
	FinishNotifier chan struct{}
}

type assignmentNotifierImpl struct {
	AssignmentNotifier
	logger *zap.SugaredLogger
}

func (n *assignmentNotifierImpl) AddAssignmentListener(ticketId string, stream pb.Frontend_WatchTicketAssignmentServer, finishNotifier chan struct{}) {
	assignmentStreams[ticketId] = append(assignmentStreams[ticketId], AssignmentStreamContainer{
		Stream:         stream,
		FinishNotifier: finishNotifier,
	})
}

func (n *assignmentNotifierImpl) notifyAssignment(match *pb.Match) {
	if match.Assignment == nil {
		n.logger.Errorw("Assignment is nil", "match", match)
		return
	}
	for _, ticket := range match.Tickets {
		for _, container := range assignmentStreams[ticket.Id] {
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

func (n *assignmentNotifierImpl) RemoveAssignmentStream(ticketId string) {
	n.logger.Info("Removing assignment stream")
	assignmentStreams[ticketId] = nil
}
