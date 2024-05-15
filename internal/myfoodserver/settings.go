package myfoodserver

import (
	"net/url"
	"time"
)

type ServerSettings struct {
	RunAddress      *url.URL
	DBFilePath      string
	SessionSecret   string
	ShutdownTimeout time.Duration
}

func NewServerSettings(
	runAddress string,
	dbFilePath string,
	sessionSecret string,
	shutdownTimeout time.Duration) (*ServerSettings, error) {

	urlRunAddress, err := url.ParseRequestURI(runAddress)
	if err != nil {
		return nil, err
	}

	return &ServerSettings{
		RunAddress:      urlRunAddress,
		DBFilePath:      dbFilePath,
		SessionSecret:   sessionSecret,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}
