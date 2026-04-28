package unittest

import (
	"testing"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
	"github.com/stretchr/testify/assert"
)

func TestNoticeService(t *testing.T) {
	_, appInstance := SetupTestDB()
	svc := notice.NewNoticeService(appInstance)

	creatorID := "ADMIN001"
	now := time.Now()

	t.Run("创建公告", func(t *testing.T) {
		req := &dto.NoticeCreateReq{
			Title:       "测试公告",
			Content:     "测试内容",
			NoticeType:  dto.NoticeTypeCampus,
			IsTop:       dto.IsTopNo,
			PublishTime: now,
		}
		err := svc.Create(creatorID, req)
		assert.NoError(t, err)
	})

	t.Run("查询公告列表", func(t *testing.T) {
		res, err := svc.GetListByPage(&dto.NoticeListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, res.Total, int64(1))
	})

	t.Run("更新公告", func(t *testing.T) {
		res, _ := svc.GetListByPage(&dto.NoticeListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		noticeList := res.List.([]*dto.NoticeResult)
		testNotice := noticeList[0]

		updateReq := &dto.NoticeUpdateReq{
			ID:         testNotice.ID,
			Title:      "更新后的标题",
			Content:    "更新后的内容",
			NoticeType: dto.NoticeTypeCampus,
			IsTop:      dto.IsTopNo,
		}
		err := svc.Update(updateReq)
		assert.NoError(t, err)

		detail, _ := svc.GetDetail(testNotice.ID)
		assert.Equal(t, "更新后的标题", detail.Title)
	})

	t.Run("删除公告", func(t *testing.T) {
		req := &dto.NoticeCreateReq{
			Title:       "待删除公告",
			Content:     "内容",
			NoticeType:  dto.NoticeTypeLogistics,
			IsTop:       dto.IsTopNo,
			PublishTime: now,
		}
		svc.Create(creatorID, req)

		res, _ := svc.GetListByPage(&dto.NoticeListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		noticeList := res.List.([]*dto.NoticeResult)
		latestNotice := noticeList[0]

		err := svc.Delete([]string{latestNotice.ID})
		assert.NoError(t, err)
	})

	t.Run("公告浏览量", func(t *testing.T) {
		req := &dto.NoticeCreateReq{
			Title:       "浏览量测试",
			Content:     "内容",
			NoticeType:  dto.NoticeTypeLogistics,
			IsTop:       dto.IsTopNo,
			PublishTime: now,
		}
		svc.Create(creatorID, req)

		res, _ := svc.GetListByPage(&dto.NoticeListPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		})
		noticeList := res.List.([]*dto.NoticeResult)
		latestNotice := noticeList[0]

		detail, err := svc.GetDetail(latestNotice.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, detail.ViewCount, 0)
	})
}
