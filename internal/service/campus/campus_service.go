package campus

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type CampusService struct {
	app *app.App
}

func NewCampusService(app *app.App) *CampusService {
	return &CampusService{app: app}
}

func (s *CampusService) Create(req *dto.CampusCreateReq) error {
	campus := &model.Campus{
		CampusName: req.CampusName,
		Address:    req.Address,
		Contact:    req.Contact,
		Phone:      req.Phone,
		Remark:     req.Remark,
	}
	return s.app.DB.Create(campus).Error
}

func (s *CampusService) Update(req *dto.CampusUpdateReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var campus model.Campus
		if err := tx.First(&campus, "id = ?", req.ID).Error; err != nil {
			return errors.New("校区信息不存在")
		}
		campus.CampusName = req.CampusName
		campus.Address = req.Address
		campus.Contact = req.Contact
		campus.Phone = req.Phone
		campus.Remark = req.Remark
		return tx.Save(&campus).Error
	})
}

func (s *CampusService) Delete(ids []string) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.Building{}).Where("campus_id IN ?", ids).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该校区下存在楼栋信息，无法删除")
		}
		return tx.Delete(&model.Campus{}, "id IN ?", ids).Error
	})
}

func (s *CampusService) GetListByPage(req *dto.CampusListPageReq) (*dto.PageResult, error) {
	var list []*model.Campus
	var total int64
	db := s.app.DB.Model(&model.Campus{})

	if req.CampusName != "" {
		db = db.Where("campus_name LIKE ?", "%"+req.CampusName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.CampusResult
	for _, v := range list {
		var buildingCount int64
		s.app.DB.Model(&model.Building{}).Where("campus_id = ?", v.ID).Count(&buildingCount)

		results = append(results, &dto.CampusResult{
			ID:            v.ID,
			CampusName:    v.CampusName,
			Address:       v.Address,
			Contact:       v.Contact,
			Phone:         v.Phone,
			Remark:        v.Remark,
			BuildingCount: int(buildingCount),
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *CampusService) GetDetail(id string) (*dto.CampusResult, error) {
	var campus model.Campus
	if err := s.app.DB.First(&campus, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("校区信息不存在")
		}
		return nil, err
	}

	var buildingCount int64
	s.app.DB.Model(&model.Building{}).Where("campus_id = ?", campus.ID).Count(&buildingCount)

	return &dto.CampusResult{
		ID:            campus.ID,
		CampusName:    campus.CampusName,
		Address:       campus.Address,
		Contact:       campus.Contact,
		Phone:         campus.Phone,
		Remark:        campus.Remark,
		BuildingCount: int(buildingCount),
	}, nil
}

func (s *CampusService) GetAll() ([]*dto.CampusResult, error) {
	var list []*model.Campus
	if err := s.app.DB.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.CampusResult
	for _, v := range list {
		results = append(results, &dto.CampusResult{
			ID:         v.ID,
			CampusName: v.CampusName,
			Address:    v.Address,
			Contact:    v.Contact,
			Phone:      v.Phone,
		})
	}
	return results, nil
}
