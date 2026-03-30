package main

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/configs"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/pkg/logger"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/router"
	"go.uber.org/zap"
)

func main() {
	// 配置初始化
	fmt.Println("配置初始化中....")
	if err := configs.InitConfig(); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	// 日志初始化
	fmt.Println("日志初始化中....")
	var logErr error = nil
	if configs.IsDev() {
		logErr = logger.Init(logger.NewDevelopmentConfig())
	}
	if configs.IsProd() {
		logcfg := &logger.Config{
			Level:         configs.GlobalConfig.Log.Level,
			EnableConsole: configs.GlobalConfig.Log.EnableConsole,
			Filename:      configs.GlobalConfig.Log.Filename,
			Encoding:      configs.GlobalConfig.Log.Encoding,
		}
		logErr = logger.Init(logcfg)
	}
	if logErr != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", logErr))
	}
	defer logger.Sync()

	// 数据库初始化
	fmt.Printf("数据库初始化中....  当前环境: %s\n", configs.GlobalConfig.App.Env)
	if err := dao.InitDB(); err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 初始化 Redis（可选，后续迭代添加）
	// fmt.Println("初始化Redis....")
	// if err := cache.InitRedis(); err != nil {
	// 	panic(fmt.Sprintf("初始化Redis失败: %v", err))
	// }

	// 注册路由
	fmt.Println("注册路由....")
	r := router.InitRouter()
	logger.Info("服务启动",
		zap.String("env", configs.GlobalConfig.App.Env),
		zap.Int("port", configs.GlobalConfig.App.Port),
	)

	fmt.Println("服务启动中....")
	if err := r.Run(fmt.Sprintf(":%d", configs.GlobalConfig.App.Port)); err != nil {
		panic(fmt.Sprintf("启动服务失败: %v", err))
	}

}
