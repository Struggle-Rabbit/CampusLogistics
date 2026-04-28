package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/campus"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/building"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadCampusRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	campusCtl := campus.NewCampusController(srv)

	campusGroup := api.Group("/campus")
	{
		campusGroup.POST("/create", middleware.PermissionValidator("campus:create"), campusCtl.Create)
		campusGroup.POST("/update", middleware.PermissionValidator("campus:update"), campusCtl.Update)
		campusGroup.POST("/del", middleware.PermissionValidator("campus:del"), campusCtl.Delete)
		campusGroup.GET("/list", middleware.PermissionValidator("campus:list"), campusCtl.GetListByPage)
		campusGroup.GET("/detail", middleware.PermissionValidator("campus:detail"), campusCtl.GetDetail)
		campusGroup.GET("/all", middleware.PermissionValidator("campus:list"), campusCtl.GetAll)
	}
}

func LoadBuildingRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	buildingCtl := building.NewBuildingController(srv)

	buildingGroup := api.Group("/building")
	{
		buildingGroup.POST("/create", middleware.PermissionValidator("building:create"), buildingCtl.Create)
		buildingGroup.POST("/update", middleware.PermissionValidator("building:update"), buildingCtl.Update)
		buildingGroup.POST("/del", middleware.PermissionValidator("building:del"), buildingCtl.Delete)
		buildingGroup.GET("/list", middleware.PermissionValidator("building:list"), buildingCtl.GetListByPage)
		buildingGroup.GET("/detail", middleware.PermissionValidator("building:detail"), buildingCtl.GetDetail)
		buildingGroup.GET("/byCampus", middleware.PermissionValidator("building:list"), buildingCtl.GetBuildingsByCampus)
		buildingGroup.POST("/import", middleware.PermissionValidator("building:import"), buildingCtl.ImportBuildings)
		buildingGroup.GET("/export", middleware.PermissionValidator("building:export"), buildingCtl.ExportBuildings)
	}
}
