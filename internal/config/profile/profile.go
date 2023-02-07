package profile

import (
	v1 "agones.dev/agones/pkg/apis/allocation/v1"
	"kurushimi/internal/notifier"
	"kurushimi/pkg/pb"
	"time"
)

type ModeProfile struct {
	Name            string        `json:"name"`
	FleetName       string        `json:"fleetName"`
	GameName        string        `json:"gameName"`
	MatchmakingRate time.Duration `json:"matchmakingRate"`

	Selector func(profile ModeProfile, match *pb.Match) *v1.GameServerAllocation `json:"-"`

	MinPlayers    int                                                                                                                                                     `json:"minPlayers"`
	MaxPlayers    int                                                                                                                                                     `json:"maxPlayers"`
	MatchFunction func(notifier notifier.Notifier, profile ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) `json:"-"`
}
