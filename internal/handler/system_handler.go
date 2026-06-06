package handler

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type systemHandler struct {
	svc *system.Service
}

func NewSystemHandler(svc *system.Service) *systemHandler {
	return &systemHandler{svc: svc}
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
func (h *systemHandler) GetOperationLogListByPage(c *gin.Context) {
	var optLogReq dto.OperationLogByPageReq
	if !utils.ShouldBind(c, &optLogReq) {
		return
	}

	res, err := h.svc.GetOperationLogListByPage(&optLogReq)
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
func (h *systemHandler) RefreshToken(c *gin.Context) {
	refreshToken, isExistence := c.GetQuery("refresh_token")
	if !isExistence {
		utils.Fail(c, "token参数为必填")
		return
	}

	res, err := h.svc.RefreshToken(refreshToken)
	if err != nil {
		utils.Unauth(c, "token过期或无效")
		return
	}
	utils.Success(c, res, "获取成功")
}
