package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	JWT   JWTConfig   `mapstructure:"jwt"`
	Log   LogConfig   `mapstructure:"log"`
}

type AppConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type LogConfig struct {
	Level           string `mapstructure:"level"`            // 日志级别: debug/info/warn/error/fatal
	Encoding        string `mapstructure:"encoding"`         // 编码格式: console(开发)/json(生产)
	EnableConsole   bool   `mapstructure:"enable_console"`   // 是否输出到控制台
	Filename        string `mapstructure:"filename"`         // 日志文件路径（为空则不输出文件）
	MaxSize         int    `mapstructure:"max_size"`         // 单个文件最大大小(MB)
	MaxBackups      int    `mapstructure:"max_backups"`      // 保留旧文件最大数量
	MaxAge          int    `mapstructure:"max_age"`          // 保留旧文件最大天数
	Compress        bool   `mapstructure:"compress"`         // 是否压缩旧文件
	EnableCaller    bool   `mapstructure:"enable_caller"`    // 是否显示调用者信息
	StacktraceLevel string `mapstructure:"stacktrace_level"` // 记录堆栈的最低级别
}

type MySQLConfig struct {
	DSN             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	AccessExpire  int64  `mapstructure:"access_expire"`
	RefreshExpire int64  `mapstructure:"refresh_expire"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig() error {
	// 获取命令行中ENV的值
	env := os.Getenv("ENV")
	// 若值等于空则默认是开发环境dev
	if env == "" {
		env = "dev"
	}

	v := viper.New()                                   // 待学习
	v.SetConfigName(fmt.Sprintf("config-%s", env))     // 设置配置文件是哪一个，根据获取道德ENV的值来确定
	v.SetConfigType("yaml")                            // 设置配置文件类型
	v.AddConfigPath("configs")                         // 设置需要读取的配置文件的路径，在哪个文件夹下
	v.AutomaticEnv()                                   // 待学习
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 覆盖配置中下划线替换为点，类似于DB_Name替换为DB.Name  待学习

	if err := v.ReadInConfig(); err != nil { // 读取文件，读取配置
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	GlobalConfig = &Config{}
	// 解析文件配置，映射到结构体中
	if err := v.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	// bs, _ := json.MarshalIndent(GlobalConfig, "", "  ")
	// fmt.Println("全局配置：\n", string(bs))
	return nil
}

// IsDev 是否开发环境
func IsDev() bool {
	return GlobalConfig.App.Env == "development" || GlobalConfig.App.Env == "dev"
}

// IsProd 是否生产环境
func IsProd() bool {
	return GlobalConfig.App.Env == "production" || GlobalConfig.App.Env == "prod"
}
