package notice

import (
	"errors"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type NoticeService struct {
	app *app.App
}

func NewNoticeService(app *app.App) *NoticeService {
	return &NoticeService{app: app}
}

func (s *NoticeService) Create(creatorID string, req *dto.NoticeCreateReq) error {
	if req.IsTop == dto.IsTopYes {
		if err := s.checkTopLimit(); err != nil {
			return err
		}
	}

	notice := &model.Notice{
		Title:       req.Title,
		Content:     req.Content,
		NoticeType:  req.NoticeType,
		IsTop:       req.IsTop,
		PublishTime: &req.PublishTime,
		CreatorID:   creatorID,
		Attachments: req.Attachments,
		ViewCount:   0,
	}

	if err := s.app.DB.Create(notice).Error; err != nil {
		return err
	}

	if req.IsTop == dto.IsTopYes {
		return s.reorderTopNotices()
	}

	return nil
}

func (s *NoticeService) Update(req *dto.NoticeUpdateReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var notice model.Notice
		if err := tx.First(&notice, "id = ?", req.ID).Error; err != nil {
			return errors.New("公告信息不存在")
		}

		if req.IsTop == dto.IsTopYes && notice.IsTop != dto.IsTopYes {
			if err := s.checkTopLimit(); err != nil {
				return err
			}
		}

		notice.Title = req.Title
		notice.Content = req.Content
		notice.NoticeType = req.NoticeType
		notice.IsTop = req.IsTop
		if !req.PublishTime.IsZero() {
			notice.PublishTime = &req.PublishTime
		}
		notice.Attachments = req.Attachments

		if err := tx.Save(&notice).Error; err != nil {
			return err
		}

		if req.IsTop == dto.IsTopYes {
			return s.reorderTopNotices()
		}

		return nil
	})
}

func (s *NoticeService) Delete(ids []string) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var notice model.Notice
			if err := tx.First(&notice, "id = ?", id).Error; err != nil {
				return errors.New("公告信息不存在: " + id)
			}
		}
		return tx.Delete(&model.Notice{}, "id IN ?", ids).Error
	})
}

func (s *NoticeService) GetListByPage(req *dto.NoticeListPageReq) (*dto.PageResult, error) {
	var list []*model.Notice
	var total int64
	db := s.app.DB.Model(&model.Notice{})

	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.NoticeType != 0 {
		db = db.Where("notice_type = ?", req.NoticeType)
	}
	if req.IsTop != 0 {
		db = db.Where("is_top = ?", req.IsTop)
	}
	if req.StartTime != "" {
		db = db.Where("publish_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("publish_time <= ?", req.EndTime)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Order("is_top DESC, publish_time DESC").Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.NoticeResult
	for _, v := range list {
		creatorName := s.getCreatorName(v.CreatorID)
		results = append(results, &dto.NoticeResult{
			ID:             v.ID,
			Title:          v.Title,
			Content:        v.Content,
			NoticeType:     v.NoticeType,
			NoticeTypeName: dto.GetNoticeTypeName(v.NoticeType),
			IsTop:          v.IsTop,
			IsTopName:      dto.GetIsTopName(v.IsTop),
			PublishTime:    *v.PublishTime,
			ViewCount:      v.ViewCount,
			CreatorID:      v.CreatorID,
			CreatorName:    creatorName,
			Attachments:    v.Attachments,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *NoticeService) GetPublicList(req *dto.NoticeListPageReq) (*dto.PageResult, error) {
	var list []*model.Notice
	var total int64
	now := time.Now()

	db := s.app.DB.Model(&model.Notice{}).Where("publish_time <= ?", now)

	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.NoticeType != 0 {
		db = db.Where("notice_type = ?", req.NoticeType)
	}
	if req.StartTime != "" {
		db = db.Where("publish_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("publish_time <= ?", req.EndTime)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Order("is_top DESC, publish_time DESC").Find(&list).Error; err != nil {
		return nil, err
	}

	var results []*dto.NoticePublicResult
	for _, v := range list {
		results = append(results, &dto.NoticePublicResult{
			ID:             v.ID,
			Title:          v.Title,
			Content:        v.Content,
			NoticeType:     v.NoticeType,
			NoticeTypeName: dto.GetNoticeTypeName(v.NoticeType),
			PublishTime:    *v.PublishTime,
			ViewCount:      v.ViewCount,
			Attachments:    v.Attachments,
		})
	}

	return &dto.PageResult{
		List:        results,
		Total:       total,
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
	}, nil
}

func (s *NoticeService) GetDetail(id string) (*dto.NoticeResult, error) {
	var notice model.Notice
	if err := s.app.DB.First(&notice, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("公告信息不存在")
		}
		return nil, err
	}

	s.app.DB.Model(&notice).Update("view_count", notice.ViewCount+1)

	creatorName := s.getCreatorName(notice.CreatorID)

	return &dto.NoticeResult{
		ID:             notice.ID,
		Title:          notice.Title,
		Content:        notice.Content,
		NoticeType:     notice.NoticeType,
		NoticeTypeName: dto.GetNoticeTypeName(notice.NoticeType),
		IsTop:          notice.IsTop,
		IsTopName:      dto.GetIsTopName(notice.IsTop),
		PublishTime:    *notice.PublishTime,
		ViewCount:      notice.ViewCount + 1,
		CreatorID:      notice.CreatorID,
		CreatorName:    creatorName,
		Attachments:    notice.Attachments,
		CreatedAt:      notice.CreatedAt,
		UpdatedAt:      notice.UpdatedAt,
	}, nil
}

func (s *NoticeService) SetTop(req *dto.NoticeTopReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		var notice model.Notice
		if err := tx.First(&notice, "id = ?", req.ID).Error; err != nil {
			return errors.New("公告信息不存在")
		}

		if req.IsTop == dto.IsTopYes && notice.IsTop != dto.IsTopYes {
			if err := s.checkTopLimit(); err != nil {
				return err
			}
		}

		if err := tx.Model(&notice).Update("is_top", req.IsTop).Error; err != nil {
			return err
		}

		if req.IsTop == dto.IsTopYes {
			return s.reorderTopNotices()
		}

		return nil
	})
}

func (s *NoticeService) checkTopLimit() error {
	var count int64
	s.app.DB.Model(&model.Notice{}).Where("is_top > ?", 0).Count(&count)
	if count >= dto.MaxTopNotice {
		return errors.New("置顶公告最多3条，请先取消其他公告的置顶状态")
	}
	return nil
}

func (s *NoticeService) reorderTopNotices() error {
	var topNotices []*model.Notice
	if err := s.app.DB.Where("is_top > ?", 0).Order("publish_time DESC").Limit(dto.MaxTopNotice).Find(&topNotices).Error; err != nil {
		return err
	}

	for i, notice := range topNotices {
		priority := dto.MaxTopNotice - i
		if err := s.app.DB.Model(notice).Update("is_top", priority).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *NoticeService) getCreatorName(creatorID string) string {
	var user model.SysUser
	if err := s.app.DB.Where("user_code = ?", creatorID).First(&user).Error; err == nil {
		return user.Name
	}
	return creatorID
}
