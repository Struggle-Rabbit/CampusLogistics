package system

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
)

func (s *SystemService) GetOperationLogListByPage(req *dto.OperationLogByPageReq) (*dto.PageResult, error) {
	var total int64
	db := s.app.DB.Model(&model.SysOperationLog{})

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.IP != "" {
		db.Where("ip = ?", req.IP)
	}

	if req.UserID != "" {
		db.Where("user_id = ?", req.UserID)
	}

	if !req.OperationTimeStart.IsZero() && req.OperationTimeEnd.IsZero() {
		db.Where("operation_at >= ? AND operation_at <= ?", req.OperationTimeStart, req.OperationTimeEnd)
	}
	var list []model.SysOperationLog

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:        list,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}
