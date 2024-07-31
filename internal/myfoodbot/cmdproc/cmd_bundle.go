package cmdproc

import (
	"github.com/devldavydov/myfood/internal/common/messages"
	"go.uber.org/zap"
)

func (r *CmdProcessor) processBundle(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid bundle command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.bundleSetCommand(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid bundle command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) bundleSetCommand(cmdParts []string, userID int64) []CmdResponse {
	return nil
}
