package profile

import (
	v1 "agones.dev/agones/pkg/apis/allocation/v1"
	"kurushimi/pkg/pb"
)

type ModeProfile struct {
	Name      string `json:"name"`
	FleetName string `json:"fleetName"`
	GameName  string `json:"gameName"`

	Selector func(profile ModeProfile, match *pb.Match) *v1.GameServerAllocation `json:"-"`

	MinPlayers    int                                                                                                                         `json:"minPlayers"`
	MaxPlayers    int                                                                                                                         `json:"maxPlayers"`
	MatchFunction func(profile ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) `json:"-"`
}
