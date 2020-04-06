package logger

import (
	"github.com/astaxie/beego"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugVerbosity = "debugVerbosity"
)

func GetLogger() *zap.Logger {
	b, _ := beego.AppConfig.Bool(debugVerbosity)
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLevel(b)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	log, _ := cfg.Build()
	return log
}

func getLevel(debugVerbosity bool) zapcore.Level {
	if debugVerbosity {
		return zap.DebugLevel
	}
	return zap.InfoLevel
}
