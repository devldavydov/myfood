package myfoodbot

import "time"

type ServiceSettings struct {
	Token       string
	PollTimeOut time.Duration
	BuildCommit string
	DBFilePath  string
}

func NewServiceSettings(token string, pollTimeout time.Duration, dbFilePath, buildVersion string) (*ServiceSettings, error) {
	return &ServiceSettings{
		Token:       token,
		PollTimeOut: pollTimeout,
		BuildCommit: buildVersion,
		DBFilePath:  dbFilePath,
	}, nil
}
