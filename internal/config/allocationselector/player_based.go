package allocationselector

import (
	"agones.dev/agones/pkg/apis"
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/pkg/pb"
	"math"
)

// CreatePlayerBasedSelector selects a GameServer where there is no 'match'.
// This could be a singleplayer game (e.g. marathon) or a stateless drop-in drop-out game (e.g. the lobby)
func CreatePlayerBasedSelector(cfg *liveconfig.GameModeConfig, match *pb.Match, playerCount int64) *allocatorv1.GameServerAllocation {
	fleetName := cfg.FleetName

	return &allocatorv1.GameServerAllocation{
		Spec: allocatorv1.GameServerAllocationSpec{
			Scheduling: apis.Packed,
			Selectors: []allocatorv1.GameServerSelector{
				{
					LabelSelector: v1.LabelSelector{
						MatchLabels: map[string]string{
							"agones.dev/fleet": fleetName,
						},
					},
					Players: &allocatorv1.PlayerSelector{
						MinAvailable: playerCount,
						MaxAvailable: math.MaxInt,
					},
					GameServerState: &AllocatedState,
				},
				{
					LabelSelector: v1.LabelSelector{
						MatchLabels: map[string]string{"agones.dev/fleet": fleetName},
					},
					GameServerState: &ReadyState,
				},
			},
			MetaPatch: allocatorv1.MetaPatch{
				Annotations: createPatchedAnnotations(match),
			},
		},
	}
}
