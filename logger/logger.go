package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

func init() {
	enc := zap.NewProductionEncoderConfig()
	enc.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	Log, _ = zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		EncoderConfig:     enc,
		DisableCaller:     true,
		DisableStacktrace: true,
		OutputPaths:       []string{"stderr"},
	}.Build()
}

func SyncLog() {
	Log.Sync()
}
