安装

```bash
go get -u go.uber.org/zap
```

# 基本使用

`zap`库的使用与其他的日志库非常相似。

先创建一个`logger`，然后调用各个级别的方法记录日志（`Debug/Info/Error/Warn`）。




```go
func dev() {
  logger, _ := zap.NewDevelopment()
  logger.Info("dev this is info")
  logger.Warn("dev this is warn")
  logger.Error("dev this is error")
}

func test() {
  logger := zap.NewExample()
  logger.Info("exam this is info")
  logger.Warn("exam this is warn")
  logger.Error("exam this is error")
}

func prod() {
  logger, _ := zap.NewProduction()
  logger.Info("prod this is info")
  logger.Warn("prod this is warn")
  logger.Error("prod this is error")
}
```

dev模式下，日志格式是text格式，并且warn和error会有栈信息

example模式下，格式是json，并且字段只有level和msg

prod模式下，格式也是json，多一个时间和函数位置字段，生产环境上用json格式的日志更方便排查

# 设置日志级别

```go
// 使用 zap 的 NewDevelopmentConfig 快速配置
cfg := zap.NewDevelopmentConfig()
cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
// 创建 logger
logger, _ := cfg.Build()
logger.Debug("this is dev debug log")
logger.Info("this is dev info log")
logger.Warn("this is dev warn log")
logger.Error("this is dev error log")
logger.Fatal("this is dev fatal log")
```

# 时间格式化

默认的时间要么是带了时区的，要么就是时间戳，不太美观

如果只是想把时间变成我们喜欢的，可以使用如下配置

```go
// 使用 zap 的 NewDevelopmentConfig 快速配置
cfg := zap.NewDevelopmentConfig()
cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式

// 创建 logger
logger, _ := cfg.Build()

logger.Info("dev this is info")
logger.Warn("dev this is warn")
logger.Error("dev this is error")
```


# 标准日志和Sugar日志

标准的`*zap.Logger`只有这几个方法

如果想输出格式化字符串，还得在里面套函数，比较麻烦

所以可以用Sugar方法得到一个加强版实例

这里面方法就比较多了，常用的就是写格式化字符串哪些方法

# 结构化日志

`zap`支持通过`Field`的形式记录结构化日志，方便分析和查询

```go
logger, _ := zap.NewDevelopment()
logger.Info("this is info",
  zap.String("username", "admin"),
  zap.Int("user_id", 42),
  zap.Bool("active", true),
)
```


# 输出美化


info，warn，error显示不同的颜色，看起来好看些

变色的关键 颜色控制字符

```go
fmt.Printf("\033[31mthis is 红色\n\033[0m")
  fmt.Printf("\033[32mthis is 绿色\n\033[0m")
  fmt.Printf("\033[33mthis is 黄色\n\033[0m")
  fmt.Printf("\033[34mthis is 蓝色\n\033[0m")
  fmt.Printf("\033[35mthis is 紫色\n\033[0m")
  fmt.Printf("\033[36mthis is 青色\n\033[0m")
```


```go
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
```

# 加日志前缀

一般如果有多个项目日志，可能会在日志的前面加上项目的名称


```go
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
```

# 全局日志

因为zap推崇的还是以对象的形式使用日志

但是有些时候想要在应用程序的任何地方都可以直接使用的日志实例，那么可以用到全局日志

```go
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
```

L方法返回的是标准zap实例，S方法返回的是superZap的实例，superZap主要多了模板字符串方法

# 日志双写

常见的：控制台和日志文件双写

使用zapcore.NewTee可以组合多个core实例

```go
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
  zap.ReplaceGlobals(logger)
}
```

或者

使用 `zapcore.NewMultiWriteSyncer`

```go
// 初始化全局日志
func initLogger() {
  // 使用 zap 的 NewDevelopmentConfig 快速配置
  cfg := zap.NewDevelopmentConfig()
  cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
  file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
  writeSyncer := zapcore.NewMultiWriteSyncer(
    zapcore.AddSync(os.Stdout),
    zapcore.AddSync(file),
  )
  // 创建 Core
  core := zapcore.NewCore(
    zapcore.NewConsoleEncoder(cfg.EncoderConfig),
    zapcore.AddSync(writeSyncer), // 输出到控制台
    zapcore.DebugLevel,         // 设置日志级别
  )
  // 创建 Logger
  logger := zap.New(core, zap.AddCaller())
  zap.ReplaceGlobals(logger)
}
```

