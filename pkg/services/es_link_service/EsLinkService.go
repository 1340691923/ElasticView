package es_link_service

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/model"
	"sync"
)

// EsLinkModel es连接信息表
type EsLinkService struct {
	esLinkList []model.EsLinkModel
	lock       *sync.RWMutex
	sqlx       *sqlstore.SqlStore
	logger     *logger.AppLogger
}

func (this *EsLinkService) SetEsLinkList(esLinkList []model.EsLinkModel) {
	this.lock.Lock()
	this.esLinkList = esLinkList
	this.lock.Unlock()
}

func (this *EsLinkService) EsLinkList() []model.EsLinkModel {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.esLinkList
}

func NewEsLinkService(sqlx *sqlstore.SqlStore, logger *logger.AppLogger) *EsLinkService {
	return &EsLinkService{lock: new(sync.RWMutex), sqlx: sqlx, logger: logger}
}

// 刷新eslink表数据到内存
func (this *EsLinkService) FlushEsLinkList() (err error) {
	list, err := this.GetListAction()
	if err != nil {
		return
	}
	this.SetEsLinkList(list)
	return
}

// 获取列表信息
func (this *EsLinkService) GetListAction() (esLinkList []model.EsLinkModel, err error) {
	sql, args, err := sqlstore.SqlBuilder.
		Select("*").
		From("es_link").ToSql()
	if err != nil {
		return
	}

	err = this.sqlx.Select(&esLinkList, sql, args...)
	if err != nil {
		return
	}
	return
}
