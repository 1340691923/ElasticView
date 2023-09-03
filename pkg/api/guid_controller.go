package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 引导控制器
type GuidController struct {
	*BaseController
	log             *logger.AppLogger
	jwtSvr          *jwt_svr.Jwt
	esClientService *es.EsClientService
	sqlx            *sqlstore.SqlStore
}

func NewGuidController(baseController *BaseController, log *logger.AppLogger, jwtSvr *jwt_svr.Jwt, esClientService *es.EsClientService, sqlx *sqlstore.SqlStore) *GuidController {
	return &GuidController{BaseController: baseController, log: log, jwtSvr: jwtSvr, esClientService: esClientService, sqlx: sqlx}
}

// 完成新手引导
func (this *GuidController) Finish(ctx *gin.Context) {
	c, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	gmGuidModel := model.GmGuidModel{}
	err = ctx.Bind(&gmGuidModel)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	_, err = sqlstore.SqlBuilder.
		Insert(gmGuidModel.TableName()).
		SetMap(util.Map{
			"uid":       c.ID,
			"guid_name": gmGuidModel.GuidName,
			"created":   time.Now().Format(util.TimeFormat),
		}).RunWith(this.sqlx).Exec()

	if err != nil && (strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "UNIQUE")) {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
}

// 是否完成新手引导
func (this *GuidController) IsFinish(ctx *gin.Context) {
	c, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	gmGuidModel := model.GmGuidModel{}
	err = ctx.Bind(&gmGuidModel)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	sql, args, err := sqlstore.SqlBuilder.
		Select("count(*)").
		From(gmGuidModel.TableName()).
		Where(sqlstore.Eq{
			"uid":       c.ID,
			"guid_name": gmGuidModel.GuidName,
		}).ToSql()

	if err != nil {
		this.Error(ctx, err)
		return
	}
	var count int
	err = this.sqlx.Get(&count, sql, args...)
	if util.FilterMysqlNilErr(err) {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, count > 0)
}
