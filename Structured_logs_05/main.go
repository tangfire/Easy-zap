package main

import "go.uber.org/zap"

func main() {
	logger, _ := zap.NewDevelopment()
	logger.Info("this is log")
	logger.Info("this is log", zap.String("name", "tangfire"), zap.Int("age", 24), zap.Bool("isok", true))

	logger1, _ := zap.NewProduction()
	logger1.Info("this is log")
	logger1.Info("this is log", zap.String("name", "tangfire"), zap.Int("age", 24), zap.Bool("isok", true))
}
