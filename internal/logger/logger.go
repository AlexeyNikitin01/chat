package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"task4/internal/config"
)

var logger *zap.Logger

func zapConfig(cfg *config.Config) (*zap.Config, error) {
	c := zap.Config{
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	switch cfg.LogLevel {
	case "info":
		c.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "error":
		c.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}
	return &c, nil
}

func Init(cfg *config.Config) {
	c, err := zapConfig(cfg)
	if err != nil {
		panic(err)
	}
	
	logger, err = c.Build()
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
}

func Instance() *zap.Logger {
	return logger
}
