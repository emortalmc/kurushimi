package notifier

import (
	"context"
	"fmt"
	"github.com/emortalmc/grpc-api-specs/gen/go/service/player_tracker"
	"github.com/emortalmc/grpc-api-specs/gen/go/service/server_discovery"
	"github.com/emortalmc/grpc-api-specs/gen/go/service/velocity/player_transporter"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kurushimi/internal/utils/kubernetes"
	"kurushimi/pkg/pb"
	"os"
)

var (
	trackerEnabled = isTrackerEnabled()

	namespace = os.Getenv("NAMESPACE")

	kubeClient          = kubernetes.KubeClient
	playerTrackerClient = createPlayerTrackerClient()
)

func notifyTransport(ctx context.Context, match *pb.Match) error {
	if !trackerEnabled {
		return status.Error(codes.Unavailable, "Match notifications are not available")
	}

	// Get the player IDs ignoring tickets marked to not notify the proxy
	var players []string
	for _, ticket := range match.Tickets {
		if ticket.NotifyProxy != nil && ticket.GetNotifyProxy() == false {
			continue
		}
		players = append(players, ticket.GetPlayerId())
	}
	if len(players) == 0 {
		return nil
	}

	resp, err := playerTrackerClient.GetPlayerServers(ctx, &player_tracker.PlayersRequest{PlayerIds: players})
	if err != nil {
		return status.Errorf(codes.Internal, "Couldn't get player servers: %s", err)
	}

	// flip the map to map[string][]string (map[proxyId][]playerId)
	var serverPlayers = make(map[string][]string)
	for pId, server := range resp.PlayerServers {
		sId := server.ProxyId
		if serverPlayers[sId] == nil {
			serverPlayers[sId] = make([]string, 0)
		}
		serverPlayers[sId] = append(serverPlayers[sId], pId)
	}

	for sId, pIds := range serverPlayers {
		pod, err := kubeClient.CoreV1().Pods(namespace).Get(ctx, sId, v1.GetOptions{})
		if err != nil {
			return status.Errorf(codes.Internal, "Error retrieving Pod from Kubernetes API: %s", err)
		}
		ip := pod.Status.PodIP
		port, err := getGrpcPort(pod)

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
		if err != nil {
			return status.Errorf(codes.Internal, "Error communicating with Velocity gRPC: %s", err)
		}
		client := player_transporter.NewVelocityPlayerTransporterClient(conn)

		assignment := match.Assignment
		_, err = client.SendToServer(ctx, &player_transporter.TransportRequest{
			Server: &server_discovery.ConnectableServer{
				Id:      assignment.ServerId,
				Address: assignment.ServerAddress,
				Port:    assignment.ServerPort,
			},
			PlayerIds: pIds,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "Error communicating with Velocity gRPC: %s", err)
		}
	}

	return nil
}

func createPlayerTrackerClient() player_tracker.PlayerTrackerClient {
	conn, err := grpc.Dial("localhost:50502", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		logger.Error("Failed to connect to Player Tracker", zap.Error(err))
	}
	return player_tracker.NewPlayerTrackerClient(conn)
}

func getGrpcPort(pod *v12.Pod) (int32, error) {
	container := pod.Spec.Containers[0]
	for _, port := range container.Ports {
		if port.Name == "grpc" {
			return port.ContainerPort, nil
		}
	}
	return 0, fmt.Errorf("no grpc port found for server %s", pod.Name)
}

func isTrackerEnabled() bool {
	v := os.Getenv("TRACKER_ENABLED")
	return v == "true"
}
