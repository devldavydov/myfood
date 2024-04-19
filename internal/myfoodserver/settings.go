package myfoodserver

import (
	"net/url"
	"time"
)

type ServerSettings struct {
	RunAddress      *url.URL
	ShutdownTimeout time.Duration
}

func NewServerSettings(runAddress string, shutdownTimeout time.Duration) (*ServerSettings, error) {
	urlRunAddress, err := url.ParseRequestURI(runAddress)
	if err != nil {
		return nil, err
	}

	return &ServerSettings{
		RunAddress:      urlRunAddress,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}
