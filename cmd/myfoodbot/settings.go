package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	bot "github.com/devldavydov/myfood/internal/myfoodbot"
)

const (
	_defaultToken       = ""
	_defaultPollTimeout = 10 * time.Second
	_defaultDBFilePath  = ""
	_defaultLogLevel    = "INFO"
)

type IDList []int64

func (r *IDList) String() string {
	return ""
}

func (r *IDList) Set(v string) error {
	iv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}
	*r = append(*r, iv)
	return nil
}

type Config struct {
	Token          string
	PollTimeOut    time.Duration
	DBFilePath     string
	LogLevel       string
	AllowedUserIDs IDList
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	config := &Config{}

	flagSet.StringVar(&config.Token, "t", _defaultToken, "Telegram API token (required)")
	flagSet.StringVar(&config.DBFilePath, "d", _defaultDBFilePath, "DB file path")
	flagSet.StringVar(&config.LogLevel, "l", _defaultLogLevel, "Log level")
	flagSet.DurationVar(&config.PollTimeOut, "p", _defaultPollTimeout, "Telegram API poll timeout")
	flagSet.Var(&config.AllowedUserIDs, "u", "Allowed User ID")

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

	if config.DBFilePath == _defaultDBFilePath {
		return nil, fmt.Errorf("invalid DB file path")
	}

	return config, nil
}

func ServiceSettingsAdapt(config *Config, buildCommit string) (*bot.ServiceSettings, error) {
	return bot.NewServiceSettings(
		config.Token,
		config.PollTimeOut,
		config.DBFilePath,
		config.AllowedUserIDs,
		buildCommit)
}
