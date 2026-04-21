package system

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/go-viper/mapstructure/v2"
)

func (s *SystemService) GetOperationLogListByPage(req *dto.OperationLogByPageReq) (*dto.PageResult, error) {
	var operationLogReq model.SysOperationLog
	var total int64
	if err := s.app.DB.Model(&model.SysOperationLog{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(req, &operationLogReq); err != nil {
		return nil, err
	}
	var list []*dto.OperationLogResult

	if err := s.app.DB.Model(&model.SysOperationLog{}).Where(&operationLogReq).Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:        list,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}
