package gsallocation

import (
	v12 "agones.dev/agones/pkg/apis/allocation/v1"
	v1 "agones.dev/agones/pkg/client/clientset/versioned/typed/allocation/v1"
	"context"
	"fmt"
	"github.com/emortalmc/kurushimi/internal/utils"
	pb "github.com/emortalmc/proto-specs/gen/go/model/matchmaker"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"sync"
)

func AllocateServers(ctx context.Context, allocationClient v1.GameServerAllocationInterface,
	allocations map[*pb.Match]*v12.GameServerAllocation) map[*pb.Match]error {

	errors := make(map[*pb.Match]error)
	wg := sync.WaitGroup{}
	wg.Add(len(allocations))

	for match, allocation := range allocations {
		go func(fMatch *pb.Match, fAllocation *v12.GameServerAllocation) {
			defer wg.Done()
			if err := AllocateServer(ctx, allocationClient, fMatch, fAllocation); err != nil {
				errors[fMatch] = err
			}
		}(match, allocation)
	}

	wg.Wait()
	return errors
}

func AllocateServer(ctx context.Context, allocationClient v1.GameServerAllocationInterface,
	match *pb.Match, allocationReq *v12.GameServerAllocation) error {

	resp, err := allocationClient.Create(ctx, allocationReq, v13.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create allocation: %w", err)
	}

	allocation := resp.Status
	if allocation.State != v12.GameServerAllocationAllocated {
		return fmt.Errorf("allocation was not successful: %s", allocation.State)
	}

	protocolVersion, versionName := parseVersions(resp.Annotations)

	match.Assignment = &pb.Assignment{
		ServerId:        allocation.GameServerName,
		ServerAddress:   allocation.Address,
		ServerPort:      uint32(allocation.Ports[0].Port),
		ProtocolVersion: protocolVersion,
		VersionName:     versionName,
	}
	return err
}

func parseVersions(annotations map[string]string) (protocolVersion *int64, versionName *string) {
	protocolStr, ok := annotations["agones.dev/emc-protocol-version"]
	if ok {
		protocol, err := strconv.Atoi(protocolStr)
		if err == nil {
			protocolVersion = utils.PointerOf(int64(protocol))
		}
	}

	version, ok := annotations["agones.dev/emc-version-name"]
	if ok {
		versionName = &version
	}

	return protocolVersion, versionName
}
