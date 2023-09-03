package gm_operater_log

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/util"
	"go.uber.org/zap"
)

type GmOperaterLogService struct {
	log  *logger.AppLogger
	sqlx *sqlstore.SqlStore
}

func NewGmOperaterLogService(log *logger.AppLogger, sqlx *sqlstore.SqlStore) *GmOperaterLogService {
	return &GmOperaterLogService{log: log, sqlx: sqlx}
}

func (this *GmOperaterLogService) List(reqData dto.GmOperaterLogList) (list []model.GmOperaterLog, count int, err error) {
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
		OperaterRoleId: reqData.RoleId,
		OperaterId:     reqData.UserId,
		OperaterAction: operater_action,
		FilterDate:     reqData.Date,
		Sqlx:           this.sqlx,
	}
	listP := &list
	err = model.SearchList(gmOperaterModel, page, limit, "*", listP, this.sqlx, this.log)
	if err != nil {
		return
	}
	count, err = model.Count(gmOperaterModel, this.sqlx, this.log)
	if err != nil {
		return
	}

	for index := range list {

		body, err := util.GzipUnCompress(list[index].Body)
		if err != nil {
			this.log.Error("err", zap.Error(err))
			continue
		}

		list[index].BodyStr = body
	}
	return
}
