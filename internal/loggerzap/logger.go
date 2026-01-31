package loggerzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger() (*ZapLogger, error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	cfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.WarnLevel),
		Development:   false,
		DisableCaller: false,
		Encoding:      "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			CallerKey:      "source",
			MessageKey:     "message",
			LineEnding:     "\n",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"logs/app.log"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := zap.Must(cfg.Build())

	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) GetLogger() *zap.Logger {
	return l.logger
}

func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func (l *ZapLogger) Error(err string, field ...zap.Field) {
	l.logger.Error(err, field...)
}

func (l *ZapLogger) Info(msg string, field ...zap.Field) {
	l.logger.Info(msg, field...)
}
