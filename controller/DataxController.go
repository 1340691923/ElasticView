package controller

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/data_conversion"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"github.com/gofiber/fiber/v2"
)

type DataxController struct {
	BaseController
}

func (this DataxController) LinkInfoList(ctx *fiber.Ctx) error {

	var dataxInfoListReq request.DataxInfoListReq

	if err := ctx.BodyParser(&dataxInfoListReq); err != nil {
		return this.Error(ctx, err)
	}

	if dataxInfoListReq.Page == 0 {
		dataxInfoListReq.Page = 1
	}

	if dataxInfoListReq.Limit == 0 {
		dataxInfoListReq.Limit = 10
	}

	sqlBuilder := db.SqlBuilder.Select("*").From("datax_link_info")

	countbuilder := db.SqlBuilder.Select("count(*)").From("datax_link_info")

	if dataxInfoListReq.Typ != "" {
		sqlBuilder = sqlBuilder.Where(db.Eq{"typ": dataxInfoListReq.Typ})
		countbuilder = countbuilder.Where(db.Eq{"typ": dataxInfoListReq.Typ})
	}

	if dataxInfoListReq.Remark != "" {
		sqlBuilder = sqlBuilder.Where(db.Like{"remark": db.CreateLike(dataxInfoListReq.Remark)})
		countbuilder = countbuilder.Where(db.Like{"remark": db.CreateLike(dataxInfoListReq.Remark)})
	}

	var count int

	err := countbuilder.RunWith(db.Sqlx).QueryRow().Scan(&count)
	if err != nil {
		return this.Error(ctx, err)
	}

	sql, args, err := sqlBuilder.
		OrderBy("id desc").
		Limit(uint64(dataxInfoListReq.Limit)).
		Offset(db.CreatePage(dataxInfoListReq.Page, dataxInfoListReq.Limit)).
		RunWith(db.Sqlx).
		ToSql()

	if err != nil {
		return this.Error(ctx, err)
	}

	list := []model.DataxLinkInfoModel{}

	err = db.Sqlx.Select(&list, sql, args...)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, util.Map{
		"data":  list,
		"count": count,
	})
}

func (this DataxController) InsertLink(ctx *fiber.Ctx) error {
	var reqData request.DataxInfoInsertReq
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	if _, err := db.SqlBuilder.Insert("datax_link_info").SetMap(db.Model2Map(db.Model2MapParmas{
		M:                reqData,
		NeedZeroByInt:    true,
		NeedZeroByString: true,
		CreateTimeCol:    "created",
		UpdateTimeCol:    "updated",
	})).RunWith(db.Sqlx).Exec(); err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this DataxController) DelLinkById(ctx *fiber.Ctx) error {
	var reqData request.DataxInfoDelReq
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	if _, err := db.SqlBuilder.Delete("datax_link_info").Where(db.Eq{"id": reqData.ID}).RunWith(db.Sqlx).Exec(); err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.DeleteSuccess, nil)
}

func (this DataxController) TestLink(ctx *fiber.Ctx) error {
	var reqData request.DataxInfoTestLinkReq
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}
	dataSource, err := data_conversion.NewDataSource(reqData)
	if err != nil {
		return this.Error(ctx, err)
	}
	err = dataSource.Ping()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.LinkSuccess, nil)
}
