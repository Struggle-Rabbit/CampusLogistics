package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/campus"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/dorm"
	"github.com/stretchr/testify/assert"
)

func TestDormService(t *testing.T) {
	_, appInstance := SetupTestDB()
	campusSvc := campus.NewCampusService(appInstance)
	buildingSvc := building.NewBuildingService(appInstance)
	dormSvc := dorm.NewDormService(appInstance)

	campusReq := &dto.CampusCreateReq{
		CampusName: "测试校区",
		Address:    "测试地址",
	}
	campusSvc.Create(campusReq)

	campusList, _ := campusSvc.GetListByPage(&dto.CampusListPageReq{
		PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
	})
	testCampus := campusList.List.([]*dto.CampusResult)[0]

	buildingReq := &dto.BuildingCreateReq{
		CampusID:     testCampus.ID,
		BuildingNo:   "D01",
		BuildingName: "1号宿舍楼",
		FloorCount:   6,
		RoomCount:    60,
	}
	buildingSvc.Create(buildingReq)

	buildingList, _ := buildingSvc.GetBuildingsByCampus(testCampus.ID)
	testBuilding := buildingList[0]

	t.Run("创建宿舍", func(t *testing.T) {
		req := &dto.DormCreateReq{
			BuildingID: testBuilding.ID,
			RoomNo:     "101",
			Floor:      1,
			RoomType:   dto.DormRoomType4,
			MaxCount:   4,
		}
		err := dormSvc.Create(req)
		assert.NoError(t, err)

		res, err := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
	})

	t.Run("宿舍编号唯一性校验", func(t *testing.T) {
		req := &dto.DormCreateReq{
			BuildingID: testBuilding.ID,
			RoomNo:     "101",
			Floor:      1,
			RoomType:   dto.DormRoomType4,
			MaxCount:   4,
		}
		err := dormSvc.Create(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "该楼栋下宿舍编号已存在")
	})

	t.Run("分配宿舍", func(t *testing.T) {
		listRes, _ := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		dormList := listRes.List.([]*dto.DormResult)
		testRoom := dormList[0]

		assignReq := &dto.DormAssignReq{
			RoomID: testRoom.ID,
			UserID: "STUDENT001",
		}
		err := dormSvc.AssignDorm(assignReq)
		assert.NoError(t, err)

		detail, _ := dormSvc.GetDetail(testRoom.ID)
		assert.Equal(t, 1, detail.CurrentCount)
	})

	t.Run("同一用户不能重复分配宿舍", func(t *testing.T) {
		listRes, _ := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		dormList := listRes.List.([]*dto.DormResult)
		testRoom := dormList[0]

		assignReq := &dto.DormAssignReq{
			RoomID: testRoom.ID,
			UserID: "STUDENT001",
		}
		err := dormSvc.AssignDorm(assignReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "该用户已在其他宿舍入住")
	})

	t.Run("退宿", func(t *testing.T) {
		listRes, _ := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		dormList := listRes.List.([]*dto.DormResult)
		testRoom := dormList[0]

		checkoutReq := &dto.DormCheckOutReq{
			RoomID: testRoom.ID,
			UserID: "STUDENT001",
		}
		err := dormSvc.CheckOut(checkoutReq)
		assert.NoError(t, err)

		detail, _ := dormSvc.GetDetail(testRoom.ID)
		assert.Equal(t, 0, detail.CurrentCount)
	})

	t.Run("宿舍容量预警", func(t *testing.T) {
		listRes, _ := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		dormList := listRes.List.([]*dto.DormResult)
		testRoom := dormList[0]

		dormSvc.AssignDorm(&dto.DormAssignReq{
			RoomID: testRoom.ID,
			UserID: "STUDENT002",
		})
		dormSvc.AssignDorm(&dto.DormAssignReq{
			RoomID: testRoom.ID,
			UserID: "STUDENT003",
		})
		dormSvc.AssignDorm(&dto.DormAssignReq{
			RoomID: testRoom.ID,
			UserID: "STUDENT004",
		})

		detail, _ := dormSvc.GetDetail(testRoom.ID)
		assert.Equal(t, 3, detail.CurrentCount)
		assert.GreaterOrEqual(t, detail.FillRate, 75.0)
	})

	t.Run("删除有学生的宿舍失败", func(t *testing.T) {
		listRes, _ := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		dormList := listRes.List.([]*dto.DormResult)
		testRoom := dormList[0]

		err := dormSvc.Delete([]string{testRoom.ID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "宿舍仍有学生入住")
	})

	t.Run("退宿后删除成功", func(t *testing.T) {
		listRes, _ := dormSvc.GetListByPage(&dto.DormListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		dormList := listRes.List.([]*dto.DormResult)
		testRoom := dormList[0]

		for i := 2; i <= 4; i++ {
			dormSvc.CheckOut(&dto.DormCheckOutReq{
				RoomID: testRoom.ID,
				UserID: "STUDENT00" + string(rune('0'+i)),
			})
		}

		err := dormSvc.Delete([]string{testRoom.ID})
		assert.NoError(t, err)
	})
}
