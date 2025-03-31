package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
)

// 定义前缀
const logPrefix = "[MyApp] "

// 自定义 Encoder
type prefixedEncoder struct {
	zapcore.Encoder
}

func (e *prefixedEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 先调用原始的 EncodeEntry 方法生成日志行
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// 在日志行的最前面添加前缀
	logLine := buf.String()
	buf.Reset()
	buf.AppendString(logPrefix + logLine)

	return buf, nil
}
func dev() {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
	// 创建自定义的 Encoder
	encoder := &prefixedEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器
	}
	// 创建 Core
	core := zapcore.NewCore(
		encoder,                    // 使用自定义的 Encoder
		zapcore.AddSync(os.Stdout), // 输出到控制台
		zapcore.DebugLevel,         // 设置日志级别
	)

	// 创建 Logger
	logger := zap.New(core, zap.AddCaller())

	logger.Info("dev this is info")
	logger.Warn("dev this is warn")
	logger.Error("dev this is error")
}

func main() {
	dev()
}
