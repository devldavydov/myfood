package myfoodbot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

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
	return s.cmdProc.process(c, c.Text(), c.Sender().ID)
}
