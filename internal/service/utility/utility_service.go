package utility

import (
	"errors"
	"math"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewUtilityService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(req *dto.UtilityCreateReq) error {
	if err := s.checkRoomExists(req.RoomID); err != nil {
		return err
	}

	var existCount int64
	s.db.Model(&model.DormUtility{}).Where("room_id = ? AND year = ? AND month = ?", req.RoomID, req.Year, req.Month).Count(&existCount)
	if existCount > 0 {
		return errors.New("该宿舍本月水电记录已存在")
	}

	amount := s.calculateAmount(req.WaterUsage, req.ElectricUsage)

	utility := &model.DormUtility{
		RoomID:        req.RoomID,
		Year:          req.Year,
		Month:         req.Month,
		WaterUsage:    req.WaterUsage,
		ElectricUsage: req.ElectricUsage,
		Amount:        amount,
		PayStatus:     1,
	}
	return s.db.Create(utility).Error
}

func (s *Service) Update(req *dto.UtilityUpdateReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var utility model.DormUtility
		if err := tx.First(&utility, "id = ?", req.ID).Error; err != nil {
			return errors.New("水电费记录不存在")
		}

		if utility.PayStatus == 2 {
			return errors.New("已缴费记录不允许修改")
		}

		amount := s.calculateAmount(req.WaterUsage, req.ElectricUsage)
		utility.WaterUsage = req.WaterUsage
		utility.ElectricUsage = req.ElectricUsage
		utility.Amount = amount

		return tx.Save(&utility).Error
	})
}

func (s *Service) Delete(ids []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var utility model.DormUtility
			if err := tx.First(&utility, "id = ?", id).Error; err != nil {
				return errors.New("水电费记录不存在")
			}
			if utility.PayStatus == 2 {
				return errors.New("已缴费记录不允许删除")
			}
		}
		return tx.Delete(&model.DormUtility{}, "id IN ?", ids).Error
	})
}

