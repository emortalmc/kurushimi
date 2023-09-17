package service

import (
	"context"
	"fmt"
	"github.com/emortalmc/kurushimi/internal/config"
	"github.com/emortalmc/kurushimi/internal/kafka"
	"github.com/emortalmc/kurushimi/internal/lobbycontroller"
	"github.com/emortalmc/kurushimi/internal/repository"
	"github.com/emortalmc/live-config-parser/golang/pkg/liveconfig"
	"github.com/emortalmc/proto-specs/gen/go/grpc/matchmaker"
	"github.com/emortalmc/proto-specs/gen/go/grpc/party"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

func RunServices(ctx context.Context, logger *zap.SugaredLogger, wg *sync.WaitGroup, cfg *config.Config,
	repo repository.Repository, notifier kafka.Notifier, gameModeController liveconfig.GameModeConfigController,
	lobbyCtrl lobbycontroller.LobbyController, partyService party.PartyServiceClient,
	partySettingsService party.PartySettingsServiceClient) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatalw("failed to listen", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpczap.UnaryServerInterceptor(logger.Desugar(), grpczap.WithLevels(func(code codes.Code) zapcore.Level {
			if code != codes.Internal && code != codes.Unavailable && code != codes.Unknown {
				return zapcore.DebugLevel
			} else {
				return zapcore.ErrorLevel
			}
		})),
	))

	if cfg.Development {
		reflection.Register(s)
	}

	matchmaker.RegisterMatchmakerServer(s, newMatchmakerService(logger, repo, notifier, gameModeController, lobbyCtrl,
		partyService, partySettingsService))
	logger.Infow("listening for gRPC requests", "port", cfg.Port)

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatalw("failed to serve", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.GracefulStop()
	}()
}
