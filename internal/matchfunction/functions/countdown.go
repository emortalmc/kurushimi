package functions

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/notifier"
	"kurushimi/internal/utils/math"
	"kurushimi/pkg/pb"
	"time"
)

// MakeCountdownMatches
// 1. Fill existing PendingMatches with new tickets from the pool
// Communicate to players in the PendingMatch the time until the match is made and the tickets in the match
// 2. Create new PendingMatches from the remaining tickets in the pool
// NOTE: tickets is all active tickets, even if they are already in a PendingMatch
func MakeCountdownMatches(notifier notifier.Notifier, profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) {
	if len(tickets) == 0 {
		return nil, pendingMatches, nil
	}
	pendingMatches = handleLeavers(notifier, profile, pendingMatches, tickets)

	// filter out tickets already in a pending match
	tickets = filterTickets(tickets, pendingMatches)

	tickets, pendingMatches = fillPendingMatches(notifier, profile, tickets, pendingMatches)

	if len(tickets) == 0 {
		return nil, pendingMatches, nil
	}

	tickets, madePendingMatches := makePendingMatches(notifier, profile, tickets)

	pendingMatches = append(pendingMatches, madePendingMatches...)

	return nil, pendingMatches, nil
}

// handleLeavers remove tickets from pending matches that have left the game.
// returns: updated pending matches as some has been removed
func handleLeavers(notifier notifier.Notifier, profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) []*pb.PendingMatch {
	for _, pendingMatch := range pendingMatches {
		updated := false
		newMatchTickets := make([]*pb.Ticket, 0)
		for _, matchTicket := range pendingMatch.Tickets {
			found := false
			for _, ticket := range tickets {
				if matchTicket.Id == ticket.Id {
					found = true
					break
				}
			}
			if !found {
				updated = true
			} else {
				newMatchTickets = append(newMatchTickets, matchTicket)
			}
		}
		if updated {
			pendingMatch.Tickets = newMatchTickets
			if len(pendingMatch.Tickets) < profile.MinPlayers {
				// We don't have to worry about the removed tickets because they're not in the queue anymore.
				notifier.NotifyCountdownCancellation(pendingMatch.Tickets)
			} else {
				notifier.NotifyCountdown(pendingMatch.Tickets, pendingMatch.TeleportTime)
			}
		}
	}

	return pendingMatches
}

func filterTickets(tickets []*pb.Ticket, pendingMatches []*pb.PendingMatch) []*pb.Ticket {
	filteredTickets := make([]*pb.Ticket, 0)
	for _, ticket := range tickets {
		found := false
		for _, pendingMatch := range pendingMatches {
			for _, matchTicket := range pendingMatch.Tickets {
				if ticket.Id == matchTicket.Id {
					found = true
					break
				}
			}
		}
		if !found {
			filteredTickets = append(filteredTickets, ticket)
		}
	}
	return filteredTickets
}

func fillPendingMatches(notifier notifier.Notifier, profile profile.ModeProfile, tickets []*pb.Ticket, pendingMatches []*pb.PendingMatch) ([]*pb.Ticket, []*pb.PendingMatch) {
	for _, pendingMatch := range pendingMatches {
		if len(pendingMatch.Tickets) >= profile.MaxPlayers {
			continue
		}

		// fill the pending match
		updated := false
		for len(pendingMatch.Tickets) < profile.MaxPlayers && len(tickets) > 0 {
			pendingMatch.Tickets = append(pendingMatch.Tickets, tickets[0])
			tickets = tickets[1:]
			updated = true
		}

		if updated {
			notifier.NotifyCountdown(pendingMatch.Tickets, pendingMatch.TeleportTime)
		}

		if len(tickets) == 0 {
			break
		}
	}
	return tickets, pendingMatches
}

func makePendingMatches(notifier notifier.Notifier, profile profile.ModeProfile, tickets []*pb.Ticket) ([]*pb.Ticket, []*pb.PendingMatch) {
	pendingMatches := make([]*pb.PendingMatch, 0)

	for len(tickets) >= profile.MinPlayers {
		maxIndex := math.Min(len(tickets), profile.MaxPlayers)
		pendingMatch := &pb.PendingMatch{
			Id:           uuid.New().String(),
			ProfileName:  profile.Name,
			Tickets:      tickets[:maxIndex],
			TeleportTime: timestamppb.New(time.Now().Add(10 * time.Second)),
		}
		pendingMatches = append(pendingMatches, pendingMatch)
		tickets = tickets[maxIndex:]

		notifier.NotifyCountdown(pendingMatch.Tickets, pendingMatch.TeleportTime)
	}

	return tickets, pendingMatches
}

func newMatch(id uuid.UUID, profile profile.ModeProfile, tickets []*pb.Ticket) *pb.Match {
	match := &pb.Match{
		Id:      id.String(),
		Tickets: tickets,
	}

	return match
}
