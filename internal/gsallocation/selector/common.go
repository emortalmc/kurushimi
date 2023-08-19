package selector

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocatorv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"encoding/base64"
	"github.com/emortalmc/kurushimi/pkg/pb"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	AllocatedState = agonesv1.GameServerStateAllocated
	ReadyState     = agonesv1.GameServerStateReady

	notOudatedExpression = v1.LabelSelectorRequirement{
		Key:      "emortal.dev/outdated",
		Operator: v1.LabelSelectorOpDoesNotExist,
	}
)

func createReadySelector(fleetName string) allocatorv1.GameServerSelector {
	return allocatorv1.GameServerSelector{
		LabelSelector: v1.LabelSelector{
			MatchLabels: map[string]string{
				"agones.dev/fleet": fleetName,
			},

			MatchExpressions: []v1.LabelSelectorRequirement{notOudatedExpression},
		},

		GameServerState: &ReadyState,
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
