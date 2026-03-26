package main

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/configs"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/pkg/logger"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化配置
	if err := configs.InitConfig("configs/app.yaml"); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	// 2. 初始化日志
	if err := logger.InitLogger(); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}

	// 3. 初始化数据库
	if err := dao.InitDB(); err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 4. 初始化 Redis（可选，后续迭代添加）
	// if err := cache.InitRedis(); err != nil {
	// 	panic(fmt.Sprintf("初始化Redis失败: %v", err))
	// }

	// 5. 初始化 Gin 引擎
	r := gin.Default()
	r.Use(middleware.ExceptionMiddleware()) // 全局异常处理
	r.Use(middleware.CorsMiddleware())      // 跨域处理

	// 6. 注册路由
	router.InitRouter(r)

	// 7. 启动服务
	addr := fmt.Sprintf(":%d", config.GlobalConfig.App.Port)
	logger.Info("服务启动成功，监听地址: %s", addr)
	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("启动服务失败: %v", err))
	}
}
