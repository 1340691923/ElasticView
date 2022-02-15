package controller

import (
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	. "github.com/gofiber/fiber/v2"
)

// DSL语法查询历史记录
type DslHistoryController struct {
	BaseController
}

// 查询DSL历史记录列表
func (this DslHistoryController) ListAction(ctx *Ctx) error {
	c, err := jwt.ParseToken(ctx.Get("X-Token"))
	if err != nil {
		return this.Error(ctx, err)
	}
	gmDslHistoryModel := model.GmDslHistoryModel{}
	err = ctx.BodyParser(&gmDslHistoryModel)
	if err != nil {
		return this.Error(ctx, err)
	}
	gmDslHistoryModel.Uid = int(c.ID)

	list, err := gmDslHistoryModel.List()
	if err != nil {
		return this.Error(ctx, err)
	}
	count, err := gmDslHistoryModel.Count()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": list, "count": count})
}

// 清空DSL查询记录
func (this DslHistoryController) CleanAction(ctx *Ctx) error {
	c, err := jwt.ParseToken(ctx.Get("X-Token"))
	if err != nil {
		return this.Error(ctx, err)
	}
	gmDslHistoryModel := model.GmDslHistoryModel{}

	gmDslHistoryModel.Uid = int(c.ID)
	err = gmDslHistoryModel.Clean()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}
