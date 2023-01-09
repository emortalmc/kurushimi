package config

import (
	v1 "agones.dev/agones/pkg/apis/allocation/v1"
	"kurushimi/internal/config/profile"
	"kurushimi/internal/config/selector"
	"kurushimi/internal/matchfunction/functions"
	"kurushimi/pkg/pb"
	"time"
)

var ModeProfiles = map[string]profile.ModeProfile{
	//"marathon": {
	//	Name:       "marathon",
	//	GameName:   "game.marathon",
	//	FleetName:  "marathon",
	//	MinPlayers: 1,
	//	MaxPlayers: 100, // The server knows this is a singleplayer game, we can reduce load by using backfilling.
	//	Selector: func(profile profile.ModeProfile, match *pb.Match) *v1.GameServerAllocation {
	//		return selector.CommonPlayerBasedSelector(profile, match, int64(len(match.Tickets)))
	//	},
	//	MatchFunction: func(profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) {
	//		matches, err := functions.MakeInstantMatches(profile, tickets)
	//		return matches, nil, err
	//	},
	//},
	"lobby": {
		Name:            "lobby",
		GameName:        "game.lobby",
		FleetName:       "lobby",
		MatchmakingRate: 500 * time.Millisecond,
		MinPlayers:      1,
		MaxPlayers:      50,
		Selector: func(profile profile.ModeProfile, match *pb.Match) *v1.GameServerAllocation {
			return selector.CommonPlayerBasedSelector(profile, match, int64(len(match.Tickets)))
		},
		MatchFunction: func(profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) {
			matches, err := functions.MakeInstantMatches(profile, tickets)
			return matches, nil, err
		},
	},
	"block_sumo": {
		Name:            "block_sumo",
		GameName:        "game.block_sumo",
		FleetName:       "block-sumo",
		MatchmakingRate: 2 * time.Second,
		MinPlayers:      2,
		MaxPlayers:      12,
		Selector: func(profile profile.ModeProfile, match *pb.Match) *v1.GameServerAllocation {
			return selector.CommonSelector(profile, match)
		},
		MatchFunction: func(profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) {
			return functions.MakeCountdownMatches(profile, pendingMatches, tickets)
		},
	},
	"parkour_tag": {
		Name:            "parkour_tag",
		GameName:        "game.parkour_tag",
		FleetName:       "parkour-tag",
		MatchmakingRate: 2 * time.Second,
		MinPlayers:      2,
		MaxPlayers:      12,
		Selector: func(profile profile.ModeProfile, match *pb.Match) *v1.GameServerAllocation {
			return selector.CommonSelector(profile, match)
		},
		MatchFunction: func(profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) {
			return functions.MakeCountdownMatches(profile, pendingMatches, tickets)
		},
	},
	//"minesweeper": {
	//	Name:       "minesweeper",
	//	GameName:   "game.minesweeper",
	//	FleetName:  "minesweeper",
	//	MinPlayers: 1,
	//	MaxPlayers: 5,
	//	Selector: func(profile profile.ModeProfile, match *pb.Match) *v1.GameServerAllocation {
	//		return selector.CommonPlayerBasedSelector(profile, match, int64(len(match.Tickets)))
	//	},
	//	MatchFunction: func(profile profile.ModeProfile, pendingMatches []*pb.PendingMatch, tickets []*pb.Ticket) ([]*pb.Match, []*pb.PendingMatch, error) {
	//		matches, err := functions.MakeInstantMatches(profile, tickets)
	//		return matches, nil, err
	//	},
	//},
}
