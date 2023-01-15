package main

import (
	"context"
	"go.uber.org/zap"
	"kurushimi/internal/config"
	"kurushimi/internal/config/dynamic"
	"kurushimi/internal/director"
	"kurushimi/internal/frontend"
	"kurushimi/internal/messaging"
	notifier2 "kurushimi/internal/notifier"
	"kurushimi/internal/statestore"
	"kurushimi/internal/utils/kubernetes"
	"os"
	"os/signal"
)

func main() {
	logger := createLogger()
	zap.ReplaceGlobals(logger.Desugar())
	defer logger.Sync()

	logger.Infow("Starting Kurushimi", "profileCount", len(config.ModeProfiles), "profiles", config.ModeProfiles)

	cfg, err := dynamic.ParseConfig()
	if err != nil {
		logger.Fatalw("Failed to parse config", err)
	}

	kubernetes.Init()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	messenger, err := messaging.NewRabbitMQ(cfg)
	if err != nil {
		logger.Fatalw("Failed to create messenger", err)
	}

	app := director.KurushimiApplication{
		StateStore: statestore.NewRedis(cfg),
		Config:     cfg,
		Messaging:  *messenger,
		Logger:     logger,
	}
	notifier := notifier2.New(messenger)

	frontend.Init(ctx, app)
	director.Init(ctx, notifier, app)

	select {
	case <-ctx.Done():
		return
	}
}

func createLogger() *zap.SugaredLogger {
	_, production := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	var err error
	var unsugared *zap.Logger
	if production {
		unsugared, err = zap.NewProduction()
	} else {
		unsugared, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	return unsugared.Sugar()
}
