package myfoodbot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (s *Service) onStart(c tele.Context) error {
	return c.Send(
		fmt.Sprintf(
			"Hello, %s!\nWelcome to MyFoodBot!\nBuild: %s\n\nEnter 'h' for help",
			c.Sender().Username,
			s.settings.BuildCommit,
		),
	)
}

func (s *Service) onText(c tele.Context) error {
	return s.cmdProc.process(c, c.Text(), c.Sender().ID)
}
