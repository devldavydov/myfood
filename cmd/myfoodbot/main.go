package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/devldavydov/myfoodbot/internal/myfoodbot"
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

	serviceSettings, err := ServiceSettingsAdapt(config, buildCommit)
	if err != nil {
		return fmt.Errorf("failed to create bot service settings: %w", err)
	}

	service, err := bot.NewService(*serviceSettings)
	if err != nil {
		return fmt.Errorf("failed to create bot service: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	return service.Run(ctx)
}
