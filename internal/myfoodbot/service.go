package myfoodbot

import (
	"context"
	"fmt"

	"github.com/devldavydov/myfood/internal/myfoodbot/cmdproc"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Service struct {
	settings ServiceSettings
	cmdProc  *cmdproc.CmdProcessor
	logger   *zap.Logger
}

func NewService(settings ServiceSettings, logger *zap.Logger) (*Service, error) {
	stg, err := storage.NewStorageSQLite(settings.DBFilePath, logger)
	if err != nil {
		return nil, err
	}

	return &Service{settings: settings, cmdProc: cmdproc.NewCmdProcessor(stg, logger)}, nil
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

func (s *Service) onStart(c tele.Context) error {
	return c.Send(
		fmt.Sprintf(
			"Привет, %s!\nДобро пожаловать в MyFoodBot!\nСборка: %s\nОтправь 'h' для помощи",
			c.Sender().Username,
			s.settings.BuildCommit,
		),
	)
}

func (s *Service) onText(c tele.Context) error {
	return s.cmdProc.Process(c, c.Text(), c.Sender().ID)
}
