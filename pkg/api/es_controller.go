package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/tidwall/gjson"

	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/cat_service"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_service"
	"github.com/1340691923/ElasticView/pkg/services/index_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/cch123/elasticsql"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Es 基本操作
type EsController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	catService      *cat_service.CatService
	esService       *es_service.EsService
	jwtSvr          *jwt_svr.Jwt
	sqlx            *sqlstore.SqlStore
}

func NewEsController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, catService *cat_service.CatService, esService *es_service.EsService, jwtSvr *jwt_svr.Jwt, sqlx *sqlstore.SqlStore) *EsController {
	return &EsController{BaseController: baseController, log: log, esClientService: esClientService, catService: catService, esService: esService, jwtSvr: jwtSvr, sqlx: sqlx}
}

// Ping
func (this *EsController) PingAction(ctx *gin.Context) {
	esConnect := new(model.EsConnect)
	err := ctx.Bind(esConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if esConnect.Pwd != "" {
		pwd, decrptErr := this.esClientService.EsPwdESBDecrypt(esConnect.Pwd)
		if decrptErr == nil {
			esConnect.Pwd = pwd
		}
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esSvr := this.esService
	res, err := esSvr.Ping(ctx, esI)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if res.Version.Number == "" {
		this.Error(ctx, errors.New("ES地址OK，但是密码验证失败"))
		return
	}

	this.Success(ctx, response.OperateSuccess, res)
}

// Es 的CAT API
func (this *EsController) CatAction(ctx *gin.Context) {

	esCat := new(dto.EsCat)
	err := ctx.Bind(&esCat)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esCat.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	catSvr := this.catService

	var data *proto.Response

	switch esCat.Cat {
	case "CatHealth":
		data, err = catSvr.CatHealth(ctx, esI)
	case "CatShards":
		data, err = catSvr.CatShards(ctx, esI)
	case "CatCount":
		data, err = catSvr.CatCount(ctx, esI)
	case "CatAllocation":
		data, err = catSvr.CatAllocation(ctx, esI)
	case "CatAliases":
		data, err = catSvr.CatAliases(ctx, esI)
	case "CatIndices":
		data, err = catSvr.CatIndices(ctx, esI, []string{"store.size:desc"}, esCat.IndexBytesFormat)
	case "CatSegments":
		data, err = catSvr.IndicesSegmentsRequest(ctx, esI)

	case "CatStats":
		data, err = catSvr.ClusterStats(ctx, esI)
	case "Node":
		data, err = catSvr.CatNodes(ctx, esI)
	}

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if data.StatusErr() != nil {
		this.Error(ctx, data.StatusErr())
		return
	}

	this.Success(ctx, response.SearchSuccess, data.JsonRawMessage())

}

func (this *EsController) RunDslAction(ctx *gin.Context) {

	esRest := new(dto.EsRest)
	err := ctx.Bind(&esRest)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(esRest.EsConnect)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esRest.Method = strings.ToUpper(esRest.Method)

	if esRest.Method == "GET" {
		c, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
		if err != nil {
			this.Error(ctx, err)
			return
		}

		gmDslHistoryModel := model.GmDslHistoryModel{
			Uid:    int(c.UserID),
			Method: esRest.Method,
			Path:   esRest.Path,
			Body:   esRest.Body,
		}
		gmDslHistoryModel.Sqlx = this.sqlx
		err = gmDslHistoryModel.Insert()

		if err != nil {
			this.Error(ctx, err)
			return
		}
	}

	u, err := url.Parse(esRest.Path)

	if err != nil {
		this.Error(ctx, err)
		return
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
		this.Error(ctx, err)
		return
	}

	res, err := esI.PerformRequest(ctx, req)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if res.StatusCode() != 200 && res.StatusCode() != 201 {
		this.Output(ctx.Writer, util.Map{
			"code": res.StatusCode(),
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode())),
			"data": res.JsonRawMessage(),
		})
		return
	}

	this.Success(ctx, response.OperateSuccess, res.JsonRawMessage())
}

// SQL 转换为 DSL
func (this *EsController) SqlToDslAction(ctx *gin.Context) {

	sql := gjson.GetBytes(this.getPostBody(ctx), "sql").String()

	dsl, table, err := elasticsql.ConvertPretty(sql)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "转换成功!", util.Map{
		"dsl":       dsl,
		"tableName": table,
	})
}

// 一些索引的操作
func (this *EsController) OptimizeAction(ctx *gin.Context) {
	esOptimize := new(dto.EsOptimize)
	err := ctx.Bind(&esOptimize)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esOptimize.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	indexSvr := index_service.NewIndexService()
	if esOptimize.IndexName == "" {
		esOptimize.IndexName = "*"
	}
	switch esOptimize.Command {
	case "_refresh":
		err = indexSvr.Refresh(ctx, esI, []string{esOptimize.IndexName})
	case "_cache/clear":
		err = indexSvr.CacheClear(ctx, esI, []string{esOptimize.IndexName})
	case "_flush":
		err = indexSvr.Flush(ctx, esI, []string{esOptimize.IndexName})
	case "_forcemerge":
		err = indexSvr.IndicesForcemerge(ctx, esI, []string{esOptimize.IndexName})
	case "open":
		err = indexSvr.Open(ctx, esI, []string{esOptimize.IndexName})
	case "close":
		err = indexSvr.Close(ctx, esI, []string{esOptimize.IndexName})
	case "empty":
		err = indexSvr.Empty(ctx, esI, []string{esOptimize.IndexName})
	}

	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 将索引恢复为可写状态   由于不可抗力，ES禁止写后，默认不会自动恢复
func (this *EsController) RecoverCanWrite(ctx *gin.Context) {
	esConnectID := new(dto.EsConnectID)
	err := ctx.Bind(&esConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esConnectID.EsConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esSvr := this.esService

	err = esSvr.RecoverCanWrite(ctx, esI)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return

}
