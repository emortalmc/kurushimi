package app

import (
	"context"
	"go.uber.org/zap"
	"kurushimi/internal/config"
	"kurushimi/internal/director"
	"kurushimi/internal/frontend"
	"kurushimi/internal/notifier"
	"kurushimi/internal/rabbitmq"
	"kurushimi/internal/rabbitmq/messenger"
	"kurushimi/internal/statestore"
	"kurushimi/internal/utils/kubernetes"
)

func Run(ctx context.Context, cfg *config.Config, logger *zap.SugaredLogger) {
	logger.Infow("starting kurushimi", "profileCount", len(config.ModeProfiles), "profiles", config.ModeProfiles)

	kubernetes.Init()

	rabbitConn, err := rabbitmq.NewConnection(cfg.RabbitMq)
	if err != nil {
		logger.Fatalw("failed to connect to rabbitmq", err)
	}

	rmqMessenger, err := messenger.NewRabbitMQMessenger(logger, rabbitConn)
	if err != nil {
		logger.Fatalw("failed to create messenger", err)
	}

	stateStore := statestore.NewRedis(cfg.Redis)

	notif := notifier.NewNotifier(logger, rmqMessenger)

	frontend.Init(ctx, stateStore, notif)
	director.Init(ctx, logger, notif, stateStore, cfg.Namespace)

	select {
	case <-ctx.Done():
		return
	}
}
