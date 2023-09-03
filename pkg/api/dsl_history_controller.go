package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
)

// DSL语法查询历史记录
type DslHistoryController struct {
	jwtSvr *jwt_svr.Jwt
	sqlx   *sqlstore.SqlStore
	log    *logger.AppLogger
	*BaseController
}

func NewDslHistoryController(jwtSvr *jwt_svr.Jwt, sqlx *sqlstore.SqlStore, log *logger.AppLogger, baseController *BaseController) *DslHistoryController {
	return &DslHistoryController{jwtSvr: jwtSvr, sqlx: sqlx, log: log, BaseController: baseController}
}

// 查询DSL历史记录列表
func (this *DslHistoryController) ListAction(ctx *gin.Context) {
	c, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}
	gmDslHistoryModel := model.GmDslHistoryModel{}
	gmDslHistoryModel.Sqlx = this.sqlx
	err = ctx.Bind(&gmDslHistoryModel)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	gmDslHistoryModel.Uid = int(c.UserID)

	list, err := gmDslHistoryModel.List()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	count, err := gmDslHistoryModel.Count()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, util.Map{"list": list, "count": count})
	return
}

// 清空DSL查询记录
func (this *DslHistoryController) CleanAction(ctx *gin.Context) {
	c, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}
	gmDslHistoryModel := model.GmDslHistoryModel{}
	gmDslHistoryModel.Sqlx = this.sqlx
	gmDslHistoryModel.Uid = int(c.UserID)
	err = gmDslHistoryModel.Clean()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return
}