# 日志切片

一般情况下，按照时间分片，每天的日志放到一个日志文件里面去

需要自定义Write方法

```go
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
  zap.ReplaceGlobals(logger)
}
```

# 日志按照level分片

把error的日志单独分出来，放到一个文件里面去

```go
// 自定义 Encoder
type levelEncoder struct {
  zapcore.Encoder
  errFile *os.File
}

func (e *levelEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
  // 先调用原始的 EncodeEntry 方法生成日志行
  buf, err := e.Encoder.EncodeEntry(entry, fields)
  if err != nil {
    return nil, err
  }
  switch entry.Level {
  case zapcore.ErrorLevel:
    if e.errFile == nil {
      file, _ := os.OpenFile("err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
      e.errFile = file
    }
    e.errFile.WriteString(buf.String())

  }

  return buf, nil
}
func main() {
  // 使用 zap 的 NewDevelopmentConfig 快速配置
  cfg := zap.NewDevelopmentConfig()
  cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
  // 创建自定义的 Encoder
  encoder := &levelEncoder{
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

  logger.Info("this is info")
  logger.Warn("this is warn")
  logger.Error("this is error1")
  logger.Error("this is error2")
}
```


以下是对这段代码的详细解释，结合了Zap日志库的自定义编码器实现原理和关键功能：

---

### 一、代码结构解析
#### 1. 自定义编码器 `levelEncoder`
```go
type levelEncoder struct {
    zapcore.Encoder       // 内嵌标准编码器
    errFile     *os.File  // 错误日志文件句柄
}
```
- **功能**：扩展Zap的编码器，实现错误日志分离。
- **继承与扩展**：内嵌`zapcore.Encoder`继承基础功能，新增`errFile`字段存储错误日志文件对象。

#### 2. 核心方法 `EncodeEntry`
```go
func (e *levelEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
    buf, err := e.Encoder.EncodeEntry(entry, fields) // 调用原始方法生成日志
    if entry.Level == zapcore.ErrorLevel {
        if e.errFile == nil { // 首次触发错误日志时创建文件
            file, _ := os.OpenFile("err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
            e.errFile = file
        }
        e.errFile.WriteString(buf.String()) // 写入错误日志文件
    }
    return buf, nil
}
```
- **日志分离逻辑**：仅对`Error`及以上级别日志写入`err.log`文件，其他日志仍输出到控制台。
- **延迟初始化**：首次遇到错误日志时创建文件，避免空文件占用资源。

#### 3. 日志初始化配置
```go
cfg := zap.NewDevelopmentConfig()
cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 时间格式
encoder := &levelEncoder{Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig)}
core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
logger := zap.New(core, zap.AddCaller())
```
- **开发环境配置**：`NewDevelopmentConfig`默认启用彩色输出和调用者信息。
- **时间格式化**：通过`EncodeTime`自定义时间格式为"年-月-日 时:分:秒"。

---

### 二、功能特性
#### 1. 控制台与文件双输出
- **控制台输出**：所有日志（包括错误）通过`os.Stdout`显示。
- **文件输出**：错误日志额外写入`err.log`，便于后续错误追踪。

#### 2. 日志级别处理
- **级别判定**：通过`entry.Level`判断是否为`ErrorLevel`，实现级别分片。
- **支持级别范围**：`DebugLevel`作为最低级别，允许记录`Debug`到`Fatal`的所有日志。

---

### 三、潜在问题与改进建议
#### 1. 并发安全问题
- **竞态条件**：多协程同时写入文件可能导致数据混乱。**改进方案**：使用`sync.Mutex`保护文件写入操作。
- **文件句柄泄漏**：未在程序退出时关闭文件。**修复方案**：添加`defer e.errFile.Close()`。

#### 2. 功能扩展建议
- **日志切割**：集成`lumberjack`库实现按大小/时间切割文件（参考网页7）：
  ```go
  import "gopkg.in/natefinch/lumberjack.v2"
  // 替换 OpenFile 逻辑为：
  lumberjackLogger := &lumberjack.Logger{
      Filename: "err.log", MaxSize: 100, MaxAge: 7}
  e.errFile = zapcore.AddSync(lumberjackLogger)
  ```
