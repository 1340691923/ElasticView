package api

import (
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/jwt"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	. "github.com/gofiber/fiber/v2"
)

type SearchController struct {
	BaseController
}

type GetIndexConfigsReq struct {
	EsConnectID int  `json:"es_connect"`
	Limit       int  `json:"limit"`
	Page        int  `json:"page"`
	All         bool `json:"all"`
}

type SetIndexConfigReq struct {
	EsConnectID int    `json:"es_connect"`
	IndexName   string `json:"indexName"`
	Remark      string `json:"remark"`
}

func (this SearchController) SetIndexConfig(ctx *Ctx) error {

	req := SetIndexConfigReq{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return this.Error(ctx, err)
	}
	m := model.SearchConfig{}

	//REPLACE INTO `visits` (ip, VALUE) VALUES ($ip, 0);

	_, err = db.Sqlx.Exec("REPLACE into "+m.TableName()+" (index_name, remark, es_connect) values   (?,?,?)",
		req.IndexName, req.Remark, req.EsConnectID)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this SearchController) GetIndexConfigs(ctx *Ctx) error {

	req := GetIndexConfigsReq{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return this.Error(ctx, err)
	}
	m := model.SearchConfig{}
	m.EsConnect = req.EsConnectID

	if req.All {
		list, err := m.All()
		if err != nil {
			return this.Error(ctx, err)
		}
		return this.Success(ctx, response.SearchSuccess, util.Map{"list": list})
	}

	m.Limit = req.Limit
	m.Page = req.Page

	list, err := m.List()
	if err != nil {
		return this.Error(ctx, err)
	}
	count, err := m.Count()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, util.Map{"list": list, "count": count})
}

func (this SearchController) SearchLog(ctx *Ctx) error {
	c, err := jwt.ParseToken(ctx.Get("X-Token"))
	if err != nil {
		return this.Error(ctx, err)
	}

	gmGuidModel := model.GmGuidModel{}
	err = ctx.BodyParser(&gmGuidModel)
	if err != nil {
		return this.Error(ctx, err)
	}
	sql, args, err := db.SqlBuilder.
		Select("count(*)").
		From(gmGuidModel.TableName()).
		Where(db.Eq{
			"uid":       c.ID,
			"guid_name": gmGuidModel.GuidName,
		}).ToSql()

	if err != nil {
		return this.Error(ctx, err)
	}
	var count int
	err = db.Sqlx.Get(&count, sql, args...)
	if util.FilterMysqlNilErr(err) {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, count > 0)
}
