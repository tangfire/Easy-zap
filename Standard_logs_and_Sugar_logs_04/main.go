package main

import (
	"fmt"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	logger.Info(fmt.Sprintf("my name is %s", "tangfire"))

	sl := logger.Sugar()
	sl.Infof("my name is %s", "tangfire")
}
