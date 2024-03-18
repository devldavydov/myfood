package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	lvl, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return nil, err
	}
	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	return cfg.Build()
}
