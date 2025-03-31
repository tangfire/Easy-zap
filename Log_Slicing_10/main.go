package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

// 自定义日志写入器
type dynamicLogWriter struct {
	mu         sync.Mutex
	currentDay string
	file       *os.File
	logDir     string
}

func (w *dynamicLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 检查是否需要切换到新的日志文件
	currentDay := time.Now().Format("2006-01-02")
	if currentDay != w.currentDay {
		// 关闭当前日志文件
		if w.file != nil {
			w.file.Close()
		}

		// 创建新的日志文件
		if err := os.MkdirAll(w.logDir, 0755); err != nil {
			return 0, err
		}
		filePath := w.logDir + "/app-" + currentDay + ".log"
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.file = file
		w.currentDay = currentDay
	}

	// 写入日志
	return w.file.Write(p)
}

// 初始化全局日志
func initLogger() {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
	// 创建 Logger
	writer := &dynamicLogWriter{
		logDir: "logs",
	}
	// 创建 Core
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(os.Stdout), // 输出到控制台
		zapcore.DebugLevel,         // 设置日志级别
	)
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer), // 输出到文件
		zapcore.DebugLevel,      // 设置日志级别
	)
	core := zapcore.NewTee(consoleCore, fileCore)
	// 创建 Logger
	logger := zap.New(core, zap.AddCaller())
	for i := 0; i < 10; i++ {
		logger.Sugar().Infof("this is %d log", i)
		time.Sleep(1 * time.Second)
	}

	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()
}
