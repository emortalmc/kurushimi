package allocationselector

import (
	"agones.dev/agones/pkg/apis"
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"encoding/base64"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/pkg/pb"
	"math"
)

var (
	AllocatedState = agonesv1.GameServerStateAllocated
	ReadyState     = agonesv1.GameServerStateReady
)

// contains some common selectors

func CreateAvailableSelector(cfg *liveconfig.GameModeConfig, match *pb.Match) *allocatorv1.GameServerAllocation {
	fleetName := cfg.FleetName

	return &allocatorv1.GameServerAllocation{
		Spec: allocatorv1.GameServerAllocationSpec{
			Scheduling: apis.Packed,
			Selectors: []allocatorv1.GameServerSelector{
				{
					LabelSelector: v1.LabelSelector{
						MatchLabels: map[string]string{
							"agones.dev/fleet":               fleetName,
							"agones.dev/sdk-should-allocate": "true",
						},
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

func createPatchedAnnotations(match *pb.Match) map[string]string {
	allocationData := &pb.AllocationData{Match: match}
	allocationDataBytes, err := proto.Marshal(allocationData)
	if err != nil {
		panic(err)
	}

	stringData := base64.StdEncoding.EncodeToString(allocationDataBytes)

	return map[string]string{
		"emortal.dev/allocation-data": stringData,
	}
}
