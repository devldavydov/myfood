package cmdproc

import (
	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/myfoodbot/helpdoc"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) helpCommand(c tele.Context) []CmdResponse {
	docRd, err := helpdoc.GetHelpDocument("help")
	if err != nil {
		r.logger.Error(
			"help command error",
			zap.Int64("userid", c.Sender().ID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(docRd),
		MIME:     "text/html",
		FileName: "help.html",
	})
}
