package protoutils

import "kurushimi/pkg/pb"

func GetMatchPlayerCount(match *pb.Match) int64 {
	count := 0

	for _, ticket := range match.Tickets {
		count += len(ticket.PlayerIds)
	}

	return int64(count)
}
