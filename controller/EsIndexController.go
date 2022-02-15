package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	. "github.com/gofiber/fiber/v2"

	"github.com/olivere/elastic"
)

// Es 索引控制器
type EsIndexController struct {
	BaseController
}

//创建索引
func (this EsIndexController) CreateAction(ctx *Ctx) error {
	esIndexInfo := es.EsIndexInfo{}
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}
	var res interface{}
	if esIndexInfo.Types == "update" {
		res, err = esClinet.IndexPutSettings(esIndexInfo.IndexName, esIndexInfo.Settings)
		if err != nil {
			return this.Error(ctx, err)
		}

	} else {
		res, err = esClinet.(*es.EsClientV6).Client.CreateIndex(esIndexInfo.IndexName).BodyJson(map[string]interface{}{
			"settings": esIndexInfo.Settings,
		}).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

// 删除索引
func (this EsIndexController) DeleteAction(ctx *Ctx) error {
	esIndexInfo := es.EsIndexInfo{}
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))

	}
	_, err = esClinet.DeleteIndex(strings.Split(esIndexInfo.IndexName, ","))
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//获取索引配置信息
func (this EsIndexController) GetSettingsAction(ctx *Ctx) error {
	esIndexInfo := es.EsIndexInfo{}
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esIndexInfo.IndexName == "" {
		this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}

	res, err := esClinet.(*es.EsClientV6).Client.IndexGetSettings(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res[esIndexInfo.IndexName].Settings)
}

//获取所有的索引配置信息
func (this EsIndexController) GetSettingsInfoAction(ctx *Ctx) error {
	esIndexInfo := es.EsIndexInfo{}
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}

	res, err := esClinet.(*es.EsClientV6).Client.IndexGetSettings(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

// 获取别名
func (this EsIndexController) GetAliasAction(ctx *Ctx) error {
	esAliasInfo := es.EsAliasInfo{}
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esAliasInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}

	aliasRes, err := esClinet.(*es.EsClientV6).Client.Aliases().Index(esAliasInfo.IndexName).Do(ctx.Context())

	return this.Success(ctx, response.OperateSuccess, aliasRes.Indices[esAliasInfo.IndexName].Aliases)
}

// 操作别名
func (this EsIndexController) OperateAliasAction(ctx *Ctx) error {
	esAliasInfo := es.EsAliasInfo{}
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	const Add = 1
	const Delete = 2
	const MoveToAnotherIndex = 3
	const PatchAdd = 4
	var res interface{}
	switch esAliasInfo.Types {
	case Add:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		}
		res, err = esClinet.(*es.EsClientV6).Client.Alias().Add(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(ctx.Context())
	case Delete:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		}
		res, err = esClinet.(*es.EsClientV6).Client.Alias().Remove(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(ctx.Context())
	case MoveToAnotherIndex:
		res, err = esClinet.(*es.EsClientV6).Client.Alias().Action(elastic.NewAliasAddAction(esAliasInfo.AliasName).Index(esAliasInfo.NewIndexList...)).Do(ctx.Context())
	case PatchAdd:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		}
		wg := sync.WaitGroup{}
		NewAliasNameListLen := len(esAliasInfo.NewAliasNameList)
		if len(esAliasInfo.NewAliasNameList) > 10 {
			err = errors.New("别名列表数量不能大于10")
			break
		} else {
			wg.Add(NewAliasNameListLen)
			for _, aliasName := range esAliasInfo.NewAliasNameList {
				go func(aliasName string) {
					defer wg.Done()
					res, err = esClinet.(*es.EsClientV6).Client.Alias().
						Add(esAliasInfo.IndexName, aliasName).
						Do(context.TODO())
				}(aliasName)
			}
			wg.Wait()
		}
	default:
		err = es.ReqParmasValid
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

// 重建索引
func (this EsIndexController) ReindexAction(ctx *Ctx) error {
	esReIndexInfo := es.EsReIndexInfo{}
	err := ctx.BodyParser(&esReIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esReIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	reindex := esClinet.(*es.EsClientV6).Client.Reindex()
	urlValues := esReIndexInfo.UrlValues
	if urlValues.WaitForActiveShards != "" {
		reindex = reindex.WaitForActiveShards(urlValues.WaitForActiveShards)
	}
	if urlValues.Slices != 0 {
		reindex = reindex.Slices(urlValues.Slices)
	}
	if urlValues.Refresh != "" {
		reindex = reindex.Refresh(urlValues.Refresh)
	}
	if urlValues.Timeout != "" {
		reindex = reindex.Timeout(urlValues.Refresh)
	}
	if urlValues.RequestsPerSecond != 0 {
		reindex = reindex.RequestsPerSecond(urlValues.RequestsPerSecond)
	}
	if urlValues.WaitForCompletion != nil {
		reindex = reindex.WaitForCompletion(*urlValues.WaitForCompletion)
	}

	res, err := reindex.Body(esReIndexInfo.Body).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

// 得到所有的索引名
func (this EsIndexController) IndexNamesAction(ctx *Ctx) error {
	esConnect := es.EsConnectID{}
	err := ctx.BodyParser(&esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esConnect.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	catIndicesResponse, err := esClinet.(*es.EsClientV6).Client.CatIndices().Human(true).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	indexNames := []string{}

	for _, catIndices := range catIndicesResponse {
		indexNames = append(indexNames, catIndices.Index)
	}

	return this.Success(ctx, response.SearchSuccess, indexNames)
}

// 获取索引的Stats
func (this EsIndexController) StatsAction(ctx *Ctx) error {
	esIndexInfo := es.EsIndexInfo{}
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))

	}

	res, err := esClinet.(*es.EsClientV6).Client.IndexStats(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsIndexController) CatStatusAction(ctx *Ctx) error {
	esIndexInfo := es.EsIndexInfo{}
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(ctx.Context(), elastic.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_cat/indices/%s?h=status", esIndexInfo.IndexName),
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res.Body)
}
