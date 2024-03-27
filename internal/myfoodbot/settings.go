package myfoodbot

import "time"

type ServiceSettings struct {
	Token          string
	PollTimeOut    time.Duration
	BuildCommit    string
	DBFilePath     string
	AllowedUserIDs []int64
}

func NewServiceSettings(
	token string,
	pollTimeout time.Duration,
	dbFilePath string,
	alloweUserIDs []int64,
	buildVersion string) (*ServiceSettings, error) {
	return &ServiceSettings{
		Token:          token,
		PollTimeOut:    pollTimeout,
		BuildCommit:    buildVersion,
		DBFilePath:     dbFilePath,
		AllowedUserIDs: alloweUserIDs,
	}, nil
}
