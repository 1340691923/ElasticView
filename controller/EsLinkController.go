package controller

import (
	"encoding/json"

	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	. "github.com/gofiber/fiber/v2"

	jsoniter "github.com/json-iterator/go"
)

// Es 连接管理控制器
type EsLinkController struct {
	BaseController
}

// 获取Es连接列表
func (this EsLinkController) ListAction(ctx *Ctx) error {

	esLinkModel := model.EsLinkModel{}

	list, err := esLinkModel.GetListAction()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, list)
}

func (this EsLinkController) OptAction(ctx *Ctx) error {

	type Opt struct {
		ID     int64  `json:"id"`
		Remark string `json:"remark"`
	}

	var optList []Opt

	for _, esLink := range model.EsLinkList {
		optList = append(optList, Opt{ID: esLink.ID, Remark: esLink.Remark})
	}

	return this.Success(ctx, response.SearchSuccess, optList)

}

// 新增Es连接
func (this EsLinkController) InsertAction(ctx *Ctx) error {

	var esLinkModel model.EsLinkModel
	err := ctx.BodyParser(&esLinkModel)
	if err != nil {
		return this.Error(ctx, err)
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	esB, err := json.Marshal(esLinkModel)

	if err != nil {
		return this.Error(ctx, err)
	}
	insertMap := map[string]interface{}{}
	err = json.Unmarshal(esB, &insertMap)
	if err != nil {
		return this.Error(ctx, err)
	}

	delete(insertMap, "created")
	delete(insertMap, "updated")

	_, err = db.SqlBuilder.
		Insert("es_link").
		SetMap(insertMap).
		RunWith(db.Sqlx).
		Exec()
	if err != nil {
		return this.Error(ctx, err)
	}
	err = esLinkModel.FlushEsLinkList()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

// 修改Es连接信息
func (this EsLinkController) UpdateAction(ctx *Ctx) error {
	var esLinkModel model.EsLinkModel
	err := ctx.BodyParser(&esLinkModel)
	if err != nil {
		return this.Error(ctx, err)
	}

	esB, err := json.Marshal(esLinkModel)

	if err != nil {
		return this.Error(ctx, err)
	}
	insertMap := map[string]interface{}{}
	err = json.Unmarshal(esB, &insertMap)
	if err != nil {
		return this.Error(ctx, err)
	}

	delete(insertMap, "id")
	delete(insertMap, "created")
	delete(insertMap, "updated")

	_, err = db.SqlBuilder.
		Update("es_link").
		SetMap(insertMap).
		Where(db.Eq{"id": esLinkModel.ID}).
		RunWith(db.Sqlx).
		Exec()
	if err != nil {
		return this.Error(ctx, err)
	}

	esCache := es.NewEsCache()
	esCache.Rem(int(esLinkModel.ID))

	err = esLinkModel.FlushEsLinkList()
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 删除es连接
func (this EsLinkController) DeleteAction(ctx *Ctx) error {

	var req struct {
		Id int `json:"id"`
	}

	err := ctx.BodyParser(&req)
	if err != nil {
		return this.Error(ctx, err)
	}

	_, err = db.SqlBuilder.
		Delete("es_link").
		Where(db.Eq{"id": req.Id}).RunWith(db.Sqlx).Exec()

	if err != nil {
		return this.Error(ctx, err)
	}

	esCache := es.NewEsCache()
	esCache.Rem(req.Id)
	esLinkModel := model.EsLinkModel{}
	err = esLinkModel.FlushEsLinkList()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.DeleteSuccess, nil)
}
