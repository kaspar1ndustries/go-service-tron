package logger

import "go.uber.org/zap"

var (
	// log *zap.Logger
	sugar *zap.SugaredLogger
)

func init() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	sugar = log.Sugar()
}

func Info(format string) {
	sugar.Infow("Custom error",
		"error", format,
	)
}

func Error(format string) {
	sugar.Errorw("Custom error",
		"error", format,
	)
}
