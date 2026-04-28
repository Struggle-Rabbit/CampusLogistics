package building

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type BuildingController struct {
	srv *service.ServiceProvider
}

func NewBuildingController(srv *service.ServiceProvider) *BuildingController {
	return &BuildingController{srv: srv}
}

// Create 创建楼栋
// @Summary 创建楼栋接口
// @Description 创建新的楼栋信息，包含楼栋编号、名称、楼层数、房间数等，编号唯一性自动校验
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param data body dto.BuildingCreateReq true "楼栋创建请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/create [post]
func (s *BuildingController) Create(c *gin.Context) {
	var req dto.BuildingCreateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.BuildingService.Create(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "创建成功")
}

// Update 更新楼栋
// @Summary 更新楼栋接口
// @Description 根据楼栋ID更新楼栋信息，包括楼栋编号、名称、楼层数、房间数等
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param data body dto.BuildingUpdateReq true "楼栋更新请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/update [post]
func (s *BuildingController) Update(c *gin.Context) {
	var req dto.BuildingUpdateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.BuildingService.Update(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "更新成功")
}

// Delete 删除楼栋
// @Summary 删除楼栋接口
// @Description 根据楼栋ID列表删除楼栋信息，删除前会检查是否有关联的宿舍
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "删除参数，ids为楼栋ID数组"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/del [post]
func (s *BuildingController) Delete(c *gin.Context) {
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

	if err := s.srv.BuildingService.Delete(idStrs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// GetListByPage 分页查询楼栋列表
// @Summary 分页查询楼栋列表接口
// @Description 分页查询楼栋信息列表，支持按校区ID、楼栋编号、楼栋名称筛选
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param data body dto.BuildingListPageReq false "查询参数"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.BuildingResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/list [get]
func (s *BuildingController) GetListByPage(c *gin.Context) {
	var req dto.BuildingListPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := s.srv.BuildingService.GetListByPage(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetDetail 获取楼栋详情
// @Summary 获取楼栋详情接口
// @Description 根据楼栋ID获取楼栋的详细信息，包含所属校区名称和已分配房间数
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param id query string true "楼栋ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.BuildingResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/detail [get]
func (s *BuildingController) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "请提供楼栋ID")
		return
	}

	res, err := s.srv.BuildingService.GetDetail(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetBuildingsByCampus 获取校区下楼栋列表
// @Summary 获取校区下楼栋列表接口
// @Description 根据校区ID获取该校区下的所有楼栋信息，用于下拉框选择
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param campus_id query string true "校区ID"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.BuildingResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/byCampus [get]
func (s *BuildingController) GetBuildingsByCampus(c *gin.Context) {
	campusID := c.Query("campus_id")
	if campusID == "" {
		utils.Fail(c, "请提供校区ID")
		return
	}

	res, err := s.srv.BuildingService.GetBuildingsByCampus(campusID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// ImportBuildings 批量导入楼栋
// @Summary 批量导入楼栋接口
// @Description 通过CSV文件批量导入楼栋信息，文件格式：校区名称,楼栋编号,楼栋名称,楼层数,房间数,备注
// @Tags 楼栋管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV文件"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/import [post]
func (s *BuildingController) ImportBuildings(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, "请上传文件")
		return
	}

	count, err := s.srv.BuildingService.ImportBuildings(file)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, map[string]interface{}{"success_count": count}, "导入成功")
}

// ExportBuildings 批量导出楼栋
// @Summary 批量导出楼栋接口
// @Description 导出楼栋信息为CSV文件，支持按校区和楼栋编号筛选
// @Tags 楼栋管理
// @Accept json
// @Produce json
// @Param data query dto.BuildingExportReq false "导出筛选参数"
// @Success 200 {file} file "CSV文件"
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/building/export [get]
func (s *BuildingController) ExportBuildings(c *gin.Context) {
	var req dto.BuildingExportReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	csvData, err := s.srv.BuildingService.ExportBuildings(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=buildings.csv")
	c.String(200, csvData)
}
