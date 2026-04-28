package unittest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/repair"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRepairController(t *testing.T) {
	_, appInstance := SetupTestDB()
	srvProvider := service.NewServiceProvider(appInstance)
	ctl := repair.NewRepairController(srvProvider)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 模拟中间件设置 userID
	r.Use(func(c *gin.Context) {
		c.Set("userID", "TEST_USER_001")
		c.Next()
	})

	api := r.Group("/api/v1/repair")
	{
		api.POST("/submit", ctl.RepairOrderSubmit)
		api.GET("/list", ctl.GetListByPage)
		api.GET("/detail", ctl.GetDetailById)
		api.POST("/update", ctl.UpdateRepairOrder)
		api.POST("/record", ctl.OrderRecord)
		api.POST("/del", ctl.DelRepairOrder)
	}

	var orderID string

	t.Run("提交报修单接口", func(t *testing.T) {
		req := dto.RepairOrderSubmitReq{
			RepairType:  1,
			Address:     "测试地点",
			Description: "测试描述",
			Images:      []string{"http://image.com/1.jpg"},
			Contact:     "张三",
			Phone:       "13800138000",
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("POST", "/api/v1/repair/submit", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("查询列表接口", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("GET", "/api/v1/repair/list?currentPage=1&pageSize=10", nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		list := data["list"].([]interface{})
		assert.NotEmpty(t, list)
		orderID = list[0].(map[string]interface{})["id"].(string)
	})

	t.Run("查询详情接口", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("GET", "/api/v1/repair/detail?id="+orderID, nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, orderID, data["id"])
	})

	t.Run("更新报修单接口", func(t *testing.T) {
		req := dto.UpdateRepairOrderSubmitReq{
			ID:          orderID,
			Description: "更新后的描述",
			RepairType:  2,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("POST", "/api/v1/repair/update", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("状态流转接口", func(t *testing.T) {
		req := dto.RecordReq{
			ID:     orderID,
			Status: 2,
			Remark: "开始处理",
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("POST", "/api/v1/repair/record", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("查询列表接口-带筛选", func(t *testing.T) {
		w := httptest.NewRecorder()
		// 测试时间范围和状态筛选
		reqHttp, _ := http.NewRequest("GET", "/api/v1/repair/list?currentPage=1&pageSize=10&status=2&start_time=2020-01-01&end_time=2030-01-01", nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		list, ok := data["list"].([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, list)
	})

	t.Run("查询详情接口-不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("GET", "/api/v1/repair/detail?id=999999", nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code) // 业务失败也是 200
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(400), resp["code"]) // 假设 CodeFail 是 400
	})

	t.Run("状态流转接口-非法状态", func(t *testing.T) {
		req := dto.RecordReq{
			ID:     orderID,
			Status: 99, // 非法状态
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("POST", "/api/v1/repair/record", bytes.NewBuffer(body))
		reqHttp.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(400), resp["code"])
	})

	t.Run("删除接口", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqHttp, _ := http.NewRequest("POST", "/api/v1/repair/del?id="+orderID, nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
