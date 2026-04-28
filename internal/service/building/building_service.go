package building

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type BuildingService struct {
	app *app.App
}

func NewBuildingService(app *app.App) *BuildingService {
	return &BuildingService{app: app}
}

func (s *BuildingService) Create(req *dto.BuildingCreateReq) error {
	if err := s.checkBuildingNoUnique("", req.BuildingNo); err != nil {
		return err
	}
	if err := s.checkCampusExists(req.CampusID); err != nil {
		return err
	}

	building := &model.Building{
		CampusID:     req.CampusID,
		BuildingNo:   req.BuildingNo,
		BuildingName: req.BuildingName,
		FloorCount:   req.FloorCount,
		RoomCount:    req.RoomCount,
		Remark:       req.Remark,
	}
	return s.app.DB.Create(building).Error
}

func (s *BuildingService) Update(req *dto.BuildingUpdateReq) error {
	if err := s.checkBuildingNoUnique(req.ID, req.BuildingNo); err != nil {
		return err
	}
	if err := s.checkCampusExists(req.CampusID); err != nil {
		return err
	}

	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var building model.Building
		if err := tx.First(&building, "id = ?", req.ID).Error; err != nil {
			return errors.New("楼栋信息不存在")
		}
		building.CampusID = req.CampusID
		building.BuildingNo = req.BuildingNo
		building.BuildingName = req.BuildingName
		building.FloorCount = req.FloorCount
		building.RoomCount = req.RoomCount
		building.Remark = req.Remark
		return tx.Save(&building).Error
	})
}

func (s *BuildingService) Delete(ids []string) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.DormRoom{}).Where("building_id IN ?", ids).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该楼栋下存在宿舍信息，无法删除")
		}
		return tx.Delete(&model.Building{}, "id IN ?", ids).Error
	})
}

