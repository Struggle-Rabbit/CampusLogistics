package repair

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
)

type RepairService struct {
	app *app.App
}

func NewRepairService(app *app.App) *RepairService {
	return &RepairService{
		app: app,
	}
}

// GenerateOrderNo生成唯一报修单号
func (RepairService) GenerateOrderNo(prefix string) string {
	// 1. 获取当前时间，精确到秒 (格式: 20060102150405)
	timestamp := time.Now().Format("20060102150405")

	// 2. 生成6位随机数字 (区间: 100000 - 999999)
	// 用 rand.Intn(900000) 得到 0-899999，再加 100000 保证是6位
	randomPart := rand.Intn(900000) + 100000

	// 3. 拼接结果
	return fmt.Sprintf("%s%s%d", prefix, timestamp, randomPart)
}

func (s *RepairService) RepairOrderSubmit(userID string, req *dto.RepairOrderSubmitReq) error {
	if err := s.app.DB.Create(&model.RepairOrder{
		OrderNo:     s.GenerateOrderNo("RO"),
		UserID:      userID,
		RepairType:  req.RepairType,
		Address:     req.Address,
		Description: req.Description,
		Images:      req.Images,
		Contact:     req.Contact,
		Phone:       req.Phone,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (s *RepairService) GetListByPage(req *dto.RepairOrderListByPageReq) (*dto.PageResult, error) {
	var total int64
	var repairRes []*model.RepairOrder

	db := s.app.DB.Model(&model.RepairOrder{})

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if req.OrderNo != "" {
		db.Where("order_no = ?", req.OrderNo)
	}
	if req.Contact != "" {
		db.Where("contact = ?", req.Contact)
	}
	if req.HandlerID != "" {
		db.Where("handler_id = ?", req.HandlerID)
	}
	if req.Phone != "" {
		db.Where("phone = ?", req.Phone)
	}
	if req.RepairType != 0 {
		db.Where("repair_type = ?", req.RepairType)
	}
	if req.Status != 0 {
		db.Where("status = ?", req.Status)
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&repairRes).Error; err != nil {
		return nil, err
	}

	var dtoList []*dto.RepairOrderResult

	for _, v := range repairRes {
		dtoList = append(dtoList, &dto.RepairOrderResult{
			ID:          v.ID,
			OrderNo:     v.OrderNo,
			UserID:      v.UserID,
			RepairType:  v.RepairType,
			Address:     v.Address,
			Description: v.Description,
			Images:      v.Images,
			Contact:     v.Contact,
			Phone:       v.Phone,
			Status:      v.Status,
			HandlerID:   v.HandlerID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return &dto.PageResult{
		List:        dtoList,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *RepairService) GetDetailById(id string) (*dto.RepairOrderResult, error) {
	var repairOrder model.RepairOrder
	if err := s.app.DB.Model(&model.RepairOrder{}).Where("id = ?", id).First(&repairOrder).Error; err != nil {
		return nil, err
	}
	return &dto.RepairOrderResult{
		ID:          repairOrder.ID,
		CreatedAt:   repairOrder.CreatedAt,
		UpdatedAt:   repairOrder.UpdatedAt,
		OrderNo:     repairOrder.OrderNo,
		UserID:      repairOrder.UserID,
		RepairType:  repairOrder.RepairType,
		Address:     repairOrder.Address,
		Description: repairOrder.Description,
		Images:      repairOrder.Images,
		Contact:     repairOrder.Contact,
		Phone:       repairOrder.Phone,
		Status:      repairOrder.Status,
		HandlerID:   repairOrder.HandlerID,
	}, nil
}

func (s *RepairService) DelRepairOrderById(id string) error {
	if err := s.app.DB.Delete(&model.RepairOrder{}, id).Error; err != nil {
		return err
	}
	return nil
}

// func (s *RepairService) UpdateRepairOrder() error {
// 	if err := s.app.DB.Delete(&model.RepairOrder{}, id).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
