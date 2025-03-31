package main

import "go.uber.org/zap"

func main() {
	//logger, _ := zap.NewDevelopment()
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, _ := cfg.Build()

	logger.Debug("hello debug")
	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")
}
