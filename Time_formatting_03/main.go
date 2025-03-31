package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	//logger, _ := zap.NewDevelopment()
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logger, _ := cfg.Build()

	logger.Debug("hello debug")
	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")
}
