package functions

import (
	"github.com/google/uuid"
	"kurushimi/internal/config/profile"
	"kurushimi/pkg/pb"
	"log"
)

// MakeInstantMatches
// Immediately returns a match with all tickets in the pool
// but groups them together to reduce allocations
func MakeInstantMatches(profile profile.ModeProfile, tickets []*pb.Ticket) ([]*pb.Match, error) {
	if len(tickets) < profile.MinPlayers {
		return nil, nil
	}

	var matches []*pb.Match
	for len(tickets) >= profile.MinPlayers {
		var matchTickets []*pb.Ticket
		for i := 0; i < profile.MaxPlayers && len(tickets) > 0; i++ {
			ticket := tickets[0]
			// Remove the Tickets from this pool and add to the match proposal.
			matchTickets = append(matchTickets, ticket)
			tickets = tickets[1:]
		}

		log.Printf("Creating instant %s match with %d tickets", profile.Name, len(matchTickets))
		match := newMatch(uuid.New(), profile, matchTickets)
		matches = append(matches, match)
	}

	return matches, nil
}
