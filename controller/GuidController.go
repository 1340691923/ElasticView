package controller

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	. "github.com/gofiber/fiber/v2"
)

// 引导控制器
type GuidController struct {
	BaseController
}

// 完成新手引导
func (this GuidController) Finish(ctx *Ctx) error {
	c, err := jwt.ParseToken(ctx.Get("X-Token"))
	if err != nil {
		return this.Error(ctx, err)
	}

	gmGuidModel := model.GmGuidModel{}
	err = ctx.BodyParser(&gmGuidModel)
	if err != nil {
		return this.Error(ctx, err)
	}
	_, err = db.SqlBuilder.
		Insert(gmGuidModel.TableName()).
		SetMap(map[string]interface{}{
			"uid":       c.ID,
			"guid_name": gmGuidModel.GuidName,
		}).RunWith(db.Sqlx).Exec()

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//是否完成新手引导
func (this GuidController) IsFinish(ctx *Ctx) error {
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
