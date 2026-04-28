package dorm

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DormController struct {
	srv *service.ServiceProvider
}

func NewDormController(srv *service.ServiceProvider) *DormController {
	return &DormController{srv: srv}
}

// Create 创建宿舍
// @Summary 创建宿舍接口
// @Description 创建新的宿舍信息，包含宿舍编号、所属楼栋、楼层、房间类型、床位数等
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormCreateReq true "宿舍创建请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/create [post]
func (s *DormController) Create(c *gin.Context) {
	var req dto.DormCreateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.DormService.Create(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "创建成功")
}

// Update 更新宿舍
// @Summary 更新宿舍接口
// @Description 根据宿舍ID更新宿舍信息，包括宿舍编号、楼层、房间类型、床位数等
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormUpdateReq true "宿舍更新请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/update [post]
func (s *DormController) Update(c *gin.Context) {
	var req dto.DormUpdateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.DormService.Update(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "更新成功")
}

// Delete 删除宿舍
// @Summary 删除宿舍接口
// @Description 根据宿舍ID列表删除宿舍信息，删除前会检查是否有学生入住
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "删除参数，ids为宿舍ID数组"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/del [post]
func (s *DormController) Delete(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	ids, ok := data["ids"].([]interface{})
	if !ok || len(ids) == 0 {
		utils.Fail(c, "请选择要删除的数据")
		return
	}

	idStrs := make([]string, 0, len(ids))
	for _, id := range ids {
		if idStr, ok := id.(string); ok {
			idStrs = append(idStrs, idStr)
		}
	}

	if err := s.srv.DormService.Delete(idStrs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// GetListByPage 分页查询宿舍列表
// @Summary 分页查询宿舍列表接口
// @Description 分页查询宿舍信息列表，支持按楼栋、校区、楼层、宿舍类型、状态筛选
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormListPageReq false "查询参数"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.DormResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/list [get]
func (s *DormController) GetListByPage(c *gin.Context) {
	var req dto.DormListPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := s.srv.DormService.GetListByPage(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetDetail 获取宿舍详情
// @Summary 获取宿舍详情接口
// @Description 根据宿舍ID获取宿舍的详细信息，包含楼栋名称、校区名称、床位使用情况
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param id query string true "宿舍ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.DormResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/detail [get]
func (s *DormController) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "请提供宿舍ID")
		return
	}

	res, err := s.srv.DormService.GetDetail(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// AssignDorm 分配宿舍
// @Summary 分配宿舍接口
// @Description 为学生分配宿舍，会检查床位是否已满、用户是否已在其他宿舍入住
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormAssignReq true "分配参数，room_id为宿舍ID，user_id为用户ID"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/assign [post]
func (s *DormController) AssignDorm(c *gin.Context) {
	var req dto.DormAssignReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.DormService.AssignDorm(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "分配成功")
}

// TransferDorm 调宿
// @Summary 调宿接口
// @Description 将学生从当前宿舍调整到目标宿舍，会检查目标宿舍是否有空床位
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormTransferReq true "调宿参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/transfer [post]
func (s *DormController) TransferDorm(c *gin.Context) {
	var req dto.DormTransferReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.DormService.TransferDorm(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "调宿成功")
}

// CheckOut 退宿
// @Summary 退宿接口
// @Description 办理学生退宿手续，自动释放床位
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormCheckOutReq true "退宿参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/checkout [post]
func (s *DormController) CheckOut(c *gin.Context) {
	var req dto.DormCheckOutReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.DormService.CheckOut(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "退宿成功")
}

// GetDormUsers 获取宿舍人员列表
// @Summary 获取宿舍人员列表接口
// @Description 获取指定宿舍的所有在住/历史入住人员列表
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Param data body dto.DormUserListReq false "查询参数"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.DormUserResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/users [get]
func (s *DormController) GetDormUsers(c *gin.Context) {
	var req dto.DormUserListReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := s.srv.DormService.GetDormUsers(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetCapacityWarning 获取容量预警
// @Summary 获取宿舍容量预警接口
// @Description 获取入住率达到80%及以上的宿舍列表，用于预警提醒
// @Tags 宿舍管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=[]dto.DormResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/dorm/warning [get]
func (s *DormController) GetCapacityWarning(c *gin.Context) {
	res, err := s.srv.DormService.GetCapacityWarning()
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