func (s *Service) GetListByPage(req *dto.UtilityListPageReq) (*dto.PageResult, error) {
	var list []*model.DormUtility
	var total int64
	db := s.db.Model(&model.DormUtility{})

	if req.RoomID != "" {
		db = db.Where("room_id = ?", req.RoomID)
	}
	if req.Year != 0 {
		db = db.Where("year = ?", req.Year)
	}
	if req.Month != 0 {
		db = db.Where("month = ?", req.Month)
	}
	if req.PayStatus != 0 {
		db = db.Where("pay_status = ?", req.PayStatus)
	}

	if req.BuildingID != "" {
		db = db.Where("room_id IN ?",
			s.db.Model(&model.DormRoom{}).Select("id").Where("building_id = ?", req.BuildingID))
	}
	if req.CampusID != "" {
		db = db.Where("room_id IN ?",
			s.db.Model(&model.DormRoom{}).Select("id").Where("building_id IN ?",
				s.db.Model(&model.Building{}).Select("id").Where("campus_id = ?", req.CampusID)))
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Order("year DESC, month DESC").Find(&list).Error; err != nil {
		return nil, err
	}

	price := s.getPriceConfig()

	var results []*dto.UtilityResult
	for _, v := range list {
		room, _ := s.getRoomInfo(v.RoomID)
		roomNo := ""
		buildingID := ""
		buildingName := ""
		campusID := ""
		campusName := ""
		if room != nil {
			roomNo = room.RoomNo
			buildingID = room.BuildingID
			building, _ := s.getBuildingInfo(room.BuildingID)
			if building != nil {
				buildingName = building.BuildingName
				campus, _ := s.getCampusInfo(building.CampusID)
				if campus != nil {
					campusID = campus.ID
					campusName = campus.CampusName
				}
			}
		}

		waterAmount := roundToTwoDecimals(v.WaterUsage * price.WaterPrice)
		electricAmount := roundToTwoDecimals(v.ElectricUsage * price.ElectricPrice)

		results = append(results, &dto.UtilityResult{
			ID:             v.ID,
			RoomID:         v.RoomID,
			RoomNo:         roomNo,
			BuildingID:     buildingID,
			BuildingName:   buildingName,
			CampusID:       campusID,
			CampusName:     campusName,
			Year:           v.Year,
			Month:          v.Month,
			WaterUsage:     v.WaterUsage,
			ElectricUsage:  v.ElectricUsage,
			WaterPrice:     price.WaterPrice,
			ElectricPrice:  price.ElectricPrice,
			WaterAmount:    waterAmount,
			ElectricAmount: electricAmount,
			Amount:         v.Amount,
			PayStatus:      v.PayStatus,
			PayStatusName:  getPayStatusName(v.PayStatus),
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *Service) GetDetail(id string) (*dto.UtilityResult, error) {
	var utility model.DormUtility
	if err := s.db.First(&utility, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("水电费记录不存在")
		}
		return nil, err
	}

	price := s.getPriceConfig()

	room, _ := s.getRoomInfo(utility.RoomID)
	roomNo := ""
	buildingID := ""
	buildingName := ""
	campusID := ""
	campusName := ""
	if room != nil {
		roomNo = room.RoomNo
		buildingID = room.BuildingID
		building, _ := s.getBuildingInfo(room.BuildingID)
		if building != nil {
			buildingName = building.BuildingName
			campus, _ := s.getCampusInfo(building.CampusID)
			if campus != nil {
				campusID = campus.ID
				campusName = campus.CampusName
			}
		}
	}

	waterAmount := roundToTwoDecimals(utility.WaterUsage * price.WaterPrice)
	electricAmount := roundToTwoDecimals(utility.ElectricUsage * price.ElectricPrice)

	return &dto.UtilityResult{
		ID:             utility.ID,
		RoomID:         utility.RoomID,
		RoomNo:         roomNo,
		BuildingID:     buildingID,
		BuildingName:   buildingName,
		CampusID:       campusID,
		CampusName:     campusName,
		Year:           utility.Year,
		Month:          utility.Month,
		WaterUsage:     utility.WaterUsage,
		ElectricUsage:  utility.ElectricUsage,
		WaterPrice:     price.WaterPrice,
		ElectricPrice:  price.ElectricPrice,
		WaterAmount:    waterAmount,
		ElectricAmount: electricAmount,
		Amount:         utility.Amount,
		PayStatus:      utility.PayStatus,
		PayStatusName:  getPayStatusName(utility.PayStatus),
	}, nil
}

func (s *Service) Pay(req *dto.UtilityPayReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var utility model.DormUtility
		if err := tx.First(&utility, "id = ?", req.ID).Error; err != nil {
			return errors.New("水电费记录不存在")
		}
		if utility.PayStatus == 2 {
			return errors.New("该记录已缴费")
		}
		return tx.Model(&utility).Update("pay_status", 2).Error
	})
}

func (s *Service) BatchPay(req *dto.UtilityBatchPayReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var utility model.DormUtility
			if err := tx.First(&utility, "id = ?", id).Error; err != nil {
				return errors.New("水电费记录不存在: " + id)
			}
			if utility.PayStatus == 2 {
				continue
			}
			if err := tx.Model(&utility).Update("pay_status", 2).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) ImportData(reqs []*dto.UtilityImportReq) (int, error) {
	var successCount int

	for _, req := range reqs {
		var existCount int64
		s.db.Model(&model.DormUtility{}).Where("room_id = ? AND year = ? AND month = ?", req.RoomID, req.Year, req.Month).Count(&existCount)
		if existCount > 0 {
			continue
		}

		amount := s.calculateAmount(req.WaterUsage, req.ElectricUsage)

		utility := &model.DormUtility{
			RoomID:        req.RoomID,
			Year:          req.Year,
			Month:         req.Month,
			WaterUsage:    req.WaterUsage,
			ElectricUsage: req.ElectricUsage,
			Amount:        amount,
			PayStatus:     1,
		}
		if err := s.db.Create(utility).Error; err == nil {
			successCount++
		}
	}

	return successCount, nil
}

func (s *Service) UpdatePrice(req *dto.UtilityPriceReq) error {
	var price model.UtilityPrice
	s.db.FirstOrCreate(&price, model.UtilityPrice{})

	return s.db.Model(&price).Updates(map[string]interface{}{
		"water_price":    req.WaterPrice,
		"electric_price": req.ElectricPrice,
	}).Error
}

func (s *Service) GetPrice() (*dto.UtilityPriceResult, error) {
	var price model.UtilityPrice
	if err := s.db.First(&price).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			price = model.UtilityPrice{
				WaterPrice:    3.5,
				ElectricPrice: 0.6,
			}
			s.db.Create(&price)
		} else {
			return nil, err
		}
	}

	return &dto.UtilityPriceResult{
		WaterPrice:    price.WaterPrice,
		ElectricPrice: price.ElectricPrice,
	}, nil
}

