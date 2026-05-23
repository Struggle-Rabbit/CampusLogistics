package repair

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/repair"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RepairController interface {
	RepairOrderSubmit(c *gin.Context)
	GetListByPage(c *gin.Context)
	GetDetailById(c *gin.Context)
	DelRepairOrderById(c *gin.Context)
	UpdateRepairOrder(c *gin.Context)
	OrderRecord(c *gin.Context)
}

type RepairControllerProvider struct {
	RepairSvc repair.RepairService
}

// RepairOrderSubmit 提交报修单
// @Summary 提交报修单
// @Description 学生或教职工提交报修申请
// @Tags 报修管理
// @Accept json
// @Produce json
// @Param data body dto.RepairOrderSubmitReq true "报修信息"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/repair/submit [post]
func (s *RepairControllerProvider) RepairOrderSubmit(c *gin.Context) {
	var req dto.RepairOrderSubmitReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	userID, _ := c.Get("userID")
	if err := s.RepairSvc.RepairOrderSubmit(userID.(string), &req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "报修提交成功")
}

// GetListByPage 分页查询报修单
// @Summary 分页查询报修单
// @Description 支持按单号、类型、状态、时间范围等条件筛选
// @Tags 报修管理
// @Accept json
// @Produce json
// @Param data query dto.RepairOrderListByPageReq true "查询条件"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.RepairOrderResult}}
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/repair/list [get]
func (s *RepairControllerProvider) GetListByPage(c *gin.Context) {
	var req dto.RepairOrderListByPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := s.RepairSvc.GetListByPage(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetDetailById 查询报修单详情
// @Summary 查询报修单详情
// @Description 根据ID查询详情，包含流转记录
// @Tags 报修管理
// @Accept json
// @Produce json
// @Param id query string true "报修单ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.RepairOrderResult}
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/repair/detail [get]
func (s *RepairControllerProvider) GetDetailById(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "参数ID不能为空")
		return
	}

	res, err := s.RepairSvc.GetDetailById(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// UpdateRepairOrder 更新报修单信息
// @Summary 更新报修单信息
// @Description 仅限待分配状态编辑
// @Tags 报修管理
// @Accept json
// @Produce json
// @Param data body dto.UpdateRepairOrderSubmitReq true "更新信息"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/repair/update [post]
func (s *RepairControllerProvider) UpdateRepairOrder(c *gin.Context) {
	var req dto.UpdateRepairOrderSubmitReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.RepairSvc.UpdateRepairOrder(req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "更新成功")
}

// OrderRecord 报修单状态流转
// @Summary 报修单状态流转
// @Description 记录处理过程并更新状态
// @Tags 报修管理
// @Accept json
// @Produce json
// @Param data body dto.RecordReq true "流转信息"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/repair/record [post]
func (s *RepairControllerProvider) OrderRecord(c *gin.Context) {
	var req dto.RecordReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	userID, _ := c.Get("userID")
	req.UserID = userID.(string)

	if err := s.RepairSvc.OrderRecord(req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "状态更新成功")
}

// DelRepairOrder 删除报修单
// @Summary 删除报修单
// @Description 软删除报修单及关联记录
// @Tags 报修管理
// @Accept json
// @Produce json
// @Param id query string true "报修单ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/repair/del [post]
func (s *RepairControllerProvider) DelRepairOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "参数ID不能为空")
		return
	}

	if err := s.RepairSvc.DelRepairOrderById(id); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "删除成功")
}
