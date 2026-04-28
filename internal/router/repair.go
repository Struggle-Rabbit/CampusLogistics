package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/repair"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadRepairRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	repairCtl := repair.NewRepairController(srv)

	repairGroup := api.Group("/repair")
	{
		// 学生/教职工提交报修
		repairGroup.POST("/submit", middleware.PermissionValidator("repair:submit"), repairCtl.RepairOrderSubmit)
		// 分页查询报修单
		repairGroup.GET("/list", middleware.PermissionValidator("repair:list"), repairCtl.GetListByPage)
		// 查询详情
		repairGroup.GET("/detail", middleware.PermissionValidator("repair:detail"), repairCtl.GetDetailById)
		// 更新报修信息 (待分配状态)
		repairGroup.POST("/update", middleware.PermissionValidator("repair:update"), repairCtl.UpdateRepairOrder)
		// 状态流转处理
		repairGroup.POST("/record", middleware.PermissionValidator("repair:record"), repairCtl.OrderRecord)
		// 删除报修单
		repairGroup.POST("/del", middleware.PermissionValidator("repair:del"), repairCtl.DelRepairOrder)
	}
}
