package cmdproc

import (
	"strings"
	"time"

	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type CmdProcessor struct {
	stg    storage.Storage
	tz     *time.Location
	logger *zap.Logger
}

func NewCmdProcessor(stg storage.Storage, tz *time.Location, logger *zap.Logger) *CmdProcessor {
	return &CmdProcessor{stg: stg, tz: tz, logger: logger}
}

func (r *CmdProcessor) Process(c tele.Context, cmd string, userID int64) error {
	cmdParts := []string{}
	for _, part := range strings.Split(cmd, ",") {
		cmdParts = append(cmdParts, strings.Trim(part, " "))
	}

	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.String("reason", "empty command"),
			zap.String("command", cmd),
			zap.Int64("userid", userID),
		)
		return c.Send(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "h":
		resp = r.helpCommand(c)
	case "w":
		resp = r.processWeight(cmdParts[1:], userID)
	case "f":
		resp = r.processFood(cmdParts[1:], userID)
	case "j":
		resp = r.processJournal(cmdParts[1:], userID)
	case "cc":
		resp = r.calcCalCommand(cmdParts[1:])
	case "us":
		resp = r.processUserSettings(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid command",
			zap.String("reason", "unknown command"),
			zap.String("command", cmd),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	for _, rItem := range resp {
		if err := c.Send(rItem.what, rItem.opts...); err != nil {
			return err
		}
	}

	return nil
}

func (r *CmdProcessor) Stop() {
	if err := r.stg.Close(); err != nil {
		r.logger.Error("storage close error", zap.Error(err))
	}
}

type CmdResponse struct {
	what any
	opts []any
}

func NewCmdResponse(what any, opts ...any) CmdResponse {
	return CmdResponse{what: what, opts: opts}
}

func NewSingleCmdResponse(what any, opts ...any) []CmdResponse {
	return []CmdResponse{
		{what: what, opts: opts},
	}
}
