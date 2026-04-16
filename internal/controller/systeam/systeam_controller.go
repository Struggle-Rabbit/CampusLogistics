package systeam

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SysteamController struct {
	srv *service.ServiceProvider
}

func NewSysteamController(srv *service.ServiceProvider) *SysteamController {
	return &SysteamController{
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
func (s *SysteamController) GetOperationLogListByPage(c *gin.Context) {
	var optLogReq dto.OperationLogByPageReq
	_ = utils.ShouldBind(c, &optLogReq)

	res, err := s.srv.SysteamService.GetOperationLogListByPage(&optLogReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
