package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	srv "github.com/devldavydov/myfood/internal/myfoodserver"
)

const (
	_defaultRunAddress      = "http://127.0.0.1:8080"
	_defaultShutdownTimeout = 15 * time.Second
	_defaultLogLevel        = "INFO"
	_defaultDBFilePath      = ""
)

type Config struct {
	RunAddress      string
	ShutdownTimeout time.Duration
	DBFilePath      string
	LogLevel        string
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	config := &Config{}

	flagSet.StringVar(&config.RunAddress, "a", _defaultRunAddress, "Server run address")
	flagSet.StringVar(&config.DBFilePath, "d", _defaultDBFilePath, "DB file path")
	flagSet.StringVar(&config.LogLevel, "l", _defaultLogLevel, "Log level")
	flagSet.DurationVar(&config.ShutdownTimeout, "t", _defaultShutdownTimeout, "Server shutdown timeout")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}

	err := flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	if config.DBFilePath == _defaultDBFilePath {
		return nil, fmt.Errorf("invalid DB file path")
	}

	return config, nil
}

func ServiceSettingsAdapt(config *Config) (*srv.ServerSettings, error) {
	return srv.NewServerSettings(
		config.RunAddress,
		config.DBFilePath,
		config.ShutdownTimeout)
}
