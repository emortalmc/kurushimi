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

	app := director.KurushimiApplication{
		StateStore: statestore.NewRedis(cfg.Redis),
		Config:     cfg,
		Messaging:  rmqMessenger,
		Logger:     logger,
	}
	notifier := notifier.New(messenger)

	frontend.Init(ctx, app)
	director.Init(ctx, notifier, app)

	select {
	case <-ctx.Done():
		return
	}
}
