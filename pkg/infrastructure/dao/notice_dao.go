package dao

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type NoticeDao struct {
	orm *orm.Gorm
}

func NewNoticeDao(orm *orm.Gorm) *NoticeDao {
	return &NoticeDao{orm: orm}
}

func (dao *NoticeDao) CreateNotice(notice *model.Notice, targetIDs []int) (*model.Notice, error) {
	// 开始事务
	tx := dao.orm.Begin()

	// 1. 创建通知
	if err := tx.Create(notice).Error; err != nil {
		tx.Rollback() // 回滚事务
		return nil, fmt.Errorf("failed to create notice: %v", err)
	}

	// 2. 如果 target_type 是 all，则不需要插入任何 NoticeTarget 记录
	if notice.TargetType != "all" && len(targetIDs) > 0 {
		// 3. 插入 NoticeTarget 记录
		for _, targetID := range targetIDs {
			target := &model.NoticeTarget{
				NoticeID: notice.ID,
				TargetID: targetID,
			}
			if err := tx.Create(target).Error; err != nil {
				tx.Rollback() // 回滚事务
				return nil, fmt.Errorf("failed to create notice target: %v", err)
			}
		}
	}

	// 提交事务
	tx.Commit()
	return notice, nil
}

// GetNoticeByID 根据 ID 获取通知
func (dao *NoticeDao) GetNoticeByID(id int) (*model.Notice, error) {
	var notice model.Notice
	err := dao.orm.First(&notice, id).Error
	return &notice, err
}

func (dao *NoticeDao) MarkUserNoticeAsRead(ctx context.Context, userID, noticeID int) error {
	readLog := model.NoticeReadLog{
		UserID:   userID,
		NoticeID: noticeID,
		ReadAt:   time.Now(),
	}
	// 忽略重复插入错误
	return dao.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "notice_id"}},
		DoNothing: true,
	}).WithContext(ctx).Create(&readLog).Error
}

// DeleteNotice 根据 ID 删除通知
func (dao *NoticeDao) DeleteNotice(id int) error {
	return dao.orm.Transaction(func(tx *gorm.DB) error {
		// 删除 NoticeTarget 记录
		if err := tx.Where("notice_id = ?", id).Delete(&model.NoticeTarget{}).Error; err != nil {
			return fmt.Errorf("failed to delete NoticeTarget: %w", err)
		}

		if err := tx.Where("notice_id = ?", id).Delete(&model.NoticeReadLog{}).Error; err != nil {
			return fmt.Errorf("failed to delete  NoticeReadLog: %w", err)
		}

		// 删除 Notice 记录
		if err := tx.Delete(&model.Notice{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete notice: %w", err)
		}

		return nil
	})
}

// GetUserNoticesWithReadStatus 获取用户可见的通知列表并附带阅读状态
// 返回: 通知列表, 总数, 错误
// readType 0：全部 1：未读 2 已读
func (dao *NoticeDao) GetUserNoticesWithReadStatus(
	userID int, roleIDs []int, readType int, title string, page, pageSize int) (notices []model.Notice,
	total int64, err error) {

	var readNoticeIDs []int

	if err = dao.orm.Model(&model.NoticeReadLog{}).
		Where("user_id = ?", userID).
		Pluck("notice_id", &readNoticeIDs).Error; err != nil {
		return
	}

	readMap := make(map[int]bool)
	for _, id := range readNoticeIDs {
		readMap[id] = true
	}

	// 创建基础查询
	baseQuery := dao.orm.Model(&model.Notice{}).
		Select("notices.*").
		Joins("LEFT JOIN notice_targets ON notice_targets.notice_id = notices.id").
		Where("(notices.target_type = 'all') OR "+
			"(notices.target_type = 'roles' AND notice_targets.target_id IN ?) OR "+
			"(notices.target_type = 'users' AND notice_targets.target_id = ?)",
			roleIDs, userID)

	countQuery := dao.orm.Model(&model.Notice{}).
		Joins("LEFT JOIN notice_targets ON notice_targets.notice_id = notices.id").
		Where("(notices.target_type = 'all') OR "+
			"(notices.target_type = 'roles' AND notice_targets.target_id IN ?) OR "+
			"(notices.target_type = 'users' AND notice_targets.target_id = ?)",
			roleIDs, userID)
	if title != "" {
		baseQuery.Where(" title like ? ", "%"+title+"%")
		countQuery.Where(" title like ? ", "%"+title+"%")
	}

	switch readType {
	case 1:
		//已读
		if len(readNoticeIDs) > 0 {
			baseQuery.Where(" notices.id  in ? ", readNoticeIDs)
			countQuery.Where(" notices.id  in ? ", readNoticeIDs)
		} else {
			baseQuery.Where(" notices.id  in ? ", []int{-1})
			countQuery.Where(" notices.id  in ? ", []int{-1})
		}

	case 2:
		//未读
		if len(readNoticeIDs) > 0 {
			baseQuery.Where(" notices.id not in ? ", readNoticeIDs)
			countQuery.Where(" notices.id not in ? ", readNoticeIDs)
		}

	}

	// 1. 查询总数
	if err = countQuery.Count(&total).Error; err != nil {
		return
	}

	// 2. 查询分页数据
	query := baseQuery.
		Order("notices.id DESC").
		Scopes(Paginate(page, pageSize))

	if err = query.Find(&notices).Error; err != nil {
		return
	}

	for i := range notices {
		notices[i].IsRead = readMap[notices[i].ID]
	}

	return
}

// Paginate 分页scope
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (dao *NoticeDao) Truncate() {
	eg := errgroup.Group{}
	eg.Go(func() error {
		dao.orm.Migrator().DropTable(&model.Notice{})
		dao.orm.Migrator().CreateTable(&model.Notice{})
		return nil
	})
	eg.Go(func() error {
		dao.orm.Migrator().DropTable(&model.NoticeReadLog{})
		dao.orm.Migrator().CreateTable(&model.NoticeReadLog{})
		return nil
	})
	eg.Go(func() error {
		dao.orm.Migrator().DropTable(&model.NoticeTarget{})
		dao.orm.Migrator().CreateTable(&model.NoticeTarget{})
		return nil
	})
	eg.Wait()
}
