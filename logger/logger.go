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

func Info(msg string) {
	sugar.Infof(msg)
}

func Error(err error) {
	sugar.Error(err)
}

func Errorw(msg string, kv interface{}) {
	sugar.Errorw(msg, kv)
}

func Fatal(msg string) {
	sugar.Fatal(msg)
}