func (s *Service) GetStatistics(campusID string, year, month int) (*dto.UtilityStatResult, error) {
	var results []struct {
		TotalWaterUsage    float64
		TotalElectricUsage float64
		TotalAmount        float64
		UnpaidCount        int64
		UnpaidAmount       float64
	}

	db := s.db.Model(&model.DormUtility{})

	if campusID != "" {
		db = db.Where("room_id IN ?",
			s.db.Model(&model.DormRoom{}).Select("id").Where("building_id IN ?",
				s.db.Model(&model.Building{}).Select("id").Where("campus_id = ?", campusID)))
	}
	if year != 0 {
		db = db.Where("year = ?", year)
	}
	if month != 0 {
		db = db.Where("month = ?", month)
	}

	db.Select("SUM(water_usage) as total_water_usage, SUM(electric_usage) as total_electric_usage, SUM(amount) as total_amount, SUM(CASE WHEN pay_status = 1 THEN 1 ELSE 0 END) as unpaid_count, SUM(CASE WHEN pay_status = 1 THEN amount ELSE 0 END) as unpaid_amount").Scan(&results)

	if len(results) == 0 {
		return &dto.UtilityStatResult{}, nil
	}

	return &dto.UtilityStatResult{
		TotalWaterUsage:    results[0].TotalWaterUsage,
		TotalElectricUsage: results[0].TotalElectricUsage,
		TotalAmount:        results[0].TotalAmount,
		UnpaidCount:        int(results[0].UnpaidCount),
		UnpaidAmount:       results[0].UnpaidAmount,
	}, nil
}

func (s *Service) GetUnpaidWarning() ([]*dto.UtilityResult, error) {
	var currentYear = time.Now().Year()
	var currentMonth = int(time.Now().Month())

	var list []*model.DormUtility
	if err := s.db.Where("year = ? AND month < ? AND pay_status = ?", currentYear, currentMonth+2, 1).Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.UtilityResult
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

		results = append(results, &dto.UtilityResult{
			ID:           v.ID,
			RoomID:       v.RoomID,
			RoomNo:       roomNo,
			BuildingName: buildingName,
			CampusName:   campusName,
			Year:         v.Year,
			Month:        v.Month,
			Amount:       v.Amount,
			PayStatus:    v.PayStatus,
		})
	}

	return results, nil
}

func (s *Service) checkRoomExists(roomID string) error {
	var count int64
	s.db.Model(&model.DormRoom{}).Where("id = ?", roomID).Count(&count)
	if count == 0 {
		return errors.New("宿舍信息不存在")
	}
	return nil
}

func (s *Service) calculateAmount(waterUsage, electricUsage float64) float64 {
	price := s.getPriceConfig()
	amount := waterUsage*price.WaterPrice + electricUsage*price.ElectricPrice
	return roundToTwoDecimals(amount)
}

func (s *Service) getPriceConfig() *model.UtilityPrice {
	var price model.UtilityPrice
	if err := s.db.First(&price).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			price = model.UtilityPrice{
				WaterPrice:    3.5,
				ElectricPrice: 0.6,
			}
		}
	}
	return &price
}

func getPayStatusName(status int) string {
	switch status {
	case 1:
		return "未缴费"
	case 2:
		return "已缴费"
	default:
		return "未知"
	}
}

func (s *Service) GetUserDormUtility(userID string, year, month int) (*dto.PageResult, error) {
	var dormUser model.DormUser
	if err := s.db.Where("user_id = ? AND status = ?", userID, 1).First(&dormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到用户的住宿信息")
		}
		return nil, err
	}

	var utilityList []*model.DormUtility
	var total int64

	db := s.db.Model(&model.DormUtility{}).Where("room_id = ?", dormUser.RoomID)
	if year != 0 {
		db = db.Where("year = ?", year)
	}
	if month != 0 {
		db = db.Where("month = ?", month)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Order("year DESC, month DESC").Find(&utilityList).Error; err != nil {
		return nil, err
	}

	price := s.getPriceConfig()
	var results []*dto.UtilityResult
	for _, v := range utilityList {
		results = append(results, &dto.UtilityResult{
			ID:             v.ID,
			RoomID:         v.RoomID,
			Year:           v.Year,
			Month:          v.Month,
			WaterUsage:     v.WaterUsage,
			ElectricUsage:  v.ElectricUsage,
			WaterPrice:     price.WaterPrice,
			ElectricPrice:  price.ElectricPrice,
			WaterAmount:    roundToTwoDecimals(v.WaterUsage * price.WaterPrice),
			ElectricAmount: roundToTwoDecimals(v.ElectricUsage * price.ElectricPrice),
			Amount:         v.Amount,
			PayStatus:      v.PayStatus,
			PayStatusName:  getPayStatusName(v.PayStatus),
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: 1,
		PageSize:    int(total),
	}, nil
}

func (s *Service) getRoomInfo(id string) (*model.DormRoom, error) {
	var room model.DormRoom
	if err := s.db.First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *Service) getBuildingInfo(id string) (*model.Building, error) {
	var building model.Building
	if err := s.db.First(&building, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &building, nil
}

func (s *Service) getCampusInfo(id string) (*model.Campus, error) {
	var campus model.Campus
	if err := s.db.First(&campus, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &campus, nil
}

func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}