func (s *BuildingService) GetListByPage(req *dto.BuildingListPageReq) (*dto.PageResult, error) {
	var list []*model.Building
	var total int64
	db := s.app.DB.Model(&model.Building{})

	if req.CampusID != "" {
		db = db.Where("campus_id = ?", req.CampusID)
	}
	if req.BuildingNo != "" {
		db = db.Where("building_no LIKE ?", "%"+req.BuildingNo+"%")
	}
	if req.BuildingName != "" {
		db = db.Where("building_name LIKE ?", "%"+req.BuildingName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.BuildingResult
	for _, v := range list {
		var campus model.Campus
		campusName := ""
		if err := s.app.DB.First(&campus, "id = ?", v.CampusID).Error; err == nil {
			campusName = campus.CampusName
		}

		var roomUsedCount int64
		s.app.DB.Model(&model.DormRoom{}).Where("building_id = ?", v.ID).Count(&roomUsedCount)

		results = append(results, &dto.BuildingResult{
			ID:            v.ID,
			CampusID:      v.CampusID,
			CampusName:    campusName,
			BuildingNo:    v.BuildingNo,
			BuildingName:  v.BuildingName,
			FloorCount:    v.FloorCount,
			RoomCount:     v.RoomCount,
			RoomUsedCount: int(roomUsedCount),
			Remark:        v.Remark,
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *BuildingService) GetDetail(id string) (*dto.BuildingResult, error) {
	var building model.Building
	if err := s.app.DB.First(&building, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("楼栋信息不存在")
		}
		return nil, err
	}

	var campus model.Campus
	campusName := ""
	if err := s.app.DB.First(&campus, "id = ?", building.CampusID).Error; err == nil {
		campusName = campus.CampusName
	}

	var roomUsedCount int64
	s.app.DB.Model(&model.DormRoom{}).Where("building_id = ?", building.ID).Count(&roomUsedCount)

	return &dto.BuildingResult{
		ID:            building.ID,
		CampusID:      building.CampusID,
		CampusName:    campusName,
		BuildingNo:    building.BuildingNo,
		BuildingName:  building.BuildingName,
		FloorCount:    building.FloorCount,
		RoomCount:     building.RoomCount,
		RoomUsedCount: int(roomUsedCount),
		Remark:        building.Remark,
	}, nil
}

func (s *BuildingService) ImportBuildings(file *multipart.FileHeader) (int, error) {
	src, err := file.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	reader := csv.NewReader(src)
	var successCount int
	var campusCache = make(map[string]string)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if len(record) < 5 {
			continue
		}

		campusName := strings.TrimSpace(record[0])
		buildingNo := strings.TrimSpace(record[1])
		buildingName := strings.TrimSpace(record[2])
		floorCount, _ := strconv.Atoi(strings.TrimSpace(record[3]))
		roomCount, _ := strconv.Atoi(strings.TrimSpace(record[4]))
		remark := ""
		if len(record) > 5 {
			remark = strings.TrimSpace(record[5])
		}

		campusID, ok := campusCache[campusName]
		if !ok {
			var campus model.Campus
			if err := s.app.DB.Where("campus_name = ?", campusName).First(&campus).Error; err != nil {
				continue
			}
			campusID = campus.ID
			campusCache[campusName] = campusID
		}

		var existCount int64
		s.app.DB.Model(&model.Building{}).Where("building_no = ?", buildingNo).Count(&existCount)
		if existCount > 0 {
			continue
		}

		building := &model.Building{
			CampusID:     campusID,
			BuildingNo:   buildingNo,
			BuildingName: buildingName,
			FloorCount:   floorCount,
			RoomCount:    roomCount,
			Remark:       remark,
		}
		if err := s.app.DB.Create(building).Error; err == nil {
			successCount++
		}
	}

	return successCount, nil
}

func (s *BuildingService) ExportBuildings(req *dto.BuildingExportReq) (string, error) {
	var list []*model.Building
	db := s.app.DB.Model(&model.Building{})

	if req.CampusID != "" {
		db = db.Where("campus_id = ?", req.CampusID)
	}
	if req.BuildingNo != "" {
		db = db.Where("building_no LIKE ?", "%"+req.BuildingNo+"%")
	}

	if err := db.Find(&list).Error; err != nil {
		return "", err
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	header := []string{"校区名称", "楼栋编号", "楼栋名称", "楼层数", "房间数", "备注"}
	writer.Write(header)

	for _, v := range list {
		var campus model.Campus
		campusName := ""
		if err := s.app.DB.First(&campus, "id = ?", v.CampusID).Error; err == nil {
			campusName = campus.CampusName
		}
		record := []string{
			campusName,
			v.BuildingNo,
			v.BuildingName,
			strconv.Itoa(v.FloorCount),
			strconv.Itoa(v.RoomCount),
			v.Remark,
		}
		writer.Write(record)
	}

	writer.Flush()
	return builder.String(), nil
}

func (s *BuildingService) checkBuildingNoUnique(id, buildingNo string) error {
	var count int64
	query := s.app.DB.Model(&model.Building{}).Where("building_no = ?", buildingNo)
	if id != "" {
		query = query.Where("id != ?", id)
	}
	query.Count(&count)
	if count > 0 {
		return errors.New("楼栋编号已存在")
	}
	return nil
}

func (s *BuildingService) checkCampusExists(campusID string) error {
	var count int64
	s.app.DB.Model(&model.Campus{}).Where("id = ?", campusID).Count(&count)
	if count == 0 {
		return errors.New("校区信息不存在")
	}
	return nil
}

func (s *BuildingService) GetBuildingsByCampus(campusID string) ([]*dto.BuildingResult, error) {
	var list []*model.Building
	if err := s.app.DB.Where("campus_id = ?", campusID).Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.BuildingResult
	for _, v := range list {
		results = append(results, &dto.BuildingResult{
			ID:           v.ID,
			CampusID:     v.CampusID,
			BuildingNo:   v.BuildingNo,
			BuildingName: v.BuildingName,
			FloorCount:   v.FloorCount,
			RoomCount:    v.RoomCount,
		})
	}
	return results, nil
}
