package protoutils

import "github.com/emortalmc/kurushimi/pkg/pb"

func GetMatchPlayerCount(match *pb.Match) int64 {
	count := 0

	for _, ticket := range match.Tickets {
		count += len(ticket.PlayerIds)
	}

	return int64(count)
}
