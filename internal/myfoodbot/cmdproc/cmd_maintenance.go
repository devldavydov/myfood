package cmdproc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processMaintenance(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid maintenance command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "backup":
		resp = r.backupCommand(userID)
	default:
		r.logger.Error(
			"invalid maintenance command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) backupCommand(userID int64) []CmdResponse {
	// Get backup from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*10)
	defer cancel()

	backup, err := r.stg.Backup(ctx)
	if err != nil {
		r.logger.Error(
			"backup command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Generate response.
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if err := json.NewEncoder(zw).Encode(&backup); err != nil {
		r.logger.Error(
			"backup json err",
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	if err := zw.Close(); err != nil {
		r.logger.Error(
			"gzip err",
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(&buf),
		MIME:     "application/x-gzip-compressed",
		FileName: fmt.Sprintf("backup_%s.json.gz", formatTimestamp(time.Now().In(r.tz))),
	})
}
