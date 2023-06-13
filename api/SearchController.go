package api

import (
	"encoding/json"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	es2 "github.com/1340691923/ElasticView/service/es"
	. "github.com/gofiber/fiber/v2"
	"strings"
)

type SearchController struct {
	BaseController
}

func (this SearchController) SetIndexConfig(ctx *Ctx) error {

	req := escache.SetIndexConfigReq{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return this.Error(ctx, err)
	}
	m := model.SearchConfig{}

	_, err = db.Sqlx.Exec("REPLACE into "+m.TableName()+" (index_name, remark, es_connect,input_cols,output_cols)"+
		" values   (?,?,?,?,?)",
		req.IndexName, req.Remark, req.EsConnectID, strings.Join(req.InputCols, ","), strings.Join(req.OutputCols, ","))

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this SearchController) GetIndexConfigs(ctx *Ctx) error {

	req := escache.GetIndexConfigsReq{}
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

		return this.Success(ctx, response.SearchSuccess, util.Map{"list": response.ToSearchConfig(list)})
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
	return this.Success(ctx, response.SearchSuccess, util.Map{"list": response.ToSearchConfig(list), "count": count})
}

func (this SearchController) SearchLog(ctx *Ctx) error {

	esIndexInfo := new(escache.SearchlogReq)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SearchLog(ctx, esIndexInfo)
}

func (this SearchController) SetMappingAlias(ctx *Ctx) error {

	req := escache.SetMappingAliasReq{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return this.Error(ctx, err)
	}
	mpCfg, _ := json.Marshal(req.MappingCfg)

	_, err = db.Sqlx.Exec("REPLACE into mapping_alias_config (index_name, es_connect,col_alias_map)"+
		" values (?,?,?)",
		req.IndexName, req.EsConnect, string(mpCfg))

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this SearchController) GetMappingAlias(ctx *Ctx) error {

	req := escache.GetMappingAliasReq{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return this.Error(ctx, err)
	}
	res := map[string]string{}

	mappingCfgStr := ""

	err = db.Sqlx.QueryRow("select col_alias_map from mapping_alias_config where es_connect = ? and index_name = ?;",
		req.EsConnectID, req.IndexName).Scan(&mappingCfgStr)

	if util.FilterMysqlNilErr(err) {
		return this.Error(ctx, err)
	}
	json.Unmarshal([]byte(mappingCfgStr), &res)

	return this.Success(ctx, response.SearchSuccess, util.Map{"res": res})
}
