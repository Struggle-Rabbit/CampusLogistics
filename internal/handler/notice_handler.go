package handler

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type noticeHandler struct {
	svc *notice.Service
}

func NewNoticeHandler(svc *notice.Service) *noticeHandler {
	return &noticeHandler{svc: svc}
}

// Create 创建公告
// @Summary 创建公告接口
// @Description 创建新的公告信息，支持定时发布，置顶公告最多3条
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param data body dto.NoticeCreateReq true "公告创建请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/create [post]
func (h *noticeHandler) Create(c *gin.Context) {
	var req dto.NoticeCreateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	userID, _ := c.Get("userID")
	if err := h.svc.Create(userID.(string), &req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "创建成功")
}

// Update 更新公告
// @Summary 更新公告接口
// @Description 根据公告ID更新公告信息，支持修改发布时间和置顶状态
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param data body dto.NoticeUpdateReq true "公告更新请求参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/update [post]
func (h *noticeHandler) Update(c *gin.Context) {
	var req dto.NoticeUpdateReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := h.svc.Update(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "更新成功")
}

// Delete 删除公告
// @Summary 删除公告接口
// @Description 根据公告ID列表删除公告信息
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "删除参数，ids为公告ID数组"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/del [post]
func (h *noticeHandler) Delete(c *gin.Context) {
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

	if err := h.svc.Delete(idStrs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "删除成功")
}

// GetListByPage 分页查询公告列表
// @Summary 分页查询公告列表接口
// @Description 分页查询公告信息列表，支持按标题、类型、置顶状态、时间范围筛选（管理员接口）
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param data body dto.NoticeListPageReq false "查询参数"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.NoticeResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/list [get]
func (h *noticeHandler) GetListByPage(c *gin.Context) {
	var req dto.NoticeListPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := h.svc.GetListByPage(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetPublicList 获取公开公告列表
// @Summary 获取公开公告列表接口
// @Description 获取已发布的公告列表（公开接口，无需登录），仅返回发布时间早于当前时间的公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param data body dto.NoticeListPageReq false "查询参数"
// @Success 200 {object} utils.SuccessResponse{data=dto.PageResult{list=[]dto.NoticePublicResult}}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/public [get]
func (h *noticeHandler) GetPublicList(c *gin.Context) {
	var req dto.NoticeListPageReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	res, err := h.svc.GetPublicList(&req)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetDetail 获取公告详情
// @Summary 获取公告详情接口
// @Description 根据公告ID获取公告的详细信息，自动增加浏览量
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param id query string true "公告ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.NoticeResult}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/detail [get]
func (h *noticeHandler) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "请提供公告ID")
		return
	}

	res, err := h.svc.GetDetail(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// SetTop 设置置顶
// @Summary 设置公告置顶接口
// @Description 设置或取消公告的置顶状态，置顶公告最多3条
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param data body dto.NoticeTopReq true "置顶参数"
// @Success 200 {object} utils.SuccessResponse{data=nil}
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/notice/top [post]
func (h *noticeHandler) SetTop(c *gin.Context) {
	var req dto.NoticeTopReq
	if !utils.ShouldBind(c, &req) {
		return
	}

	if err := h.svc.SetTop(&req); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, nil, "置顶设置成功")
}
