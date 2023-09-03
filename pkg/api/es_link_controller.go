package api

import (
	"encoding/json"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/1340691923/ElasticView/pkg/response"
)

// Es 连接管理控制器
type EsLinkController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	sqlx            *sqlstore.SqlStore
	esCache         *es.EsCache
	esLinkService   *es_link_service.EsLinkService
}

func NewEsLinkController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, sqlx *sqlstore.SqlStore, esCache *es.EsCache, esLinkService *es_link_service.EsLinkService) *EsLinkController {
	return &EsLinkController{BaseController: baseController, log: log, esClientService: esClientService, sqlx: sqlx, esCache: esCache, esLinkService: esLinkService}
}

// 获取Es连接列表
func (this *EsLinkController) ListAction(ctx *gin.Context) {

	list, err := this.esLinkService.GetListAction()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
}

func (this *EsLinkController) OptAction(ctx *gin.Context) {

	type Opt struct {
		ID     int64  `json:"id"`
		Remark string `json:"remark"`
	}

	var optList []Opt

	for _, esLink := range this.esLinkService.EsLinkList() {
		optList = append(optList, Opt{ID: esLink.ID, Remark: esLink.Remark})
	}

	this.Success(ctx, response.SearchSuccess, optList)

}

// 新增Es连接
func (this *EsLinkController) InsertAction(ctx *gin.Context) {

	var esLinkModel model.EsLinkModel
	err := ctx.Bind(&esLinkModel)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esB, err := json.Marshal(esLinkModel)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	insertMap := map[string]interface{}{}
	err = json.Unmarshal(esB, &insertMap)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	delete(insertMap, "created")
	delete(insertMap, "updated")
	delete(insertMap, "id")
	insertMap["created"] = time.Now().Format(util.TimeFormat)
	insertMap["updated"] = time.Now().Format(util.TimeFormat)
	if insertMap["pwd"] != "" {
		insertMap["pwd"], err = this.esClientService.EsPwdESBEncrypt(insertMap["pwd"].(string))
		if err != nil {
			this.Error(ctx, err)
			return
		}
	}

	_, err = sqlstore.SqlBuilder.
		Insert("es_link").
		SetMap(insertMap).
		RunWith(this.sqlx).
		Exec()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esLinkService.FlushEsLinkList()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 修改Es连接信息
func (this *EsLinkController) UpdateAction(ctx *gin.Context) {
	var esLinkModel model.EsLinkModel
	err := ctx.Bind(&esLinkModel)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esB, err := json.Marshal(esLinkModel)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	insertMap := map[string]interface{}{}
	err = json.Unmarshal(esB, &insertMap)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	delete(insertMap, "id")
	delete(insertMap, "created")
	delete(insertMap, "updated")
	insertMap["updated"] = time.Now().Format(util.TimeFormat)
	if insertMap["pwd"] != "" {
		insertMap["pwd"], err = this.esClientService.EsPwdESBEncrypt(insertMap["pwd"].(string))
		if err != nil {
			this.Error(ctx, err)
			return
		}
	}

	_, err = sqlstore.SqlBuilder.
		Update("es_link").
		SetMap(insertMap).
		Where(sqlstore.Eq{"id": esLinkModel.ID}).
		RunWith(this.sqlx).
		Exec()
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esCache := this.esCache
	esCache.Rem(int(esLinkModel.ID))

	err = this.esLinkService.FlushEsLinkList()
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 删除es连接
func (this *EsLinkController) DeleteAction(ctx *gin.Context) {

	var req struct {
		Id int `json:"id"`
	}

	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	_, err = sqlstore.SqlBuilder.
		Delete("es_link").
		Where(sqlstore.Eq{"id": req.Id}).RunWith(this.sqlx).Exec()

	if err != nil {
		this.Error(ctx, err)
		return
	}

	esCache := this.esCache
	esCache.Rem(req.Id)

	err = this.esLinkService.FlushEsLinkList()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.DeleteSuccess, nil)
}
