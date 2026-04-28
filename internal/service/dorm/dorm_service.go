package dorm

import (
	"errors"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type DormService struct {
	app *app.App
}

func NewDormService(app *app.App) *DormService {
	return &DormService{app: app}
}

func (s *DormService) Create(req *dto.DormCreateReq) error {
	if err := s.checkRoomNoUnique("", req.BuildingID, req.RoomNo); err != nil {
		return err
	}
	if err := s.checkBuildingExists(req.BuildingID); err != nil {
		return err
	}

	room := &model.DormRoom{
		BuildingID:   req.BuildingID,
		RoomNo:       req.RoomNo,
		Floor:        req.Floor,
		RoomType:     req.RoomType,
		MaxCount:     req.MaxCount,
		CurrentCount: 0,
		Remark:       req.Remark,
	}
	return s.app.DB.Create(room).Error
}

func (s *DormService) Update(req *dto.DormUpdateReq) error {
	if err := s.checkRoomNoUnique(req.ID, req.BuildingID, req.RoomNo); err != nil {
		return err
	}
	if err := s.checkBuildingExists(req.BuildingID); err != nil {
		return err
	}

	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var room model.DormRoom
		if err := tx.First(&room, "id = ?", req.ID).Error; err != nil {
			return errors.New("宿舍信息不存在")
		}

		if req.MaxCount < room.CurrentCount {
			return errors.New("最大入住人数不能小于当前入住人数")
		}

		room.BuildingID = req.BuildingID
		room.RoomNo = req.RoomNo
		room.Floor = req.Floor
		room.RoomType = req.RoomType
		room.MaxCount = req.MaxCount
		room.Remark = req.Remark
		return tx.Save(&room).Error
	})
}

func (s *DormService) Delete(ids []string) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var room model.DormRoom
			if err := tx.First(&room, "id = ?", id).Error; err != nil {
				return errors.New("宿舍信息不存在")
			}
			if room.CurrentCount > 0 {
				return errors.New("宿舍仍有学生入住，无法删除")
			}
		}
		return tx.Delete(&model.DormRoom{}, "id IN ?", ids).Error
	})
}

