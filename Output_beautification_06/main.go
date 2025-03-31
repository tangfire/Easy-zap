package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 定义颜色
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

// 自定义 EncodeLevel
func coloredLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(colorBlue + "DEBUG" + colorReset)
	case zapcore.InfoLevel:
		enc.AppendString(colorGreen + "INFO" + colorReset)
	case zapcore.WarnLevel:
		enc.AppendString(colorYellow + "WARN" + colorReset)
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString(colorRed + "ERROR" + colorReset)
	default:
		enc.AppendString(level.String()) // 默认行为
	}
}
func dev() {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
	cfg.EncoderConfig.EncodeLevel = coloredLevelEncoder
	// 创建 logger
	logger, _ := cfg.Build()

	logger.Info("dev this is info")
	logger.Warn("dev this is warn")
	logger.Error("dev this is error")
}

func main() {
	dev()
	//cfg := zap.NewDevelopmentConfig()
	//cfg.EncoderConfig.EncodeLevel = coloredLevelEncoder
	//
	//logger, _ := cfg.Build()
	//logger.Debug("dev this is debug")
	//logger.Info("dev this is info")
	//logger.Warn("dev this is warn")
	//logger.Error("dev this is error")
}
