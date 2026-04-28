package system

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SystemController struct {
	srv *service.ServiceProvider
}

func NewSystemController(srv *service.ServiceProvider) *SystemController {
	return &SystemController{
		srv: srv,
	}
}

// GetOperationLogListByPage 操作日志查询
// @Summary 操作日志分页查询接口
// @Description 操作日志分页查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data query dto.OperationLogByPageReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/OperationLogList [get]
func (s *SystemController) GetOperationLogListByPage(c *gin.Context) {
	var optLogReq dto.OperationLogByPageReq
	if !utils.ShouldBind(c, &optLogReq) {
		return
	}

	res, err := s.srv.SystemService.GetOperationLogListByPage(&optLogReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// RefreshToken token刷新
// @Summary token刷新
// @Description 前端使用refresh_token请求
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param refresh_token query string true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/RefreshToken [get]
func (s *SystemController) RefreshToken(c *gin.Context) {
	refresh_token, isExistence := c.GetQuery("refresh_token")
	if !isExistence {
		utils.Fail(c, "token参数为必填")
		return
	}

	res, err := s.srv.SystemService.RefreshToken(refresh_token)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
