package app

import (
	"context"
	"fmt"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	"github.com/emortalmc/proto-specs/gen/go/grpc/party"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"kurushimi/internal/config"
	"kurushimi/internal/director"
	"kurushimi/internal/kafka"
	"kurushimi/internal/lobbycontroller"
	"kurushimi/internal/repository"
	"kurushimi/internal/service"
	"kurushimi/internal/utils/kubernetes"
	"kurushimi/pkg/pb"
	"net"
	"time"
)

func Run(ctx context.Context, cfg *config.Config, logger *zap.SugaredLogger) {
	logger.Info("starting kurushimi")
	// Parse gamemode configs
	gameModeController, err := liveconfig.NewGameModeConfigController(logger)
	if err != nil {
		logger.Fatalw("failed to create game mode config controller", err)
	}

	gameModes := gameModeController.GetConfigs()

	modeNames := make([]string, 0)
	for id := range gameModes {
		modeNames = append(modeNames, id)
	}
	logger.Infow("loaded initial gamemodes", "modeCount", len(gameModes), "modes", modeNames)

	_, agonesClient := kubernetes.CreateClients()

	repo, err := repository.NewMongoRepository(ctx, cfg.MongoDB)
	if err != nil {
		logger.Fatalw("failed to connect to mongo", err)
	}

	notifier := kafka.NewKafkaNotifier(cfg.Kafka, logger)

	err = repo.HealthCheck(ctx, 5*time.Second)
	if err != nil {
		logger.Fatalw("failed to initiate mongodb", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatalw("failed to listen", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(logger.Desugar(), grpc_zap.WithLevels(func(code codes.Code) zapcore.Level {
		if code != codes.Internal && code != codes.Unavailable && code != codes.Unknown {
			return zapcore.DebugLevel
		} else {
			return zapcore.ErrorLevel
		}
	}))))

	if cfg.Development {
		reflection.Register(s)
	}

	pConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.PartyService.ServiceHost, cfg.PartyService.ServicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw("failed to connect to party service", err)
	}
	partyService := party.NewPartyServiceClient(pConn)

	pSConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.PartyService.SettingsServiceHost, cfg.PartyService.SettingsServicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw("failed to connect to party service", err)
	}
	partySettingsService := party.NewPartySettingsServiceClient(pSConn)

	allocationClient := agonesClient.AllocationV1().GameServerAllocations(cfg.Namespace)

	// Lobby controller
	lobbyCtrl := lobbycontroller.NewLobbyController(logger, cfg, notifier, allocationClient)
	go lobbyCtrl.Run(ctx)

	pb.RegisterMatchmakerServer(s, service.NewMatchmakerService(logger, repo, notifier, gameModeController, lobbyCtrl,
		partyService, partySettingsService))

	logger.Infow("started kurushimi listener", "port", cfg.Port)

	go func() {
		err = s.Serve(lis)
		if err != nil {
			logger.Fatalw("failed to serve", err)
			return
		}
	}()

	directR := director.New(logger, repo, notifier, allocationClient, gameModeController)
	directR.Start(ctx)

	<-ctx.Done()
	logger.Info("shutting down kurushimi")
	s.Stop()
}
