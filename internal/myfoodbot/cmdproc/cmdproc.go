package cmdproc

import (
	"strings"
	"time"

	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

var (
	_stgOperationTimeout = 10 * time.Second
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
		return c.Send(msgErrInvalidCommand)
	}

	var what any
	var opts []any

	switch cmdParts[0] {
	case "h":
		what, opts = r.helpCommand(c)
	case "w":
		what, opts = r.processWeight(cmdParts[1:], userID)
	case "f":
		what, opts = r.processFood(cmdParts[1:], userID)
	case "j":
		what, opts = r.processJournal(cmdParts[1:], userID)
	case "cc":
		what, opts = r.calcCalCommand(cmdParts[1:])
	case "us":
		what, opts = r.processUserSettings(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid command",
			zap.String("reason", "unknown command"),
			zap.String("command", cmd),
			zap.Int64("userid", userID),
		)
		what = msgErrInvalidCommand
	}

	return c.Send(what, opts...)
}

func (r *CmdProcessor) Stop() {
	if err := r.stg.Close(); err != nil {
		r.logger.Error("storage close error", zap.Error(err))
	}
}