- **错误处理**：忽略`OpenFile`的错误返回值，需补充错误检查。

---

### 四、运行效果示例
执行后会产生以下输出：
1. **控制台**：
   ```
   2025-04-01 14:30:00 INFO  main.go:30 this is info
   2025-04-01 14:30:00 WARN  main.go:31 this is warn
   2025-04-01 14:30:00 ERROR main.go:32 this is error1
   2025-04-01 14:30:00 ERROR main.go:33 this is error2
   ```
2. **err.log**：
   ```
   2025-04-01 14:30:00 ERROR main.go:32 this is error1
   2025-04-01 14:30:00 ERROR main.go:33 this is error2
   ```

---

### 五、扩展知识（来自搜索结果）
- **编码器配置**：可通过`EncoderConfig`自定义键名、堆栈跟踪等字段。
- **异步写入**：使用`zapcore.BufferedWriteSyncer`提升高频日志性能（参考网页7）。
- **全局替换**：调用`zap.ReplaceGlobals(logger)`使`zap.L()`全局生效。

---

以下是对 `EncodeEntry` 方法的详细解释，结合 Zap 日志库的设计原理和代码功能：

---

### 一、`EncodeEntry` 的核心作用
`EncodeEntry` 是 Zap 日志库中 `Encoder` 接口的关键方法，负责将日志条目（`zapcore.Entry`）和附加字段（`[]zapcore.Field`）序列化为字节流，并写入缓冲区（`*buffer.Buffer`）。其核心作用包括：
1. **日志格式化**：将日志的元数据（如时间、级别、调用者）和用户自定义字段转换为特定格式（如 JSON、Console 等）。
2. **扩展功能入口**：允许开发者通过重写此方法实现自定义逻辑（如错误日志分离、染色输出等）。

---

### 二、代码功能解析
用户提供的 `EncodeEntry` 方法实现了以下功能：
```go
func (e *levelEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
    // 调用原始编码器生成日志内容
    buf, err := e.Encoder.EncodeEntry(entry, fields)
    if entry.Level == zapcore.ErrorLevel {
        // 初始化错误日志文件（首次触发时创建）
        if e.errFile == nil {
            file, _ := os.OpenFile("err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
            e.errFile = file
        }
        // 将错误日志写入独立文件
        e.errFile.WriteString(buf.String())
    }
    return buf, nil
}
```

#### 1. **错误日志分离**
- **逻辑**：当检测到日志级别为 `ErrorLevel` 时，将日志内容额外写入 `err.log` 文件。
- **实现细节**：
    - 文件按需创建：首次遇到错误日志时通过 `os.OpenFile` 初始化文件句柄。
    - 写入方式：使用 `APPEND` 模式追加内容，避免覆盖历史日志。
- **效果**：普通日志仍输出到控制台，错误日志同时写入独立文件，便于监控和排查。

#### 2. **编码器继承与扩展**
- **继承原始功能**：通过调用 `e.Encoder.EncodeEntry(entry, fields)` 复用父编码器（如 Console 或 JSON 编码器）的默认逻辑。
- **扩展性**：在保留原始日志格式的基础上，仅对错误日志添加额外处理，体现了 Zap 编码器的可组合性。

---

### 三、在 Zap 日志流程中的定位
`EncodeEntry` 是 Zap 日志处理流水线的关键环节：
1. **调用时机**：由 `ioCore.Write` 方法触发，在日志级别检查通过后执行序列化。
2. **性能优化**：Zap 通过手动拼接字符串（而非反射）实现高效序列化，此方法正是性能优化的核心节点。
3. **输出控制**：生成的缓冲区内容最终会通过 `WriteSyncer` 写入目标（如控制台、文件、网络等）。

---

### 四、潜在问题与改进建议
#### 1. **并发安全问题**
- **问题**：多协程并发写入 `errFile` 可能导致数据混乱或竞态条件。
- **改进**：使用互斥锁（`sync.Mutex`）保护文件写入操作，或改用线程安全的 `zapcore.WriteSyncer`。

