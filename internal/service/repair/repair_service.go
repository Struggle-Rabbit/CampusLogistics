package repair

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
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

// GenerateOrderNo 生成唯一报修单号
func (RepairService) GenerateOrderNo(prefix string) string {
	return prefix + utils.GenStringID()
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
	var repairList []*model.RepairOrder

	db := s.app.DB.Model(&model.RepairOrder{})

	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Contact != "" {
		db = db.Where("contact LIKE ?", "%"+req.Contact+"%")
	}
	if req.HandlerID != nil && *req.HandlerID != "" {
		db = db.Where("handler_id = ?", *req.HandlerID)
	}
	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.RepairType != 0 {
		db = db.Where("repair_type = ?", req.RepairType)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.StartTime != "" {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("created_at <= ?", req.EndTime)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Order("created_at DESC").Find(&repairList).Error; err != nil {
		return nil, err
	}

	var dtoList []*dto.RepairOrderResult
	for _, v := range repairList {
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
	var order model.RepairOrder
	if err := s.app.DB.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}

	var records []model.RepairRecord
	s.app.DB.Where("order_id = ?", id).Order("created_at DESC").Find(&records)

	recordDTOs := make([]*dto.RepairRecordResult, len(records))
	for i, r := range records {
		recordDTOs[i] = &dto.RepairRecordResult{
			ID:         r.ID,
			OrderID:    r.OrderID,
			OperatorID: r.OperatorID,
			OldStatus:  r.OldStatus,
			NewStatus:  r.NewStatus,
			Remark:     r.Remark,
			CreatedAt:  r.CreatedAt,
		}
	}

	return &dto.RepairOrderResult{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		UserID:      order.UserID,
		RepairType:  order.RepairType,
		Address:     order.Address,
		Description: order.Description,
		Images:      order.Images,
		Contact:     order.Contact,
		Phone:       order.Phone,
		Status:      order.Status,
		HandlerID:   order.HandlerID,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Records:     recordDTOs,
	}, nil
}

func (s *RepairService) DelRepairOrderById(id string) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		// 软删除报修单
		if err := tx.Delete(&model.RepairOrder{}, "id = ?", id).Error; err != nil {
			return err
		}
		// 级联删除流转记录 (这里也可以选择软删除，但 RepairRecord 目前没加 DeletedAt)
		if err := tx.Where("order_id = ?", id).Delete(&model.RepairRecord{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *RepairService) UpdateRepairOrder(req dto.UpdateRepairOrderSubmitReq) error {
	var order model.RepairOrder
	if err := s.app.DB.Where("id = ?", req.ID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到订单记录")
		}
		return err
	}

	if order.Status != 1 {
		return errors.New("只有待分配可编辑")
	}

	return s.app.DB.Model(&order).Updates(map[string]interface{}{
		"repair_type": req.RepairType,
		"address":     req.Address,
		"status":      req.Status,
		"description": req.Description,
		"contact":     req.Contact,
		"phone":       req.Phone,
		"images":      req.Images,
		"handler_id":  req.HandlerID,
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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order model.RepairOrder
	if err := tx.Model(&model.RepairOrder{}).Where("id = ?", req.ID).First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到订单记录")
		}
		return err
	}
	if order.Status == 4 || order.Status == 5 || order.Status == 6 {
		tx.Rollback()
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
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("订单状态更新失败")
	}

	if err := tx.Create(&model.RepairRecord{
		OrderID:    req.ID,
		OperatorID: req.UserID,
		OldStatus:  order.Status,
		NewStatus:  req.Status,
		Remark:     req.Remark,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
