package utility

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UtilityController struct {
	srv *service.ServiceProvider
}

func NewUtilityController(srv *service.ServiceProvider) *UtilityController {
	return &UtilityController{srv: srv}
}

// Create 创建水电费记录
// @Summary 创建水电费记录接口
// @Description 录入宿舍水电费用量数据，自动计算费用，同一宿舍同月份记录不能重复
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body dto.UtilityCreateReq true "水电费记录创建请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/create [post]
func (s *UtilityController) Create(c *gin.Context) {
	var req dto.UtilityCreateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.UtilityService.Create(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "创建成功")
}

// Update 更新水电费记录
// @Summary 更新水电费记录接口
// @Description 更新宿舍水电费用量数据，已缴费记录不允许修改
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body dto.UtilityUpdateReq true "水电费记录更新请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/update [post]
func (s *UtilityController) Update(c *gin.Context) {
	var req dto.UtilityUpdateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.UtilityService.Update(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "更新成功")
}

// Delete 删除水电费记录
// @Summary 删除水电费记录接口
// @Description 删除水电费记录，已缴费记录不允许删除
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "删除参数，ids为记录ID数组"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/del [post]
func (s *UtilityController) Delete(c *gin.Context) {
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

	if err := s.srv.UtilityService.Delete(idStrs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// GetListByPage 分页查询水电费记录
// @Summary 分页查询水电费记录接口
// @Description 分页查询水电费记录列表，支持按宿舍、校区、楼栋、年月、缴费状态筛选
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body dto.UtilityListPageReq false "查询参数"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.UtilityResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/list [get]
func (s *UtilityController) GetListByPage(c *gin.Context) {
	var req dto.UtilityListPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := s.srv.UtilityService.GetListByPage(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetDetail 获取水电费详情
// @Summary 获取水电费详情接口
// @Description 根据记录ID获取水电费的详细信息，包含用量、单价、费用明细
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param id query string true "记录ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.UtilityResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/detail [get]
func (s *UtilityController) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "请提供记录ID")
		return
	}

	res, err := s.srv.UtilityService.GetDetail(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// Pay 缴费
// @Summary 缴费接口
// @Description 为指定水电费记录进行缴费操作
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body dto.UtilityPayReq true "缴费参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/pay [post]
func (s *UtilityController) Pay(c *gin.Context) {
	var req dto.UtilityPayReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.UtilityService.Pay(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "缴费成功")
}

// BatchPay 批量缴费
// @Summary 批量缴费接口
// @Description 批量为多个水电费记录进行缴费操作
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body dto.UtilityBatchPayReq true "批量缴费参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/batchPay [post]
func (s *UtilityController) BatchPay(c *gin.Context) {
	var req dto.UtilityBatchPayReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.UtilityService.BatchPay(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "批量缴费成功")
}

// UpdatePrice 更新单价配置
// @Summary 更新水电单价接口
// @Description 更新水费和电费的单价配置
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param data body dto.UtilityPriceReq true "单价配置参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/price [post]
func (s *UtilityController) UpdatePrice(c *gin.Context) {
	var req dto.UtilityPriceReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := s.srv.UtilityService.UpdatePrice(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "价格更新成功")
}

// GetPrice 获取单价配置
// @Summary 获取水电单价接口
// @Description 获取当前的水费和电费单价配置
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=dto.UtilityPriceResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/price [get]
func (s *UtilityController) GetPrice(c *gin.Context) {
	res, err := s.srv.UtilityService.GetPrice()
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetStatistics 获取统计信息
// @Summary 获取水电费统计接口
// @Description 按校区和月份统计水电费使用量和费用情况
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param campus_id query string false "校区ID"
// @Param year query int false "年份"
// @Param month query int false "月份"
// @Success 200 {object} utils.SuccessResponse{data=dto.UtilityStatResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/statistics [get]
func (s *UtilityController) GetStatistics(c *gin.Context) {
	campusID := c.Query("campus_id")
	year := 0
	month := 0

	if y := c.Query("year"); y != "" {
		year = utils.StrToInt(y)
	}
	if m := c.Query("month"); m != "" {
		month = utils.StrToInt(m)
	}

	res, err := s.srv.UtilityService.GetStatistics(campusID, year, month)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetUnpaidWarning 获取欠费预警
// @Summary 获取欠费预警接口
// @Description 获取当前欠费的水电费记录列表，用于催缴提醒
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=[]dto.UtilityResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/warning [get]
func (s *UtilityController) GetUnpaidWarning(c *gin.Context) {
	res, err := s.srv.UtilityService.GetUnpaidWarning()
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetUserDormUtility 获取用户水电费
// @Summary 获取用户水电费接口
// @Description 获取当前登录用户所在宿舍的水电费记录（需登录认证）
// @Tags 水电费管理
// @Accept json
// @Produce json
// @Param year query int false "年份"
// @Param month query int false "月份"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.UtilityResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/utility/myUtility [get]
func (s *UtilityController) GetUserDormUtility(c *gin.Context) {
	userID, _ := c.Get("user_id")
	year := 0
	month := 0

	if y := c.Query("year"); y != "" {
		year = utils.StrToInt(y)
	}
	if m := c.Query("month"); m != "" {
		month = utils.StrToInt(m)
	}

	res, err := s.srv.UtilityService.GetUserDormUtility(userID.(string), year, month)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
