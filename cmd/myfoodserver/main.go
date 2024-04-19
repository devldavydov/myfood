package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_log "github.com/devldavydov/myfood/internal/common/log"
	srv "github.com/devldavydov/myfood/internal/myfoodserver"
	"go.uber.org/zap"
)

var (
	buildDate   string
	buildCommit string
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to load configuration settings: %w", err)
	}

	logger, err := _log.NewLogger(config.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	logger.Info("Started MyFood Server", zap.String("buildDate", buildDate), zap.String("buildCommit", buildCommit))

	serviceSettings, err := ServiceSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create bot service settings: %w", err)
	}

	service, err := srv.NewService(serviceSettings, logger)
	if err != nil {
		return fmt.Errorf("failed to create server service: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	return service.Run(ctx)
}
