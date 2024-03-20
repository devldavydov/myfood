package cmdproc

import (
	"encoding/csv"
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
	logger *zap.Logger
}

func NewCmdProcessor(stg storage.Storage, logger *zap.Logger) *CmdProcessor {
	return &CmdProcessor{stg: stg, logger: logger}
}

func (r *CmdProcessor) Process(c tele.Context, cmd string, userID int64) error {
	rd := csv.NewReader(strings.NewReader(cmd))
	cmdParts, err := rd.Read()
	if err != nil {
		r.logger.Error("command parse error", zap.String("command", cmd), zap.Int64("userid", userID), zap.Error(err))
		return c.Send(msgErrInternal)
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

	switch cmdParts[0] {
	case "h":
		return r.helpCommand(c)
	case "w":
		return r.processWeight(c, cmdParts[1:], userID)
	}

	r.logger.Error(
		"invalid command",
		zap.String("reason", "unknown command"),
		zap.String("command", cmd),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}
