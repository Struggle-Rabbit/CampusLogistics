package configs

import (
	"fmt"

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
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type LogConfig struct {
	Level         string `mapstructure:"level"`
	Encoding      string `mapstructure:"encoding"`
	EnableConsole bool   `mapstructure:"enable_console"`
	Filename      string `mapstructure:"filename"`
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
func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
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