#### 2. **资源泄漏风险**
- **问题**：未在程序退出时关闭文件句柄，可能导致日志丢失。
- **改进**：在 `levelEncoder` 结构体中添加 `Close()` 方法，并在程序退出时调用 `defer e.errFile.Close()`。

#### 3. **性能优化**
- **问题**：高频错误日志场景下频繁调用 `WriteString` 可能成为瓶颈。
- **改进**：引入缓冲写入（如 `bufio.Writer`）或异步写入机制。

---

### 五、关联设计模式（来自搜索结果）
- **Encoder 链式扩展**：通过内嵌原始编码器并重写 `EncodeEntry`，实现类似装饰器模式的功能增强。
- **延迟初始化（Lazy Initialization）**：文件句柄在首次遇到错误日志时才创建，节省资源。
- **关注点分离**：将格式化和写入逻辑分离，符合单一职责原则。

---

### 附：运行流程示意图
```
Logger.Error() 
  → Core.Check() 
  → ioCore.Write() 
  → levelEncoder.EncodeEntry() 
    → 原始编码器生成日志 
    → 错误日志写入文件 
  → 缓冲区内容返回并输出到控制台
```


# 项目中的日志的使用

把以上功能做一个整合

1. 可设置级别
2. 时间格式化
3. 输出美化
4. 日志前缀
5. 日志双写
6. 日志时间分片、单独把error的日志分出来

```go
package main

import (
  "fmt"
  "go.uber.org/zap"
  "go.uber.org/zap/buffer"
  "go.uber.org/zap/zapcore"
  "os"
  "time"
)

// logEncoder 时间分片和level分片同时做
type logEncoder struct {
  zapcore.Encoder
  errFile     *os.File
  file        *os.File
  currentDate string
}

const (
  BuleColor   = "\033[34m"
  YellowColor = "\033[33m"
  RedColor    = "\033[31m"
  ResetColor  = "\033[0m"
)

func myEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
  switch level {
  case zapcore.InfoLevel:
    enc.AppendString(BuleColor + "INFO" + ResetColor)
  case zapcore.WarnLevel:
    enc.AppendString(YellowColor + "WARN" + ResetColor)
  case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
    enc.AppendString(RedColor + "ERROR" + ResetColor)
  default:
    enc.AppendString(level.String())
  }
}
func (e *logEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
  // 先调用原始的 EncodeEntry 方法生成日志行
  buff, err := e.Encoder.EncodeEntry(entry, fields)
  if err != nil {
    return nil, err
  }
  data := buff.String()
  buff.Reset()
  buff.AppendString("[myApp] " + data)
  data = buff.String()
  // 时间分片
  now := time.Now().Format("2006-01-02")
  if e.currentDate != now {
    os.MkdirAll(fmt.Sprintf("logs/%s", now), 0666)
    // 时间不同，先创建目录
    name := fmt.Sprintf("logs/%s/out.log", now)
    file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    e.file = file
    e.currentDate = now
  }

  switch entry.Level {
  case zapcore.ErrorLevel:
    if e.errFile == nil {
      name := fmt.Sprintf("logs/%s/err.log", now)
      file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
      e.errFile = file
    }
    e.errFile.WriteString(buff.String())
  }

  if e.currentDate == now {
    e.file.WriteString(data)
  }
  return buff, nil
}

func InitLog() *zap.Logger {
  // 使用 zap 的 NewDevelopmentConfig 快速配置
  cfg := zap.NewDevelopmentConfig()
  cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 替换时间格式化方式
  cfg.EncoderConfig.EncodeLevel = myEncodeLevel
  // 创建自定义的 Encoder
  encoder := &logEncoder{
    Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器
  }
  // 创建 Core
  core := zapcore.NewCore(
    encoder,                    // 使用自定义的 Encoder
    zapcore.AddSync(os.Stdout), // 输出到控制台
    zapcore.InfoLevel,          // 设置日志级别
  )
  // 创建 Logger
  logger := zap.New(core, zap.AddCaller())

  zap.ReplaceGlobals(logger)
  return logger
}

func main() {

  logger := InitLog()
  logger.Info("this is info")
  logger.Warn("this is warn")
  logger.Error("this is error1")
  logger.Error("this is error2")
  zap.S().Infof("%s xxx", "fengfeng")
}
```

