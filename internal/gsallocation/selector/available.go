package selector

import (
	"agones.dev/agones/pkg/apis"
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	pb "github.com/emortalmc/proto-specs/gen/go/model/matchmaker"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
				createReadySelector(fleetName),
			},
		},
	}
}
