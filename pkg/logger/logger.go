package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局 Logger 实例
var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

// NewDevelopmentConfig 开发环境预设配置
func NewDevelopmentConfig() *config.LogConfig {
	return &config.LogConfig{
		Level:           "debug",
		EnableConsole:   true,
		Filename:        "",
		Encoding:        "console",
		EnableCaller:    true,
		StacktraceLevel: "error",
	}
}

// NewProductionConfig 生产环境预设配置
func NewProductionConfig() *config.LogConfig {
	return &config.LogConfig{
		Level:           "info",
		EnableConsole:   false,
		Filename:        "./logs/app.log",
		MaxSize:         100,
		MaxBackups:      10,
		MaxAge:          30,
		Compress:        true,
		Encoding:        "json",
		EnableCaller:    true,
		StacktraceLevel: "error",
	}
}

// InitLogger 从 Config 结构初始化日志
func InitLogger() error {

	// 日志初始化
	var LogConfig = config.GlobalConfig.Log
	fmt.Println("日志初始化中....")
	var cfg *config.LogConfig
	if config.IsDev() {
		cfg = NewDevelopmentConfig()
	} else {
		cfg = &config.LogConfig{
			Level:         LogConfig.Level,
			EnableConsole: LogConfig.EnableConsole,
			Filename:      LogConfig.Filename,
			Encoding:      LogConfig.Encoding,
		}
	}

	// 1. 解析日志级别
	logLevel := parseLevel(cfg.Level, zapcore.InfoLevel)
	stackLevel := parseLevel(cfg.StacktraceLevel, zapcore.ErrorLevel)

	// 2. 构建编码器
	encoder := buildEncoder(cfg)

	// 3. 构建输出目标
	writeSyncer, err := buildWriteSyncer(cfg)
	if err != nil {
		return err
	}

	// 4. 创建 Core
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	// 5. 构建 Logger 选项
	opts := []zap.Option{}
	if cfg.EnableCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1)) // 跳过封装层
	}
	opts = append(opts, zap.AddStacktrace(stackLevel))

	// 6. 初始化全局 Logger
	Log = zap.New(core, opts...)
	Sugar = Log.Sugar()

	return nil
}

// Sync 刷新日志缓冲区（程序退出前调用）
func Sync() {
	_ = Log.Sync()
	_ = Sugar.Sync()
}

// ------------------------------ 快捷方法 ------------------------------
// 结构化日志（推荐性能敏感场景使用）
func Debug(msg string, fields ...zap.Field) { Log.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { Log.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Log.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Log.Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Log.Fatal(msg, fields...) }

// 格式化日志（推荐开发调试使用）
func Debugf(template string, args ...interface{}) { Sugar.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { Sugar.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { Sugar.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { Sugar.Errorf(template, args...) }
func Fatalf(template string, args ...interface{}) { Sugar.Fatalf(template, args...) }

// ------------------------------ 内部辅助函数 ------------------------------
func parseLevel(levelStr string, defaultLevel zapcore.Level) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return defaultLevel
	}
	return level
}

func buildEncoder(cfg *config.LogConfig) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if cfg.Encoding == "console" {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder // 彩色级别
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func buildWriteSyncer(cfg *config.LogConfig) (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer

	// 控制台输出
	if cfg.EnableConsole {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// 文件输出
	if cfg.Filename != "" {
		// 自动创建日志目录
		if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}

		// 配置 lumberjack 轮转
		hook := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writers = append(writers, zapcore.AddSync(hook))
	}

	// 默认输出到控制台
	if len(writers) == 0 {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	return zapcore.NewMultiWriteSyncer(writers...), nil
}
