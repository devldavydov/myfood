package myfoodbot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (s *Service) onStart(c tele.Context) error {
	return c.Send(
		fmt.Sprintf("Hello, %s!\nWelcome to MyFoodBot!\nBuild: %s", c.Sender().Username, s.settings.BuildCommit),
	)
}

func (s *Service) onText(c tele.Context) error {
	cmd := c.Text()
	return c.Send(fmt.Sprintf("Your command: %s", cmd))
}