func (s *DormService) GetListByPage(req *dto.DormListPageReq) (*dto.PageResult, error) {
	var list []*model.DormRoom
	var total int64
	db := s.app.DB.Model(&model.DormRoom{})

	if req.BuildingID != "" {
		db = db.Where("building_id = ?", req.BuildingID)
	}
	if req.Floor != 0 {
		db = db.Where("floor = ?", req.Floor)
	}
	if req.RoomType != 0 {
		db = db.Where("room_type = ?", req.RoomType)
	}

	if req.Status != 0 {
		switch req.Status {
		case 1:
			db = db.Where("current_count < max_count")
		case 2:
			db = db.Where("current_count >= max_count")
		case 3:
			db = db.Where("status = ?", 3)
		}
	}

	if req.CampusID != "" {
		db = db.Where("building_id IN ?",
			s.app.DB.Model(&model.Building{}).Select("id").Where("campus_id = ?", req.CampusID))
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.DormResult
	for _, v := range list {
		building, _ := s.getBuildingInfo(v.BuildingID)
		buildingName := ""
		campusID := ""
		campusName := ""
		if building != nil {
			buildingName = building.BuildingName
			campus, _ := s.getCampusInfo(building.CampusID)
			if campus != nil {
				campusID = campus.ID
				campusName = campus.CampusName
			}
		}

		roomTypeName := s.getRoomTypeName(v.RoomType)
		fillRate := 0.0
		if v.MaxCount > 0 {
			fillRate = float64(v.CurrentCount) / float64(v.MaxCount) * 100
		}

		status := 1
		if v.CurrentCount >= v.MaxCount {
			status = 2
		}

		results = append(results, &dto.DormResult{
			ID:           v.ID,
			BuildingID:   v.BuildingID,
			BuildingName: buildingName,
			CampusID:     campusID,
			CampusName:   campusName,
			RoomNo:       v.RoomNo,
			Floor:        v.Floor,
			RoomType:     v.RoomType,
			RoomTypeName: roomTypeName,
			MaxCount:     v.MaxCount,
			CurrentCount: v.CurrentCount,
			AvailableBed: v.MaxCount - v.CurrentCount,
			FillRate:     fillRate,
			Status:       status,
			Remark:       v.Remark,
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *DormService) GetDetail(id string) (*dto.DormResult, error) {
	var room model.DormRoom
	if err := s.app.DB.First(&room, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("宿舍信息不存在")
		}
		return nil, err
	}

	building, _ := s.getBuildingInfo(room.BuildingID)
	buildingName := ""
	campusID := ""
	campusName := ""
	if building != nil {
		buildingName = building.BuildingName
		campus, _ := s.getCampusInfo(building.CampusID)
		if campus != nil {
			campusID = campus.ID
			campusName = campus.CampusName
		}
	}

	roomTypeName := s.getRoomTypeName(room.RoomType)
	fillRate := 0.0
	if room.MaxCount > 0 {
		fillRate = float64(room.CurrentCount) / float64(room.MaxCount) * 100
	}

	status := 1
	if room.CurrentCount >= room.MaxCount {
		status = 2
	}

	return &dto.DormResult{
		ID:           room.ID,
		BuildingID:   room.BuildingID,
		BuildingName: buildingName,
		CampusID:     campusID,
		CampusName:   campusName,
		RoomNo:       room.RoomNo,
		Floor:        room.Floor,
		RoomType:     room.RoomType,
		RoomTypeName: roomTypeName,
		MaxCount:     room.MaxCount,
		CurrentCount: room.CurrentCount,
		AvailableBed: room.MaxCount - room.CurrentCount,
		FillRate:     fillRate,
		Status:       status,
		Remark:       room.Remark,
	}, nil
}

func (s *DormService) AssignDorm(req *dto.DormAssignReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var existCount int64
		if err := tx.Model(&model.DormUser{}).Where("user_id = ? AND status = ?", req.UserID, 1).Count(&existCount).Error; err != nil {
			return err
		}
		if existCount > 0 {
			return errors.New("该用户已在其他宿舍入住")
		}

		var room model.DormRoom
		if err := tx.First(&room, "id = ?", req.RoomID).Error; err != nil {
			return errors.New("宿舍信息不存在")
		}
		if room.CurrentCount >= room.MaxCount {
			return errors.New("宿舍床位已满")
		}

		now := time.Now()
		dormUser := &model.DormUser{
			RoomID:      req.RoomID,
			UserID:      req.UserID,
			CheckInTime: &now,
			Status:      1,
		}
		if err := tx.Create(dormUser).Error; err != nil {
			return err
		}

		return tx.Model(&room).Update("current_count", room.CurrentCount+1).Error
	})
}

func (s *DormService) TransferDorm(req *dto.DormTransferReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var sourceRoom model.DormRoom
		if err := tx.First(&sourceRoom, "id = ?", req.RoomID).Error; err != nil {
			return errors.New("当前宿舍信息不存在")
		}

		var targetRoom model.DormRoom
		if err := tx.First(&targetRoom, "id = ?", req.TargetRoomID).Error; err != nil {
			return errors.New("目标宿舍信息不存在")
		}
		if targetRoom.CurrentCount >= targetRoom.MaxCount {
			return errors.New("目标宿舍床位已满")
		}

		var dormUser model.DormUser
		if err := tx.Where("user_id = ? AND status = ?", req.UserID, 1).First(&dormUser).Error; err != nil {
			return errors.New("未找到该用户的住宿记录")
		}

		now := time.Now()
		dormUser.RoomID = req.TargetRoomID
		dormUser.CheckInTime = &now
		if err := tx.Save(&dormUser).Error; err != nil {
			return err
		}

		if err := tx.Model(&sourceRoom).Update("current_count", sourceRoom.CurrentCount-1).Error; err != nil {
			return err
		}
		return tx.Model(&targetRoom).Update("current_count", targetRoom.CurrentCount+1).Error
	})
}

func (s *DormService) CheckOut(req *dto.DormCheckOutReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var room model.DormRoom
		if err := tx.First(&room, "id = ?", req.RoomID).Error; err != nil {
			return errors.New("宿舍信息不存在")
		}

		var dormUser model.DormUser
		if err := tx.Where("user_id = ? AND room_id = ? AND status = ?", req.UserID, req.RoomID, 1).First(&dormUser).Error; err != nil {
			return errors.New("未找到该用户的住宿记录")
		}

		now := time.Now()
		dormUser.CheckOutTime = &now
		dormUser.Status = 2
		if err := tx.Save(&dormUser).Error; err != nil {
			return err
		}

		newCount := room.CurrentCount - 1
		if newCount < 0 {
			newCount = 0
		}
		return tx.Model(&room).Update("current_count", newCount).Error
	})
}

