package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 初始化全局日志
func initLogger() {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
	// 创建 Logger
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func dev() {
	zap.L().Info("dev this is info")
	zap.L().Warn("dev this is warn")
	zap.L().Error("dev this is error")
	zap.S().Infof("dev this is info %s", "xxx")
	zap.S().Warnf("dev this is warn %s", "xxx")
	zap.S().Errorf("dev this is error %s", "xxx")
}

func main() {
	initLogger()
	dev()
}
