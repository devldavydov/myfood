package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	bot "github.com/devldavydov/myfood/internal/myfoodbot"
)

const (
	_defaultToken       = ""
	_defaultPollTimeout = 10 * time.Second
)

type Config struct {
	Token       string
	PollTimeOut time.Duration
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	config := &Config{}

	flagSet.StringVar(&config.Token, "t", _defaultToken, "Telegram API token (required)")
	flagSet.DurationVar(&config.PollTimeOut, "p", _defaultPollTimeout, "Telegram API poll timeout")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}

	err := flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	if config.Token == _defaultToken {
		return nil, fmt.Errorf("invalid token")
	}

	return config, nil
}

func ServiceSettingsAdapt(config *Config, buildCommit string) (*bot.ServiceSettings, error) {
	return bot.NewServiceSettings(config.Token, config.PollTimeOut, buildCommit)
}
