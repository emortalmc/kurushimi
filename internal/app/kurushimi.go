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
	"github.com/emortalmc/proto-specs/gen/go/grpc/party"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(cfg *config.Config, logger *zap.SugaredLogger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	wg := &sync.WaitGroup{}

	repoCtx, repoCancel := context.WithCancel(ctx)
	repoWg := &sync.WaitGroup{}

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

	repo, err := repository.NewMongoRepository(repoCtx, repoWg, logger, &cfg.MongoDB)
	if err != nil {
		logger.Fatalw("failed to connect to mongo", err)
	}

	notifier := kafka.NewKafkaNotifier(ctx, wg, &cfg.Kafka, logger)

	kafka.NewConsumer(ctx, wg, &cfg.Kafka, logger, repo)

	err = repo.HealthCheck(ctx, 5*time.Second)
	if err != nil {
		logger.Fatalw("failed to initiate mongodb", err)
	}

	pSConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.PartyService.SettingsServiceHost, cfg.PartyService.SettingsServicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw("failed to connect to party service", err)
	}
	partySettingsService := party.NewPartySettingsServiceClient(pSConn)

	allocationClient := agonesClient.AllocationV1().GameServerAllocations(cfg.Namespace)

	pConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.PartyService.ServiceHost, cfg.PartyService.ServicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw("failed to connect to party service", err)
	}

	partyService := party.NewPartyServiceClient(pConn)

	// Lobby controller
	lobbyCtrl := lobbycontroller.NewLobbyController(ctx, wg, logger, cfg, notifier, allocationClient)

	service.RunServices(ctx, logger, wg, cfg, repo, notifier, gameModeController, lobbyCtrl, partyService, partySettingsService)

	directR := director.New(logger, repo, notifier, allocationClient, gameModeController)
	directR.Start(ctx)

	wg.Wait()

	logger.Info("shutting down repository")
	repoCancel()
	repoWg.Wait()
}
