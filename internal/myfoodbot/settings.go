package myfoodbot

import "time"

type ServiceSettings struct {
	Token          string
	PollTimeOut    time.Duration
	BuildCommit    string
	DBFilePath     string
	AllowedUserIDs []int64
	TZ             *time.Location
}

func NewServiceSettings(
	token string,
	pollTimeout time.Duration,
	dbFilePath string,
	alloweUserIDs []int64,
	stz string,
	buildVersion string) (*ServiceSettings, error) {

	tz, err := time.LoadLocation(stz)
	if err != nil {
		return nil, err
	}

	return &ServiceSettings{
		Token:          token,
		PollTimeOut:    pollTimeout,
		BuildCommit:    buildVersion,
		DBFilePath:     dbFilePath,
		AllowedUserIDs: alloweUserIDs,
		TZ:             tz,
	}, nil
}
