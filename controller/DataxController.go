package controller

import (
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/data_conversion"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"time"
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

func (this DataxController) LinkSelectOpt(ctx *fiber.Ctx) error {
	type D struct {
		Id     int    `json:"id" db:"id"`
		Remark string `json:"remark" db:"remark"`
		Typ    string `json:"typ" db:"typ"`
	}
	list := []D{}
	err := db.Sqlx.Select(&list, "select id,remark,typ from datax_link_info")
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, list)
}

func (this DataxController) Tables(ctx *fiber.Ctx) error {
	var reqData struct {
		Id int `json:"id" db:"id"`
	}
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	var obj model.DataxLinkInfoModel
	err := db.Sqlx.Get(&obj, "select * from datax_link_info where id = ?", reqData.Id)
	if err != nil {
		return this.Error(ctx, err)
	}
	dataSource, err := data_conversion.NewDataSource(request.DataxInfoTestLinkReq{
		IP:       obj.Ip,
		Port:     obj.Port,
		DbName:   obj.DbName,
		Username: obj.Username,
		Pwd:      obj.Pwd,
		Remark:   obj.Remark,
		Typ:      obj.Typ,
	})
	if err != nil {
		return this.Error(ctx, err)
	}
	tables, err := dataSource.GetTables()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, tables)
}

func (this DataxController) GetTableColumns(ctx *fiber.Ctx) error {
	var reqData struct {
		Id        int    `json:"id" db:"id"`
		TableName string `json:"table_name"`
	}
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}

	var obj model.DataxLinkInfoModel
	err := db.Sqlx.Get(&obj, "select * from datax_link_info where id = ?", reqData.Id)
	if err != nil {
		return this.Error(ctx, err)
	}
	dataSource, err := data_conversion.NewDataSource(request.DataxInfoTestLinkReq{
		IP:       obj.Ip,
		Port:     obj.Port,
		DbName:   obj.DbName,
		Username: obj.Username,
		Pwd:      obj.Pwd,
		Remark:   obj.Remark,
		Typ:      obj.Typ,
	})
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := dataSource.GetTableColumns(reqData.TableName)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}

func (this DataxController) TransferLogList(ctx *fiber.Ctx) error {
	sql, args, err := db.SqlBuilder.Select("*").From("datax_transfer_list").ToSql()
	if err != nil {
		return this.Error(ctx, err)
	}
	var list []model.DataxListModel
	err = db.Sqlx.Select(&list, sql, args...)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, list)
}

func (this DataxController) Transfer(ctx *fiber.Ctx) error {
	var reqData request.TransferReq
	if err := ctx.BodyParser(&reqData); err != nil {
		return this.Error(ctx, err)
	}
	if reqData.Remark == "" {
		return this.Error(ctx, errors.New("备注不能为空"))
	}
	if reqData.FlushInterval == 0 {
		return this.Error(ctx, errors.New("入库轮循间隔时间不能为空"))
	}
	if reqData.BufferSize == 0 {
		return this.Error(ctx, errors.New("入库批次数量不能为空"))
	}
	if reqData.BufferSize == 0 {
		return this.Error(ctx, errors.New("入库批次数量不能为空"))
	}
	if reqData.IndexName == "" {
		return this.Error(ctx, errors.New("索引名不能为空"))
	}
	if reqData.SelectTable == "" {
		return this.Error(ctx, errors.New("表名不能为空"))
	}
	if len(reqData.Cols.TableCols) == 0 {
		return this.Error(ctx, errors.New("表字段不能为空"))
	}
	for _, col := range reqData.Cols.EsCols {
		if col.Col == "" {
			return this.Error(ctx, errors.New(fmt.Sprintf("表字段（%s）所映射的ES字段不能为空", col.TbCol)))
		}
	}
	selectType, err := reqData.ParseSelectType()
	if err != nil {
		return this.Error(ctx, err)
	}

	var obj model.DataxLinkInfoModel
	err = db.Sqlx.Get(&obj, "select * from datax_link_info where id = ?", selectType.ID)
	if err != nil {
		return this.Error(ctx, err)
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	formData, _ := json.MarshalToString(reqData)
	now := time.Now().Format(util.TimeFormat)

	_, err = db.SqlBuilder.
		Insert("datax_transfer_list").
		Columns("form_data", "remark", "table_name", "index_name", "error_msg", "status", "updated", "created").
		Values(formData, reqData.Remark, reqData.SelectTable, reqData.IndexName, "无", "正在运行中...", now, now).RunWith(db.Sqlx).Exec()
	if err != nil {
		return this.Error(ctx, err)
	}
	dataSource, err := data_conversion.NewDataSource(request.DataxInfoTestLinkReq{
		IP:       obj.Ip,
		Port:     obj.Port,
		DbName:   obj.DbName,
		Username: obj.Username,
		Pwd:      obj.Pwd,
		Remark:   obj.Remark,
		Typ:      obj.Typ,
	})

	err = dataSource.Ping()
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
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
