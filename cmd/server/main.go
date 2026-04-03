package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().Format("20060102")
	userCode := fmt.Sprintf("%s00%d", now, 1)
	fmt.Println(userCode)
	// // 配置初始化
	// fmt.Println("配置初始化中....")
	// if err := config.InitConfig(); err != nil {
	// 	panic(fmt.Sprintf("初始化配置失败: %v", err))
	// }

	// if logErr := logger.InitLogger(); logErr != nil {
	// 	panic(fmt.Sprintf("初始化日志失败: %v", logErr))
	// }
	// defer logger.Sync()

	// // 数据库初始化
	// fmt.Printf("数据库初始化中....  当前环境: %s\n", config.GlobalConfig.App.Env)
	// if err := dao.InitDB(); err != nil {
	// 	panic(fmt.Sprintf("初始化数据库失败: %v", err))
	// }

	// //初始化 Redis（可选，后续迭代添加）
	// // fmt.Println("初始化Redis....")
	// // if err := cache.InitRedis(); err != nil {
	// // 	panic(fmt.Sprintf("初始化Redis失败: %v", err))
	// // }

	// if err := router.Run(); err != nil {
	// 	panic(fmt.Sprintf("服务启动失败: %v", err))
	// }
}