这段代码实现了一个基于 Uber 的 Zap 日志库的自定义日志系统，主要包含**彩色控制台输出**、**按日期分片存储日志**和**错误日志单独存储**等功能。以下是逐部分解析：

---

### 一、核心结构体 `logEncoder`
```go
type logEncoder struct {
    zapcore.Encoder          // 内嵌 Zap 的编码器接口
    errFile     *os.File     // 存储错误日志的文件句柄
    file        *os.File     // 存储普通日志的文件句柄
    currentDate string       // 记录当前日志日期（用于分片）
}
```
这是自定义编码器的核心结构，通过内嵌 `zapcore.Encoder` 继承 Zap 的编码功能，并扩展了文件分片逻辑。

---

### 二、彩色日志级别输出
```go
func myEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
    switch level {
    case zapcore.InfoLevel:
        enc.AppendString(BuleColor + "INFO" + ResetColor) // 蓝色 INFO
    case zapcore.WarnLevel:
        enc.AppendString(YellowColor + "WARN" + ResetColor) // 黄色 WARN
    case zapcore.ErrorLevel, ...:
        enc.AppendString(RedColor + "ERROR" + ResetColor) // 红色 ERROR
    }
}
```
通过 ANSI 转义码为不同日志级别添加颜色，使控制台输出更易区分。此功能在开发环境中常见（如 `zap.NewDevelopment()` 的默认行为）。

---

### 三、日志分片与文件写入
在 `EncodeEntry` 方法中实现了以下逻辑：
1. **时间分片**：
   ```go
   now := time.Now().Format("2006-01-02")
   if e.currentDate != now {
       os.MkdirAll(fmt.Sprintf("logs/%s", now), 0666) // 按日期创建目录
       name := fmt.Sprintf("logs/%s/out.log", now)   // 普通日志路径
       file, _ := os.OpenFile(name, ...)             // 打开新文件
       e.file = file
       e.currentDate = now
   }
   ```
   每天生成一个目录（如 `logs/2025-04-01/out.log`），实现按日期分片存储。

2. **错误日志分离**：
   ```go
   case zapcore.ErrorLevel:
       if e.errFile == nil {
           name := fmt.Sprintf("logs/%s/err.log", now) // 错误日志路径
           file, _ := os.OpenFile(name, ...)
           e.errFile = file
       }
       e.errFile.WriteString(buff.String()) // 写入错误日志
   ```
   将 `ERROR` 及以上级别的日志单独写入 `err.log`，便于后续监控和分析。

---

### 四、初始化配置 `InitLog`
```go
func InitLog() *zap.Logger {
    cfg := zap.NewDevelopmentConfig()
    cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 自定义时间格式
    cfg.EncoderConfig.EncodeLevel = myEncodeLevel // 绑定彩色编码函数

    encoder := &logEncoder{
        Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 控制台编码器
    }

    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(os.Stdout), // 输出到控制台
        zapcore.InfoLevel,         // 最低日志级别为 INFO
    )

    logger := zap.New(core, zap.AddCaller()) // 添加调用位置信息
    zap.ReplaceGlobals(logger)               // 替换全局 Logger
    return logger
}
```
- 使用 `NewDevelopmentConfig` 快速配置开发环境参数。
- 通过 `zap.AddCaller()` 在日志中记录调用文件名和行号，便于调试。
- `ReplaceGlobals` 允许通过 `zap.L()` 或 `zap.S()` 全局访问 Logger。

---

### 五、潜在问题与改进建议
1. **并发安全**：
    - 当前文件写入未加锁，多协程下可能导致数据混乱。可通过 `sync.Mutex` 或使用 `lumberjack` 库（如网页7示例）解决。

2. **文件关闭**：
    - 未在程序退出时关闭文件句柄，可能导致日志丢失。应在 `defer` 中调用 `file.Close()`。

3. **性能优化**：
    - 高频日志场景下频繁打开文件会影响性能。建议结合异步写入或缓冲机制。

---

### 六、运行示例
执行 `main` 函数会输出：
```log
[myApp] 2025-04-01 14:30:00 INFO  main.go:50 this is info
[myApp] 2025-04-01 14:30:00 WARN  main.go:51 this is warn
[myApp] 2025-04-01 14:30:00 ERROR main.go:52 this is error1
[myApp] 2025-04-01 14:30:00 ERROR main.go:53 this is error2
```
- 控制台显示彩色级别，普通日志写入 `logs/当前日期/out.log`，错误日志额外写入 `err.log`。