func (s *DormService) GetDormUsers(req *dto.DormUserListReq) (*dto.PageResult, error) {
	var list []*model.DormUser
	var total int64
	db := s.app.DB.Model(&model.DormUser{})

	if req.RoomID != "" {
		db = db.Where("room_id = ?", req.RoomID)
	}
	if req.UserID != "" {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	if req.UserName != "" {
		db = db.Where("user_id IN ?",
			s.app.DB.Model(&model.SysUser{}).Select("user_code").Where("name LIKE ?", "%"+req.UserName+"%"))
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.DormUserResult
	for _, v := range list {
		room, _ := s.getRoomInfo(v.RoomID)
		roomNo := ""
		buildingName := ""
		campusName := ""
		if room != nil {
			roomNo = room.RoomNo
			building, _ := s.getBuildingInfo(room.BuildingID)
			if building != nil {
				buildingName = building.BuildingName
				campus, _ := s.getCampusInfo(building.CampusID)
				if campus != nil {
					campusName = campus.CampusName
				}
			}
		}

		var user model.SysUser
		userName := ""
		userType := ""
		if err := s.app.DB.Where("user_code = ?", v.UserID).First(&user).Error; err == nil {
			userName = user.Name
			userType = user.UserType
		}

		checkInTime := ""
		if v.CheckInTime != nil {
			checkInTime = v.CheckInTime.Format("2006-01-02 15:04:05")
		}
		checkOutTime := ""
		if v.CheckOutTime != nil {
			checkOutTime = v.CheckOutTime.Format("2006-01-02 15:04:05")
		}

		results = append(results, &dto.DormUserResult{
			ID:           v.ID,
			RoomID:       v.RoomID,
			RoomNo:       roomNo,
			BuildingName: buildingName,
			CampusName:   campusName,
			UserID:       v.UserID,
			UserName:     userName,
			UserType:     userType,
			CheckInTime:  checkInTime,
			CheckOutTime: checkOutTime,
			Status:       v.Status,
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *DormService) GetCapacityWarning() ([]*dto.DormResult, error) {
	var list []*model.DormRoom
	if err := s.app.DB.Where("current_count >= max_count * 0.8").Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.DormResult
	for _, v := range list {
		building, _ := s.getBuildingInfo(v.BuildingID)
		buildingName := ""
		if building != nil {
			buildingName = building.BuildingName
		}

		fillRate := 0.0
		if v.MaxCount > 0 {
			fillRate = float64(v.CurrentCount) / float64(v.MaxCount) * 100
		}

		results = append(results, &dto.DormResult{
			ID:           v.ID,
			BuildingID:   v.BuildingID,
			BuildingName: buildingName,
			RoomNo:       v.RoomNo,
			RoomType:     v.RoomType,
			MaxCount:     v.MaxCount,
			CurrentCount: v.CurrentCount,
			FillRate:     fillRate,
			Status:       2,
		})
	}

	return results, nil
}

func (s *DormService) checkRoomNoUnique(id, buildingID, roomNo string) error {
	var count int64
	query := s.app.DB.Model(&model.DormRoom{}).Where("building_id = ? AND room_no = ?", buildingID, roomNo)
	if id != "" {
		query = query.Where("id != ?", id)
	}
	query.Count(&count)
	if count > 0 {
		return errors.New("该楼栋下宿舍编号已存在")
	}
	return nil
}

func (s *DormService) checkBuildingExists(buildingID string) error {
	var count int64
	s.app.DB.Model(&model.Building{}).Where("id = ?", buildingID).Count(&count)
	if count == 0 {
		return errors.New("楼栋信息不存在")
	}
	return nil
}

func (s *DormService) getBuildingInfo(id string) (*model.Building, error) {
	var building model.Building
	if err := s.app.DB.First(&building, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &building, nil
}

func (s *DormService) getCampusInfo(id string) (*model.Campus, error) {
	var campus model.Campus
	if err := s.app.DB.First(&campus, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &campus, nil
}

func (s *DormService) getRoomInfo(id string) (*model.DormRoom, error) {
	var room model.DormRoom
	if err := s.app.DB.First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *DormService) getRoomTypeName(roomType int) string {
	switch roomType {
	case 1:
		return "4人间"
	case 2:
		return "6人间"
	case 3:
		return "8人间"
	case 4:
		return "其他"
	default:
		return "未知"
	}
}
