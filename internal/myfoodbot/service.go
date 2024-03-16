package myfoodbot

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

type Service struct {
	settings ServiceSettings
}

func NewService(settings ServiceSettings) (*Service, error) {
	return &Service{settings: settings}, nil
}

func (s *Service) Run(ctx context.Context) error {
	pref := tele.Settings{
		Token:  s.settings.Token,
		Poller: &tele.LongPoller{Timeout: s.settings.PollTimeOut},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}

	s.setupRouting(b)
	go b.Start()

	select {
	case <-ctx.Done():
		b.Stop()
	}

	return nil
}

func (s *Service) setupRouting(b *tele.Bot) {
	b.Handle("/start", s.onStart)
	b.Handle(tele.OnText, s.onText)
}