如果需要更完整的日志滚动功能（如按大小切割），可参考网页7中 `lumberjack` 的集成方案。

---


以下是对 `InitLog()` 函数的详细技术解析，结合 Zap 日志库的设计原理和代码实现逻辑：

---

### 一、核心功能定位
该函数用于创建**高性能结构化日志系统**，主要实现以下功能：
1. **开发环境优化**：采用易读的控制台输出格式
2. **深度定制**：重定义时间格式与日志级别染色
3. **全局可用**：替换全局 Logger 实现统一调用

---

### 二、逐行代码解析
#### 1. 基础配置初始化
```go
cfg := zap.NewDevelopmentConfig()
```
- **作用**：创建开发环境默认配置（彩色输出、DEBUG级别、调用者信息）
- **特点**：相比生产环境配置（`NewProductionConfig`）：
    - 时间格式改为易读的 `ISO8601` 而非 Unix 时间戳
    - 日志级别显示为小写字符串（如 `info`）而非数字编码
    - 默认输出到 `stderr`（本示例后续覆盖为 `stdout`）

#### 2. 时间格式定制
```go
cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
```
- **实现原理**：覆盖默认的 `ISO8601` 格式，采用精确到秒的时间戳
- **格式符号**：
    - `2006`：年份固定值（Go语言时间格式约定）
    - `01`：两位月份
    - `02`：两位日期
    - `15`：24小时制小时
    - `04`：分钟
    - `05`：秒

#### 3. 日志级别染色
```go
cfg.EncoderConfig.EncodeLevel = myEncodeLevel
```
- **关联函数**：`myEncodeLevel` 通过 ANSI 转义码实现彩色输出：
  ```go
  func myEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
    switch level {
    case zapcore.InfoLevel:
        enc.AppendString("\033[34mINFO\033[0m") // 蓝色
    case zapcore.WarnLevel:
        enc.AppendString("\033[33mWARN\033[0m") // 黄色
    case zapcore.ErrorLevel,...:
        enc.AppendString("\033[31mERROR\033[0m") // 红色
    }
  }
  ```
- **开发环境优势**：通过颜色快速识别关键日志（如错误信息）

#### 4. 核心组件构建
```go
core := zapcore.NewCore(
    &logEncoder{Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig)},
    zapcore.AddSync(os.Stdout),
    zapcore.InfoLevel,
)
```
- **三要素解析**：
    - **编码器**：自定义的 `logEncoder` 包装控制台编码器，实现格式扩展（如添加 `[myApp]` 前缀）
    - **输出目标**：`AddSync(os.Stdout)` 将日志定向到控制台（覆盖默认的 `stderr`）
    - **日志级别**：`InfoLevel` 过滤 DEBUG 及更低级别日志

#### 5. 全局替换
```go
zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))
```
- **链式调用**：
    1. `zap.New()` 创建带调用者信息的 Logger（显示文件名和行号）
    2. `ReplaceGlobals()` 使 `zap.L()` 和 `zap.S()` 全局可用
- **生产建议**：配合 `zap.AddCallerSkip(1)` 调整调用栈深度

---

### 三、性能优化特性
1. **零分配设计**：通过 `buffer` 复用减少内存分配（参考 `logEncoder.EncodeEntry` 实现）
2. **同步写入**：直接使用 `os.Stdout` 而非缓冲 IO，避免日志延迟
3. **结构化存储**：原始日志保持键值对结构，便于后续分析处理

---

### 四、潜在改进方向
1. **并发安全**：文件写入操作需加 `sync.Mutex`（当前示例未涉及文件输出）
2. **错误处理**：补充 `OpenFile` 错误检查，避免静默失败
3. **日志切割**：集成 `lumberjack` 实现自动分片（参考网页4实现）：
   ```go
   writeSyncer := zapcore.AddSync(&lumberjack.Logger{
       Filename: "app.log", 
       MaxSize: 100 // MB
   })
   ```

---

