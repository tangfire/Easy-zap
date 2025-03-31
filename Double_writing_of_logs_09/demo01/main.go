package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 初始化全局日志
func initLogger() {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
	// 创建 Core
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(os.Stdout), // 输出到控制台
		zapcore.DebugLevel,         // 设置日志级别
	)
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(file), // 输出到文件
		zapcore.DebugLevel,    // 设置日志级别
	)
	core := zapcore.NewTee(consoleCore, fileCore)
	// 创建 Logger
	logger := zap.New(core, zap.AddCaller())
	logger.Info("hello world")

	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()
}
