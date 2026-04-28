package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/campus"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
	"github.com/stretchr/testify/assert"
)

func TestCampusAndBuildingService(t *testing.T) {
	_, appInstance := SetupTestDB()
	campusSvc := campus.NewCampusService(appInstance)
	buildingSvc := building.NewBuildingService(appInstance)

	t.Run("创建校区", func(t *testing.T) {
		req := &dto.CampusCreateReq{
			CampusName: "济南校区",
			Address:    "济南市历下区",
			Contact:    "张老师",
			Phone:      "0531-12345678",
		}
		err := campusSvc.Create(req)
		assert.NoError(t, err)

		res, err := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
	})

	t.Run("更新校区", func(t *testing.T) {
		listRes, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		campusList := listRes.List.([]*dto.CampusResult)
		c := campusList[0]

		updateReq := &dto.CampusUpdateReq{
			ID:         c.ID,
			CampusName: "济南校区（更新）",
			Address:    "济南市高新区",
			Contact:    "李老师",
			Phone:      "0531-87654321",
		}
		err := campusSvc.Update(updateReq)
		assert.NoError(t, err)

		detail, _ := campusSvc.GetDetail(c.ID)
		assert.Equal(t, "济南校区（更新）", detail.CampusName)
	})

	t.Run("创建楼栋", func(t *testing.T) {
		listRes, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		campusList := listRes.List.([]*dto.CampusResult)
		c := campusList[0]

		req := &dto.BuildingCreateReq{
			CampusID:     c.ID,
			BuildingNo:   "A01",
			BuildingName: "1号教学楼",
			FloorCount:   6,
			RoomCount:    100,
		}
		err := buildingSvc.Create(req)
		assert.NoError(t, err)
	})

	t.Run("楼栋编号唯一性校验", func(t *testing.T) {
		listRes, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		campusList := listRes.List.([]*dto.CampusResult)
		c := campusList[0]

		req := &dto.BuildingCreateReq{
			CampusID:     c.ID,
			BuildingNo:   "A01",
			BuildingName: "重复楼栋",
			FloorCount:   3,
			RoomCount:    50,
		}
		err := buildingSvc.Create(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "楼栋编号已存在")
	})

	t.Run("删除校区失败-有关联楼栋", func(t *testing.T) {
		listRes, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		campusList := listRes.List.([]*dto.CampusResult)
		c := campusList[0]

		err := campusSvc.Delete([]string{c.ID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "该校区下存在楼栋信息")
	})

	t.Run("获取校区下楼栋列表", func(t *testing.T) {
		listRes, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		campusList := listRes.List.([]*dto.CampusResult)
		c := campusList[0]

		buildings, err := buildingSvc.GetBuildingsByCampus(c.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(buildings), 1)
	})

	t.Run("获取楼栋详情", func(t *testing.T) {
		listRes, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		campusList := listRes.List.([]*dto.CampusResult)
		c := campusList[0]

		buildings, _ := buildingSvc.GetBuildingsByCampus(c.ID)
		detail, err := buildingSvc.GetDetail(buildings[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, "A01", detail.BuildingNo)
	})
}
