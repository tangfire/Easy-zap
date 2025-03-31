package main

import "go.uber.org/zap"

func dev() {
	logger, _ := zap.NewDevelopment()
	logger.Debug("hello debug")
	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")
	logger.Panic("hello panic")
	logger.Fatal("hello fatal")

}

func example() {
	logger := zap.NewExample()
	logger.Debug("hello debug")
	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")
	logger.Panic("hello panic")
	logger.Fatal("hello fatal")
}

func prod() {
	logger, _ := zap.NewProduction()
	logger.Debug("hello debug")
	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")
	logger.Panic("hello panic")
	logger.Fatal("hello fatal")
}

func main() {
	dev()
	//example()
	//prod()
}
