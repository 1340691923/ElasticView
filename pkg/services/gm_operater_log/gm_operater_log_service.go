package gm_operater_log

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/util"
	"go.uber.org/zap"
)

type GmOperaterLogService struct {
	log *logger.AppLogger
	orm *sqlstore.SqlStore
}

func NewGmOperaterLogService(log *logger.AppLogger, orm *sqlstore.SqlStore) *GmOperaterLogService {
	return &GmOperaterLogService{log: log, orm: orm}
}

func (this *GmOperaterLogService) List(ctx context.Context, reqData dto.GmOperaterLogList) (res []vo.GmOperaterLog, count int64, err error) {

	list := []model.GmOperaterLog{}
	if reqData.Page <= 0 {
		reqData.Page = 1
	}
	if reqData.Limit <= 0 {
		reqData.Limit = 10
	}

	page := reqData.Page
	limit := reqData.Limit

	operater_action := reqData.OperaterAction

	gmOperaterModel := &model.GmOperaterLog{

		OperaterId:     reqData.UserId,
		OperaterAction: operater_action,
	}

	builder := this.orm.Model(gmOperaterModel)

	if reqData.UserId != 0 {
		builder = builder.Where("operater_id = ?", reqData.UserId)
	}

	if reqData.OperaterAction != "" {
		builder = builder.Where("operater_action = ?", reqData.OperaterAction)
	}
	if len(reqData.Date) == 2 {
		builder = builder.Where("created >= ?", reqData.Date[0])
		builder = builder.Where("created <= ?", reqData.Date[1])
	}
	err = builder.Count(&count).Error
	if err != nil {
		return
	}

	err = builder.Order("id desc").Limit(limit).Offset(sqlstore.CreatePage(page, limit)).Scan(&list).Error
	if err != nil {
		return
	}

	for _, v := range list {

		body, err := util.GzipUnCompress(v.Body)
		if err != nil {
			this.log.Error("err", zap.Error(err))
			continue
		}
		res = append(res, vo.GmOperaterLog{
			Id:             v.Id,
			OperaterId:     v.OperaterId,
			OperaterName:   v.OperaterName,
			OperaterAction: v.OperaterAction,
			Method:         v.Method,
			Body:           body,

			Created: v.Created.Format(util.TimeFormat),
		})
	}
	return
}
