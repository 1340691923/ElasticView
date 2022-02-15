package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/es_optimize"
	. "github.com/gofiber/fiber/v2"

	"github.com/cch123/elasticsql"
	"github.com/olivere/elastic"
)

//Es 基本操作
type EsController struct {
	BaseController
}

// Ping
func (this EsController) PingAction(ctx *Ctx) error {
	esConnect := es.EsConnect{}
	err := ctx.BodyParser(&esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClient(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	data, _, err := esClinet.Ping()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, data)
}

// Es 的CAT API
func (this EsController) CatAction(ctx *Ctx) error {

	esCat := es.EsCat{}
	err := ctx.BodyParser(&esCat)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esCat.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	var data interface{}

	switch esCat.Cat {
	case "CatHealth":
		data, err = esClinet.(*es.EsClientV6).Client.CatHealth().Human(true).Do(ctx.Context())
	case "CatShards":
		data, err = esClinet.(*es.EsClientV6).Client.CatShards().Human(true).Do(ctx.Context())
	case "CatCount":
		data, err = esClinet.(*es.EsClientV6).Client.CatCount().Human(true).Do(ctx.Context())
	case "CatAllocation":
		data, err = esClinet.(*es.EsClientV6).Client.CatAllocation().Human(true).Do(ctx.Context())
	case "CatAliases":
		data, err = esClinet.(*es.EsClientV6).Client.CatAliases().Human(true).Do(ctx.Context())
	case "CatIndices":
		if esCat.IndexBytesFormat != "" {
			data, err = esClinet.(*es.EsClientV6).Client.CatIndices().Human(true).Bytes(esCat.IndexBytesFormat).Do(ctx.Context())
		} else {
			data, err = esClinet.(*es.EsClientV6).Client.CatIndices().Human(true).Do(ctx.Context())
		}
	case "CatSegments":
		data, err = esClinet.(*es.EsClientV6).Client.IndexSegments().Human(true).Do(ctx.Context())
	case "CatStats":
		data, err = esClinet.(*es.EsClientV6).Client.ClusterStats().Human(true).Do(ctx.Context())
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, data)
}

func (this EsController) RunDslAction(ctx *Ctx) error {

	esRest := es.EsRest{}
	err := ctx.BodyParser(&esRest)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esRest.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esRest.Method = strings.ToUpper(esRest.Method)
	if esRest.Method == "GET" {
		c, err := jwt.ParseToken(ctx.Get("X-Token"))
		if err != nil {
			return this.Error(ctx, err)
		}

		gmDslHistoryModel := model.GmDslHistoryModel{
			Uid:    int(c.ID),
			Method: esRest.Method,
			Path:   esRest.Path,
			Body:   esRest.Body,
		}

		err = gmDslHistoryModel.Insert()

		if err != nil {
			return this.Error(ctx, err)
		}
	}

	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(context.TODO(), elastic.PerformRequestOptions{
		Method: esRest.Method,
		Path:   esRest.Path,
		Body:   esRest.Body,
	})

	if err != nil {
		return this.Error(ctx, err)
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return this.Output(ctx, map[string]interface{}{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

// SQL 转换为 DSL
func (this EsController) SqlToDslAction(ctx *Ctx) error {
	sql := ctx.FormValue("sql")
	dsl, table, err := elasticsql.ConvertPretty(sql)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, "转换成功!", map[string]interface{}{
		"dsl":       dsl,
		"tableName": table,
	})
}

// 一些索引的操作
func (this EsController) OptimizeAction(ctx *Ctx) error {
	esOptimize := es.EsOptimize{}
	err := ctx.BodyParser(&esOptimize)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esOptimize.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	optimize := es_optimize.OptimizeFactory(esOptimize.Command)

	if optimize == nil {
		return this.Error(ctx, errors.New("不支持该指令"))

	}
	if esOptimize.IndexName != "" {
		optimize.SetIndexName(esOptimize.IndexName)
	}
	err = optimize.Do(esClinet.(*es.EsClientV6).Client)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

// 将索引恢复为可写状态   由于不可抗力，ES禁止写后，默认不会自动恢复
func (this EsController) RecoverCanWrite(ctx *Ctx) error {
	esConnect := es.EsConnectID{}
	err := ctx.BodyParser(&esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esConnect.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(ctx.Context(), elastic.PerformRequestOptions{
		Method: "PUT",
		Path:   "/_settings",
		Body: map[string]interface{}{
			"index": map[string]interface{}{
				"blocks": map[string]interface{}{
					"read_only_allow_delete": "false",
				},
			},
		},
	})

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return this.Output(ctx, map[string]interface{}{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}
