package gm_operater_log

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

type GmOperaterLogService struct {
	log           *logger.AppLogger
	orm           *orm.Gorm
	buffer        []*model.GmOperaterLog
	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

func NewGmOperaterLogService(log *logger.AppLogger, orm *orm.Gorm) *GmOperaterLogService {
	return &GmOperaterLogService{log: log, orm: orm, batchSize: 500, flushInterval: 10, bufferMutex: new(sync.RWMutex)}
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

	err = builder.Order("id desc").Limit(limit).Offset(orm.CreatePage(page, limit)).Scan(&list).Error
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
			CostTime:       v.CostTime,
			Status:         v.Status,
			Created:        v.Created.Format(util.TimeFormat),
		})
	}
	return
}

func (this *GmOperaterLogService) Flush() (err error) {

	this.bufferMutex.Lock()
	defer this.bufferMutex.Unlock()
	if len(this.buffer) == 0 {
		return nil
	}

	err = this.orm.Create(this.buffer).Error
	if err != nil {
		this.log.Error("err", zap.Error(err))
	}

	this.buffer = make([]*model.GmOperaterLog, 0, this.batchSize)
	return nil
}

func (this *GmOperaterLogService) Save(data *model.GmOperaterLog) (err error) {
	this.bufferMutex.Lock()
	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()

	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}

	return nil
}

func (this *GmOperaterLogService) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *GmOperaterLogService) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *GmOperaterLogService) Run(ctx context.Context) error {

	ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
	run := true

	for run {
		select {
		case <-ticker.C:
			if err := this.Flush(); err != nil {
				this.log.Error("err", zap.Error(err))
			}
		case <-ctx.Done():
			run = false
		}
	}

	this.FlushAll()

	return ctx.Err()
}
