package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/campus"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/dorm"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/utility"
	"github.com/stretchr/testify/assert"
)

func TestUtilityService(t *testing.T) {
	_, appInstance := SetupTestDB()
	campusSvc := campus.NewCampusService(appInstance)
	buildingSvc := building.NewBuildingService(appInstance)
	dormSvc := dorm.NewDormService(appInstance)
	utilitySvc := utility.NewUtilityService(appInstance)

	campusSvc.Create(&dto.CampusCreateReq{CampusName: "测试校区"})
	campusList, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10}})
	testCampus := campusList.List.([]*dto.CampusResult)[0]

	buildingSvc.Create(&dto.BuildingCreateReq{
		CampusID: testCampus.ID, BuildingNo: "T01", BuildingName: "测试楼", FloorCount: 3, RoomCount: 30,
	})
	buildingList, _ := buildingSvc.GetBuildingsByCampus(testCampus.ID)
	testBuilding := buildingList[0]

	dormSvc.Create(&dto.DormCreateReq{
		BuildingID: testBuilding.ID, RoomNo: "T101", Floor: 1, RoomType: dto.DormRoomType4, MaxCount: 4,
	})
	dormList, _ := dormSvc.GetListByPage(&dto.DormListPageReq{PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10}})
	testDorm := dormList.List.([]*dto.DormResult)[0]

	t.Run("创建水电费记录", func(t *testing.T) {
		req := &dto.UtilityCreateReq{
			RoomID:        testDorm.ID,
			Year:          2024,
			Month:         3,
			WaterUsage:    10.5,
			ElectricUsage: 150.0,
		}
		err := utilitySvc.Create(req)
		assert.NoError(t, err)

		res, err := utilitySvc.GetListByPage(&dto.UtilityListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
	})

	t.Run("水电费记录唯一性校验", func(t *testing.T) {
		req := &dto.UtilityCreateReq{
			RoomID: testDorm.ID, Year: 2024, Month: 3, WaterUsage: 5, ElectricUsage: 100,
		}
		err := utilitySvc.Create(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "水电记录已存在")
	})

	t.Run("获取单价配置", func(t *testing.T) {
		price, err := utilitySvc.GetPrice()
		assert.NoError(t, err)
		assert.Greater(t, price.WaterPrice, 0.0)
		assert.Greater(t, price.ElectricPrice, 0.0)
	})

	t.Run("更新单价配置", func(t *testing.T) {
		req := &dto.UtilityPriceReq{
			WaterPrice:    4.0,
			ElectricPrice: 0.65,
		}
		err := utilitySvc.UpdatePrice(req)
		assert.NoError(t, err)

		price, _ := utilitySvc.GetPrice()
		assert.Equal(t, 4.0, price.WaterPrice)
		assert.Equal(t, 0.65, price.ElectricPrice)
	})

	t.Run("更新水电费记录", func(t *testing.T) {
		listRes, _ := utilitySvc.GetListByPage(&dto.UtilityListPageReq{PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10}})
		utilList := listRes.List.([]*dto.UtilityResult)
		util := utilList[0]

		updateReq := &dto.UtilityUpdateReq{
			ID: util.ID, WaterUsage: 12.0, ElectricUsage: 180.0,
		}
		err := utilitySvc.Update(updateReq)
		assert.NoError(t, err)
	})

	t.Run("缴费", func(t *testing.T) {
		listRes, _ := utilitySvc.GetListByPage(&dto.UtilityListPageReq{PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10}})
		utilList := listRes.List.([]*dto.UtilityResult)
		util := utilList[0]

		payReq := &dto.UtilityPayReq{ID: util.ID}
		err := utilitySvc.Pay(payReq)
		assert.NoError(t, err)

		detail, _ := utilitySvc.GetDetail(util.ID)
		assert.Equal(t, 2, detail.PayStatus)
	})

	t.Run("已缴费记录不能修改", func(t *testing.T) {
		listRes, _ := utilitySvc.GetListByPage(&dto.UtilityListPageReq{PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10}})
		utilList := listRes.List.([]*dto.UtilityResult)
		util := utilList[0]

		updateReq := &dto.UtilityUpdateReq{
			ID: util.ID, WaterUsage: 20.0, ElectricUsage: 200.0,
		}
		err := utilitySvc.Update(updateReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已缴费")
	})

	t.Run("删除已缴费记录失败", func(t *testing.T) {
		listRes, _ := utilitySvc.GetListByPage(&dto.UtilityListPageReq{PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10}})
		utilList := listRes.List.([]*dto.UtilityResult)
		util := utilList[0]

		err := utilitySvc.Delete([]string{util.ID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已缴费")
	})

	t.Run("获取统计信息", func(t *testing.T) {
		stat, err := utilitySvc.GetStatistics("", 2024, 3)
		assert.NoError(t, err)
		assert.NotNil(t, stat)
	})
}
