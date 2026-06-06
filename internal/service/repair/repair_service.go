package repair

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewRepairService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GenerateOrderNo 生成唯一报修单号
func (s *Service) GenerateOrderNo(prefix string) string {
	return prefix + utils.GenStringID()
}

func (s *Service) RepairOrderSubmit(userID string, req *dto.RepairOrderSubmitReq) error {
	if err := s.db.Create(&model.RepairOrder{
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

func (s *Service) GetListByPage(req *dto.RepairOrderListByPageReq) (*dto.PageResult, error) {
	var total int64
	var repairList []*model.RepairOrder

	db := s.db.Model(&model.RepairOrder{})

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

func (s *Service) GetDetailById(id string) (*dto.RepairOrderResult, error) {
	var order model.RepairOrder
	if err := s.db.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}

	var records []model.RepairRecord
	s.db.Where("order_id = ?", id).Order("created_at DESC").Find(&records)

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

func (s *Service) DelRepairOrderById(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&model.RepairOrder{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("报修单不存在或已被删除")
		}
		if err := tx.Where("order_id = ?", id).Delete(&model.RepairRecord{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *Service) UpdateRepairOrder(req dto.UpdateRepairOrderSubmitReq) error {
	var order model.RepairOrder
	if err := s.db.Where("id = ?", req.ID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到订单记录")
		}
		return err
	}

	if order.Status != 1 {
		return errors.New("只有待分配可编辑")
	}

	return s.db.Model(&order).Updates(map[string]interface{}{
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

func (s *Service) OrderRecord(req dto.RecordReq) error {
	if req.Status < 1 || req.Status > 6 {
		return errors.New("无效的订单状态")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
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
			Where("id = ? AND status = ?", req.ID, order.Status).
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

		return tx.Create(&model.RepairRecord{
			OrderID:    req.ID,
			OperatorID: req.UserID,
			OldStatus:  order.Status,
			NewStatus:  req.Status,
			Remark:     req.Remark,
		}).Error
	})
}
