package app

import (
	"context"
	"fmt"
	"github.com/emortalmc/kurushimi/internal/config"
	"github.com/emortalmc/kurushimi/internal/director"
	"github.com/emortalmc/kurushimi/internal/kafka"
	"github.com/emortalmc/kurushimi/internal/lobbycontroller"
	"github.com/emortalmc/kurushimi/internal/repository"
	"github.com/emortalmc/kurushimi/internal/service"
	"github.com/emortalmc/kurushimi/internal/utils/kubernetes"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	"github.com/emortalmc/proto-specs/gen/go/grpc/matchmaker"
	"github.com/emortalmc/proto-specs/gen/go/grpc/party"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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

	repo, err := repository.NewMongoRepository(ctx, logger, cfg.MongoDB)
	if err != nil {
		logger.Fatalw("failed to connect to mongo", err)
	}

	notifier := kafka.NewKafkaNotifier(cfg.Kafka, logger)

	kafka.NewConsumer(ctx, cfg.Kafka, logger, repo)

	err = repo.HealthCheck(ctx, 5*time.Second)
	if err != nil {
		logger.Fatalw("failed to initiate mongodb", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatalw("failed to listen", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(grpczap.UnaryServerInterceptor(logger.Desugar(), grpczap.WithLevels(func(code codes.Code) zapcore.Level {
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

	matchmaker.RegisterMatchmakerServer(s, service.NewMatchmakerService(logger, repo, notifier, gameModeController, lobbyCtrl,
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
