package cmdproc

import (
	"github.com/devldavydov/myfood/internal/myfoodbot/helpdoc"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) helpCommand(c tele.Context) error {
	docRd, err := helpdoc.GetHelpDocument("help")
	if err != nil {
		r.logger.Error(
			"help command error",
			zap.Int64("userid", c.Sender().ID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}
	return c.Send(&tele.Document{
		File:     tele.FromReader(docRd),
		MIME:     "text/html",
		FileName: "help.html",
	})
}
