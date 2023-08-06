package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type ZapLogger struct {
	Log *zap.Logger
}

func NewZapLogger(level zapcore.Level, encoding string, outputPath, errorOutputPath []string) (*ZapLogger, error) {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Encoding:         encoding,
		OutputPaths:      outputPath,
		ErrorOutputPaths: errorOutputPath,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "lvl",
			TimeKey:    "ts",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.DateTime))
			},
			EncodeLevel: func(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(lvl.String())
			},
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{Log: logger}, nil
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}
