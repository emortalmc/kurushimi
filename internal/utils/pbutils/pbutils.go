package pbutils

import "kurushimi/pkg/pb"

func ParsePlayersFromMatch(match *pb.Match) []string {
	var players []string

	for _, ticket := range match.Tickets {
		players = append(players, ticket.PlayerId)
	}

	return players
}