该实现充分体现了 Zap 的高扩展性特点，通过组合编码器、输出目标和过滤级别，构建出符合开发需求的日志系统。生产环境建议增加异步写入和采样策略。

---

这段代码之所以能同时输出到控制台和文件，关键在于**自定义编码器 `logEncoder` 的设计**，其内部通过重写 `EncodeEntry` 方法实现了双重写入机制。以下是具体原因分析：

---

### 一、核心逻辑分离
#### 1. **`zapcore.Core` 的职责**
代码中 `zapcore.NewCore` 的第二个参数 `zapcore.AddSync(os.Stdout)` 仅负责将日志**原始数据流输出到控制台**。它的作用是将经过编码器处理后的日志内容（`buffer.Buffer` 对象）同步到标准输出（`os.Stdout`）。

#### 2. **自定义编码器的扩展**
真正的文件写入逻辑隐藏在自定义编码器 `logEncoder` 的 `EncodeEntry` 方法中：
```go
func (e *logEncoder) EncodeEntry(...) (*buffer.Buffer, error) {
    // 调用父类编码器生成控制台格式日志
    buff, _ := e.Encoder.EncodeEntry(entry, fields)
    
    // 文件写入逻辑
    now := time.Now().Format("2006-01-02")
    if e.currentDate != now {
        // 按日期创建日志目录和文件
        file, _ := os.OpenFile(fmt.Sprintf("logs/%s/out.log", now), ...)
        e.file = file
    }
    e.file.WriteString(data) // 写入普通日志文件
    
    // 错误日志分离
    if entry.Level >= zapcore.ErrorLevel {
        e.errFile.WriteString(buff.String()) // 写入错误日志文件
    }
    
    return buff // 返回的缓冲区内容会通过 Core 的 WriteSyncer 输出到控制台
}
```
**关键点**：
- `Core` 的 `WriteSyncer` 负责控制台输出
- `logEncoder` 在编码过程中**直接操作文件句柄**，实现文件写入

---

### 二、双重输出机制图解
```
日志生成 → 自定义编码器处理 → 返回缓冲区内容
          ↳ 文件写入（直接操作文件句柄）  
          ↳ 控制台输出（通过 Core 的 WriteSyncer）
```

---

### 三、设计模式分析
#### 1. **装饰器模式**
- `logEncoder` 内嵌 `zapcore.Encoder`，通过重写方法扩展功能
- 在生成控制台日志格式的基础上，追加文件写入逻辑

#### 2. **副作用式写入**
- 文件写入操作完全独立于 Zap 的 `WriteSyncer` 体系
- 直接通过 `os.File.WriteString` 实现，绕过了 Zap 的同步机制

---

### 四、潜在问题与改进建议
#### 1. **并发安全问题**
当前代码未对文件写入加锁，多协程环境下可能导致日志混乱。建议：
```go
// 在结构体中添加互斥锁
type logEncoder struct {
    ...
    mu sync.Mutex
}

// 在写入文件时加锁
e.mu.Lock()
defer e.mu.Unlock()
e.file.WriteString(data)
```

#### 2. **性能优化**
高频文件写入可能成为瓶颈，可参考网页7方案集成异步缓冲：
```go
// 使用 BufferedWriteSyncer 包装文件写入
fileSyncer := zapcore.BufferedWriteSyncer{
    WS:            zapcore.AddSync(e.file),
    Size:          4096,  // 缓冲区大小
    FlushInterval: time.Second * 5,
}
fileSyncer.Write(buff.Bytes())
```

#### 3. **标准方案对比**
若需更规范的日志切割，建议改用网页1提到的 `lumberjack` 方案：
```go
// 替换手动文件操作
lumberJackLogger := &lumberjack.Logger{
    Filename: "logs/app.log",
    MaxSize: 100, // MB
    MaxAge: 30,   // 保留天数
}
core := zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), level)
```

---

### 五、总结
该实现通过**混合使用 Zap 标准输出通道和编码器副作用写入**，实现了控制台与文件的双重输出。这种设计虽然灵活，但打破了 Zap 原有的职责分离原则，建议生产环境优先采用通过 `zapcore.NewTee` 创建多核心的标准方案。

---

以上只是zap的基本使用，掌握之后再去探索zap的高级功能，就会更加的得心应手！





