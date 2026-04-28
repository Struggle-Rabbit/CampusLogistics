package campus

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CampusController struct {
	srv *service.ServiceProvider
}

func NewCampusController(srv *service.ServiceProvider) *CampusController {
	return &CampusController{srv: srv}
}

// Create 创建校区
// @Summary 创建校区接口
// @Description 创建新的校区信息，包含校区名称、地址、联系方式等基础信息
// @Tags 校区管理
// @Accept json
// @Produce json
// @Param data body dto.CampusCreateReq true "校区创建请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/campus/create [post]
func (s *CampusController) Create(c *gin.Context) {
	var req dto.CampusCreateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.CampusService.Create(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "创建成功")
}

// Update 更新校区
// @Summary 更新校区接口
// @Description 根据校区ID更新校区信息，包括校区名称、地址、联系方式等
// @Tags 校区管理
// @Accept json
// @Produce json
// @Param data body dto.CampusUpdateReq true "校区更新请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/campus/update [post]
func (s *CampusController) Update(c *gin.Context) {
	var req dto.CampusUpdateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.CampusService.Update(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "更新成功")
}

// Delete 删除校区
// @Summary 删除校区接口
// @Description 根据校区ID列表删除校区信息，删除前会检查是否有关联的楼栋
// @Tags 校区管理
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "删除参数，ids为校区ID数组"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/campus/del [post]
func (s *CampusController) Delete(c *gin.Context) {
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

	if err := s.srv.CampusService.Delete(idStrs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// GetListByPage 分页查询校区列表
// @Summary 分页查询校区列表接口
// @Description 分页查询校区信息列表，支持按校区名称模糊搜索
// @Tags 校区管理
// @Accept json
// @Produce json
// @Param data body dto.CampusListPageReq false "查询参数，page和page_size必填"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.CampusResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/campus/list [get]
func (s *CampusController) GetListByPage(c *gin.Context) {
	var req dto.CampusListPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := s.srv.CampusService.GetListByPage(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetDetail 获取校区详情
// @Summary 获取校区详情接口
// @Description 根据校区ID获取校区的详细信息，包含楼栋数量统计
// @Tags 校区管理
// @Accept json
// @Produce json
// @Param id query string true "校区ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.CampusResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/campus/detail [get]
func (s *CampusController) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "请提供校区ID")
		return
	}

	res, err := s.srv.CampusService.GetDetail(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetAll 获取所有校区
// @Summary 获取所有校区接口
// @Description 获取所有校区的简要信息列表，用于下拉框选择等场景
// @Tags 校区管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=[]dto.CampusResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/campus/all [get]
func (s *CampusController) GetAll(c *gin.Context) {
	res, err := s.srv.CampusService.GetAll()
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
