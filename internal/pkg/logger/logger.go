package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// 全局 Logger 实例
var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

// Config 日志配置结构（支持 YAML 解析）
type Config struct {
	Level           string `yaml:"level"`            // 日志级别: debug/info/warn/error/fatal
	EnableConsole   bool   `yaml:"enable_console"`   // 是否输出到控制台
	Filename        string `yaml:"filename"`         // 日志文件路径（为空则不输出文件）
	MaxSize         int    `yaml:"max_size"`         // 单个文件最大大小(MB)
	MaxBackups      int    `yaml:"max_backups"`      // 保留旧文件最大数量
	MaxAge          int    `yaml:"max_age"`          // 保留旧文件最大天数
	Compress        bool   `yaml:"compress"`         // 是否压缩旧文件
	Encoding        string `yaml:"encoding"`         // 编码格式: console(开发)/json(生产)
	EnableCaller    bool   `yaml:"enable_caller"`    // 是否显示调用者信息
	StacktraceLevel string `yaml:"stacktrace_level"` // 记录堆栈的最低级别
}

// NewDevelopmentConfig 开发环境预设配置
func NewDevelopmentConfig() *Config {
	return &Config{
		Level:           "debug",
		EnableConsole:   true,
		Filename:        "",
		Encoding:        "console",
		EnableCaller:    true,
		StacktraceLevel: "error",
	}
}

// NewProductionConfig 生产环境预设配置
func NewProductionConfig() *Config {
	return &Config{
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

// InitFromFile 从 YAML 文件初始化日志
func InitFromFile(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return Init(&cfg)
}

// Init 从 Config 结构初始化日志
func Init(cfg *Config) error {
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

func buildEncoder(cfg *Config) zapcore.Encoder {
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

func buildWriteSyncer(cfg *Config) (zapcore.WriteSyncer, error) {
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
