package myfoodbot

import "time"

type ServiceSettings struct {
	Token       string
	PollTimeOut time.Duration
	BuildCommit string
}

func NewServiceSettings(token string, pollTimeout time.Duration, buildVersion string) (*ServiceSettings, error) {
	return &ServiceSettings{
		Token:       token,
		PollTimeOut: pollTimeout,
		BuildCommit: buildVersion,
	}, nil
}
