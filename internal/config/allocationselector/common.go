package allocationselector

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"kurushimi/pkg/pb"
)

var (
	AllocatedState = agonesv1.GameServerStateAllocated
	ReadyState     = agonesv1.GameServerStateReady
)

// contains some common selectors

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
