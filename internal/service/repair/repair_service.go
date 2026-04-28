package repair

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
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
	var repairRes []*dto.RepairOrderResult

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
	if req.HandlerID != nil {
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

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Scan(&repairRes).Error; err != nil {
		return nil, err
	}

	// var dtoList []*dto.RepairOrderResult

	// for _, v := range repairRes {
	// 	dtoList = append(dtoList, &dto.RepairOrderResult{
	// 		ID:          v.ID,
	// 		OrderNo:     v.OrderNo,
	// 		UserID:      v.UserID,
	// 		RepairType:  v.RepairType,
	// 		Address:     v.Address,
	// 		Description: v.Description,
	// 		Images:      v.Images,
	// 		Contact:     v.Contact,
	// 		Phone:       v.Phone,
	// 		Status:      v.Status,
	// 		HandlerID:   v.HandlerID,
	// 		CreatedAt:   v.CreatedAt,
	// 		UpdatedAt:   v.UpdatedAt,
	// 	})
	// }

	return &dto.PageResult{
		List:        repairRes,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *RepairService) GetDetailById(id string) (*dto.RepairOrderResult, error) {
	var repairOrder dto.RepairOrderResult
	if err := s.app.DB.Model(&model.RepairOrder{}).Where("id = ?", id).First(&repairOrder).Error; err != nil {
		return nil, err
	}
	// &dto.RepairOrderResult{
	// 	ID:          repairOrder.ID,
	// 	CreatedAt:   repairOrder.CreatedAt,
	// 	UpdatedAt:   repairOrder.UpdatedAt,
	// 	OrderNo:     repairOrder.OrderNo,
	// 	UserID:      repairOrder.UserID,
	// 	RepairType:  repairOrder.RepairType,
	// 	Address:     repairOrder.Address,
	// 	Description: repairOrder.Description,
	// 	Images:      repairOrder.Images,
	// 	Contact:     repairOrder.Contact,
	// 	Phone:       repairOrder.Phone,
	// 	Status:      repairOrder.Status,
	// 	HandlerID:   repairOrder.HandlerID,
	// }
	return &repairOrder, nil
}

func (s *RepairService) DelRepairOrderById(id string) error {
	if err := s.app.DB.Delete(&model.RepairOrder{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *RepairService) UpdateRepairOrder(req dto.UpdateRepairOrderSubmitReq) error {
	var order model.RepairOrder
	db := s.app.DB.Model(&model.RepairOrder{}).Where("id = ?", req.ID)

	if err := db.First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到订单记录")
		}
		return err
	}

	if order.Status != 1 {
		return errors.New("只有待分配可编辑")
	}

	return db.Where("id = ?", req.ID).Updates(&model.RepairOrder{
		RepairType:  req.RepairType,
		Address:     req.Address,
		Status:      req.Status,
		Description: req.Description,
		Contact:     req.Contact,
		Phone:       req.Phone,
		Images:      req.Images,
		HandlerID:   req.HandlerID,
	}).Error
}

func (s *RepairService) OrderRecord(req dto.RecordReq) error {
	if req.Status < 1 || req.Status > 6 {
		return errors.New("无效的订单状态")
	}
	tx := s.app.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback() // 任何异常都会自动回滚
	var order model.RepairOrder
	if err := tx.Model(&model.RepairOrder{}).Where("id = ?", req.ID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到订单记录")
		}
		return err
	}
	if order.Status == 4 || order.Status == 5 || order.Status == 6 {
		return errors.New("当前订单不可流转")
	}

	result := tx.Model(&model.RepairOrder{}).
		Where("id = ? AND status = ?", req.ID, order.Status). // 核心：乐观锁条件
		Select("status", "handler_id").
		Updates(&model.RepairOrder{
			Status:    req.Status,
			HandlerID: &req.UserID,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("订单状态更新失败")
	}

	if err := tx.Create(&model.RepairRecord{
		OrderID:    req.ID,
		OperatorID: req.UserID,
		OldStatus:  order.Status,
		NewStatus:  req.Status,
		Remark:     req.Remark,
	}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
