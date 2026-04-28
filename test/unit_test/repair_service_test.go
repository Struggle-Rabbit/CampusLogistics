package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/repair"
	"github.com/stretchr/testify/assert"
)

func TestRepairService(t *testing.T) {
	_, appInstance := SetupTestDB()
	svc := repair.NewRepairService(appInstance)

	var orderID string
	userID := "STUDENT001"

	t.Run("提交报修单", func(t *testing.T) {
		req := &dto.RepairOrderSubmitReq{
			RepairType:  1,
			Address:     "宿舍A栋101",
			Description: "水龙头漏水",
			Images:      []string{"img1.jpg"},
			Contact:     "张三",
			Phone:       "13811112222",
		}
		err := svc.RepairOrderSubmit(userID, req)
		assert.NoError(t, err)

		// 验证列表
		res, err := svc.GetListByPage(&dto.RepairOrderListByPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
		orderID = res.List.([]*dto.RepairOrderResult)[0].ID
	})

	t.Run("更新报修单", func(t *testing.T) {
		req := dto.UpdateRepairOrderSubmitReq{
			ID:          orderID,
			Description: "水龙头爆裂",
			Status:      1,
		}
		err := svc.UpdateRepairOrder(req)
		assert.NoError(t, err)

		detail, err := svc.GetDetailById(orderID)
		assert.NoError(t, err)
		assert.Equal(t, "水龙头爆裂", detail.Description)
	})

	t.Run("订单流转记录", func(t *testing.T) {
		req := dto.RecordReq{
			ID:     orderID,
			Status: 2, // 待处理
			UserID: "STAFF001",
			Remark: "分派给李四处理",
		}
		err := svc.OrderRecord(req)
		assert.NoError(t, err)

		detail, _ := svc.GetDetailById(orderID)
		assert.Equal(t, 2, detail.Status)
		assert.Equal(t, "STAFF001", *detail.HandlerID)
	})

	t.Run("删除订单", func(t *testing.T) {
		err := svc.DelRepairOrderById(orderID)
		assert.NoError(t, err)

		detail, err := svc.GetDetailById(orderID)
		assert.Error(t, err)
		assert.Nil(t, detail)
	})
}
