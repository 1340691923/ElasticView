package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/jwt"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/service/cat_service"
	"github.com/1340691923/ElasticView/service/es_service"
	"github.com/1340691923/ElasticView/service/index_service"
	"github.com/cch123/elasticsql"
	. "github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

// Es 基本操作
type EsController struct {
	BaseController
}

// Ping
func (this EsController) PingAction(ctx *Ctx) error {
	esConnect := new(escache.EsConnect)
	err := ctx.BodyParser(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if esConnect.Pwd != "" {
		pwd, decrptErr := escache.EsPwdESBDecrypt(esConnect.Pwd)
		if decrptErr == nil {
			esConnect.Pwd = pwd
		}
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	esSvr := es_service.NewEsService(esI)
	res, err := esSvr.Ping(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	if res.Version.Number == "" {
		return this.Error(ctx, errors.New("ES地址OK，但是密码验证失败"))
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

// Es 的CAT API
func (this EsController) CatAction(ctx *Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			//打印调用栈信息
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			logs.Logger.Sugar().Errorf("panic stack info %s", stackInfo)
			logs.Logger.Sugar().Errorf("--->HaveLoginUserSign Error:", r)
			log.Println(stackInfo)
		}
	}()
	esCat := new(escache.EsCat)
	err := ctx.BodyParser(&esCat)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esCat.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	catSvr := cat_service.NewCatService(esI)

	var data *proto.Response

	switch esCat.Cat {
	case "CatHealth":
		data, err = catSvr.CatHealth(ctx.Context())
	case "CatShards":
		data, err = catSvr.CatShards(ctx.Context())
	case "CatCount":
		data, err = catSvr.CatCount(ctx.Context())
	case "CatAllocation":
		data, err = catSvr.CatAllocation(ctx.Context())
	case "CatAliases":
		data, err = catSvr.CatAliases(ctx.Context())
	case "CatIndices":
		data, err = catSvr.CatIndices(ctx.Context(), []string{"store.size:desc"}, esCat.IndexBytesFormat)
	case "CatSegments":
		data, err = catSvr.CatSegments(ctx.Context())

	case "CatStats":
		data, err = catSvr.ClusterStats(ctx.Context())
	case "Node":
		data, err = catSvr.CatNodes(ctx.Context())
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	if data.StatusErr() != nil {
		return this.Error(ctx, data.StatusErr())
	}

	return this.Success(ctx, response.SearchSuccess, data.JsonRawMessage())

}

func (this EsController) RunDslAction(ctx *Ctx) error {

	esRest := new(escache.EsRest)
	err := ctx.BodyParser(&esRest)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := escache.GetEsClientByID(esRest.EsConnect)

	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
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

	u, err := url.Parse(esRest.Path)

	if err != nil {
		return this.Error(ctx, err)
	}

	esRest.Path = u.Path
	query := u.Query()
	query.Add("format", "json")

	var req *http.Request

	if esI.Version() > 6 {
		if len(esRest.Path) > 0 {
			if esRest.Path[0:1] != "/" {
				esRest.Path = "/" + esRest.Path
			}
		}

		if len(strings.Split(esRest.Path, "/")) == 2 || strings.Contains(esRest.Path, "/_cat") {
			esRest.Body = ""
		}

	}
	if esRest.Body == "" {
		req, err = http.NewRequest(esRest.Method, esRest.Path+"?"+query.Encode(), nil)
	} else {
		req, err = http.NewRequest(esRest.Method, esRest.Path+"?"+query.Encode(), bytes.NewReader([]byte(esRest.Body)))
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := esI.PerformRequest(ctx.Context(), req)

	if err != nil {
		return this.Error(ctx, err)
	}

	if res.StatusCode() != 200 && res.StatusCode() != 201 {
		return this.Output(ctx, util.Map{
			"code": res.StatusCode(),
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode())),
			"data": res.JsonRawMessage(),
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.JsonRawMessage())
}

// SQL 转换为 DSL
func (this EsController) SqlToDslAction(ctx *Ctx) error {
	sql := ctx.FormValue("sql")
	dsl, table, err := elasticsql.ConvertPretty(sql)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, "转换成功!", util.Map{
		"dsl":       dsl,
		"tableName": table,
	})
}

// 一些索引的操作
func (this EsController) OptimizeAction(ctx *Ctx) error {
	esOptimize := new(escache.EsOptimize)
	err := ctx.BodyParser(&esOptimize)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esOptimize.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	indexSvr := index_service.NewIndexService(esI)
	if esOptimize.IndexName == "" {
		esOptimize.IndexName = "*"
	}
	switch esOptimize.Command {
	case "_refresh":
		err = indexSvr.Refresh(ctx.Context(), []string{esOptimize.IndexName})
	case "_cache/clear":
		err = indexSvr.CacheClear(ctx.Context(), []string{esOptimize.IndexName})
	case "_flush":
		err = indexSvr.Flush(ctx.Context(), []string{esOptimize.IndexName})
	case "_forcemerge":
		err = indexSvr.IndicesForcemerge(ctx.Context(), []string{esOptimize.IndexName})
	case "open":
		err = indexSvr.Open(ctx.Context(), []string{esOptimize.IndexName})
	case "close":
		err = indexSvr.Close(ctx.Context(), []string{esOptimize.IndexName})
	case "empty":
		err = indexSvr.Empty(ctx.Context(), []string{esOptimize.IndexName})
	}

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

// 将索引恢复为可写状态   由于不可抗力，ES禁止写后，默认不会自动恢复
func (this EsController) RecoverCanWrite(ctx *Ctx) error {
	esConnectID := new(escache.EsConnectID)
	err := ctx.BodyParser(&esConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esConnectID.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	esSvr := es_service.NewEsService(esI)

	err = esSvr.RecoverCanWrite(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)

}
