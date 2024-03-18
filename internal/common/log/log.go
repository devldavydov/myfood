package log

import "go.uber.org/zap"

func NewLogger(logLevel string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	lvl, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return nil, err
	}
	cfg.Level = lvl

	return cfg.Build()
}
