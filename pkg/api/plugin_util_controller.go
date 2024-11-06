package api

import (
	"bytes"
	"database/sql"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_service"
	"github.com/1340691923/ElasticView/pkg/services/plugin_service"
	dto2 "github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	vo2 "github.com/1340691923/eve-plugin-sdk-go/ev_api/vo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type PluginUtilController struct {
	*BaseController
	pluginServer    *plugin_service.PluginService
	esClientService *es.EsClientService
	esService       *es_service.EsService
	log             *logger.AppLogger
	cfg             *config.Config
}

func NewPluginUtilController(baseController *BaseController, pluginServer *plugin_service.PluginService, esClientService *es.EsClientService, esService *es_service.EsService, log *logger.AppLogger, cfg *config.Config) *PluginUtilController {
	return &PluginUtilController{BaseController: baseController, pluginServer: pluginServer, esClientService: esClientService, esService: esService, log: log, cfg: cfg}
}

// 进行增删改等操作
func (this *PluginUtilController) ExecSql(ctx *gin.Context) {
	var reqData dto.ExecReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	rowsAffected, err := this.pluginServer.ExecSql(ctx, reqData.PluginId, reqData.Sql, reqData.Args)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res := vo.ExecSqlRes{
		RowsAffected: rowsAffected,
	}

	this.Success(ctx, response.OperateSuccess, res)
}

// 进行查询操作
func (this *PluginUtilController) SelectSql(ctx *gin.Context) {
	var reqData dto.SelectReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.pluginServer.SelectSql(ctx, reqData.PluginId, reqData.Sql, reqData.Args)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	if len(res) == 0 {
		this.Error(ctx, sql.ErrNoRows)
		return
	}
	this.Success(ctx, response.SearchSuccess, vo.SelectSqlRes{Result: res})
}

// 进行查询操作
func (this *PluginUtilController) FirstSql(ctx *gin.Context) {
	var reqData dto.SelectReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	storeRes, err := this.pluginServer.SelectSql(ctx, reqData.PluginId, reqData.Sql, reqData.Args)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	var res map[string]interface{}
	if len(storeRes) > 0 {
		res = storeRes[0]
	} else {
		this.Error(ctx, sql.ErrNoRows)
		return
	}

	this.Success(ctx, response.SearchSuccess, vo.FirstSqlRes{Result: res})
}

// 进行原生es操作
func (this *PluginUtilController) EsRunDsl(ctx *gin.Context) {
	req := new(dto2.PluginRunDsl)
	err := ctx.BindJSON(&req)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(ctx, req.EsConnectData.EsConnect, req.EsConnectData.UserID)

	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, req.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := this.esService.RunDsl(ctx, esI, req.EsConnectData.UserID, req.HttpMethod, req.Path, req.Dsl)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) Ping(ctx *gin.Context) {
	var reqData dto2.PingReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.Ping(ctx)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

// 加载需调试插件
func (this *PluginUtilController) LoadDebugPlugin(ctx *gin.Context) {

	var reqData dto.LoadDebugPlugin

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginServer.LoadDebugPlugin(ctx, reqData.ID, reqData.Addr, reqData.Pid)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
}

// 停止需调试插件
func (this *PluginUtilController) StopDebugPlugin(ctx *gin.Context) {

	var reqData dto.StopDebugPlugin

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.pluginServer.StopDebugPlugin(ctx, reqData.ID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
}

func (this *PluginUtilController) EsVersion(ctx *gin.Context) {

	var reqData dto2.EsConnectData

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnect, reqData.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnect))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	version, err := esI.EsVersion()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, version)
}

func (this *PluginUtilController) EsCatNodes(ctx *gin.Context) {
	var reqData dto2.CatNodesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCatNodes(ctx, reqData.CatNodeReqData.H)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsClusterStats(ctx *gin.Context) {
	var reqData dto2.ClusterStatsReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsClusterStats(ctx, reqData.ClusterStatsReqData.Human)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsIndicesSegmentsRequest(ctx *gin.Context) {
	var reqData dto2.IndicesSegmentsRequest

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsIndicesSegmentsRequest(ctx, reqData.IndicesSegmentsRequestData.Human)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsRefresh(ctx *gin.Context) {
	var reqData dto2.RefreshReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsRefresh(ctx, reqData.RefreshReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsOpen(ctx *gin.Context) {
	var reqData dto2.OpenReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsOpen(ctx, reqData.OpenReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsFlush(ctx *gin.Context) {
	var reqData dto2.FlushReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsFlush(ctx, reqData.FlushReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsIndicesClearCache(ctx *gin.Context) {
	var reqData dto2.IndicesClearCacheReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsIndicesClearCache(ctx, reqData.IndicesClearCacheReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsIndicesClose(ctx *gin.Context) {
	var reqData dto2.IndicesCloseReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsIndicesClose(ctx, reqData.IndicesCloseReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsIndicesForcemerge(ctx *gin.Context) {
	var reqData dto2.IndicesForcemergeReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsIndicesForcemerge(ctx, reqData.IndicesForcemergeReqData.IndexNames, reqData.IndicesForcemergeReqData.MaxNumSegments)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsDeleteByQuery(ctx *gin.Context) {
	var reqData dto2.DeleteByQueryReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsDeleteByQuery(ctx,
		reqData.DeleteByQueryReqData.IndexNames,
		reqData.DeleteByQueryReqData.Documents,
		reqData.DeleteByQueryReqData.Body,
	)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsSnapshotCreate(ctx *gin.Context) {
	var reqData dto2.SnapshotCreateReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSnapshotCreate(ctx,
		reqData.SnapshotCreateReqData.Repository,
		reqData.SnapshotCreateReqData.Snapshot,
		reqData.SnapshotCreateReqData.WaitForCompletion,
		reqData.SnapshotCreateReqData.ReqJson,
	)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsSnapshotDelete(ctx *gin.Context) {
	var reqData dto2.SnapshotDeleteReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSnapshotDelete(ctx,
		reqData.SnapshotDeleteReqData.Repository,
		reqData.SnapshotDeleteReqData.Snapshot)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsRestoreSnapshot(ctx *gin.Context) {
	var reqData dto2.RestoreSnapshotReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsRestoreSnapshot(ctx,
		reqData.RestoreSnapshotReqData.Repository,
		reqData.RestoreSnapshotReqData.Snapshot,
		reqData.RestoreSnapshotReqData.WaitForCompletion,
		reqData.RestoreSnapshotReqData.ReqJson,
	)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsSnapshotStatus(ctx *gin.Context) {
	var reqData dto2.SnapshotStatusReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSnapshotStatus(ctx,
		reqData.SnapshotStatusReqData.Repository,
		reqData.SnapshotStatusReqData.Snapshot,
		reqData.SnapshotStatusReqData.IgnoreUnavailable)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsSnapshotGetRepository(ctx *gin.Context) {
	var reqData dto2.SnapshotGetRepositoryReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSnapshotGetRepository(ctx,
		reqData.SnapshotGetRepositoryReqData.Repository)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsSnapshotCreateRepository(ctx *gin.Context) {
	var reqData dto2.SnapshotCreateRepositoryReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSnapshotCreateRepository(ctx,
		reqData.SnapshotCreateRepositoryReqData.Repository,
		reqData.SnapshotCreateRepositoryReqData.ReqJson)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsSnapshotDeleteRepository(ctx *gin.Context) {
	var reqData dto2.SnapshotDeleteRepositoryReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSnapshotDeleteRepository(ctx,
		reqData.SnapshotDeleteRepositoryReqData.Repository)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsPerformRequest(ctx *gin.Context) {
	var reqData dto2.PerformRequest

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	req := reqData.Request
	request, err := http.NewRequest(req.Method, req.URL.Path, bytes.NewReader([]byte(req.JsonBody)))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	request.Header = req.Header
	request.PostForm = req.PostForm
	request.Form = req.Form
	request.MultipartForm = req.MultipartForm
	request.Method = req.Method
	request.URL = req.URL

	res, err := esI.EsPerformRequest(ctx, request)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsGetIndices(ctx *gin.Context) {
	var reqData dto2.GetIndicesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsGetIndices(ctx,
		reqData.GetIndicesReqData.CatIndicesRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsCatHealth(ctx *gin.Context) {
	var reqData dto2.CatHealthReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCatHealth(ctx,
		reqData.CatHealthReqData.CatRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsCatShards(ctx *gin.Context) {
	var reqData dto2.CatShardsReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCatShards(ctx,
		reqData.CatShardsReqData.CatRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsCatCount(ctx *gin.Context) {
	var reqData dto2.CatCountReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCatCount(ctx,
		reqData.CatCountReqData.CatRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsCatAllocationRequest(ctx *gin.Context) {
	var reqData dto2.CatAllocationRequest

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCatAllocationRequest(ctx,
		reqData.CatAllocationRequestData.CatRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsCatAliases(ctx *gin.Context) {
	var reqData dto2.CatAliasesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCatAliases(ctx,
		reqData.CatAliasesReqData.CatRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsDelete(ctx *gin.Context) {
	var reqData dto2.DeleteReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsDelete(ctx,
		reqData.DeleteReqData.DeleteRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsUpdate(ctx *gin.Context) {
	var reqData dto2.UpdateReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsUpdate(ctx,
		reqData.UpdateReqData.UpdateRequest,
		reqData.UpdateReqData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsCreate(ctx *gin.Context) {
	var reqData dto2.CreateReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCreate(ctx,
		reqData.CreateReqData.CreateRequest,
		reqData.CreateReqData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsSearch(ctx *gin.Context) {
	var reqData dto2.SearchReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsSearch(ctx,
		reqData.SearchReqData.SearchRequest,
		reqData.SearchReqData.Query)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsIndicesPutSettingsRequest(ctx *gin.Context) {
	var reqData dto2.IndicesPutSettingsRequest

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsIndicesPutSettingsRequest(ctx,
		reqData.IndicesPutSettingsRequestData.IndexSettingsRequest,
		reqData.IndicesPutSettingsRequestData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsCreateIndex(ctx *gin.Context) {
	var reqData dto2.CreateIndexReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsCreateIndex(ctx,
		reqData.CreateIndexReqData.IndexCreateRequest,
		reqData.CreateIndexReqData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsDeleteIndex(ctx *gin.Context) {
	var reqData dto2.DeleteIndexReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsDeleteIndex(ctx,
		reqData.DeleteIndexReqData.IndicesDeleteRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsReindex(ctx *gin.Context) {
	var reqData dto2.ReindexReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsReindex(ctx,
		reqData.ReindexReqData.ReindexRequest,
		reqData.ReindexReqData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsIndicesGetSettingsRequest(ctx *gin.Context) {
	var reqData dto2.IndicesGetSettingsRequestReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsIndicesGetSettingsRequest(ctx,
		reqData.IndicesGetSettingsRequestReqData.IndicesGetSettingsRequest)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsPutMapping(ctx *gin.Context) {
	var reqData dto2.PutMappingReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsPutMapping(ctx,
		reqData.PutMappingReqData.IndicesPutMappingRequest,
		reqData.PutMappingReqData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsGetMapping(ctx *gin.Context) {
	var reqData dto2.GetMappingReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsGetMapping(ctx,
		reqData.GetMappingReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsGetAliases(ctx *gin.Context) {
	var reqData dto2.GetAliasesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsGetAliases(ctx,
		reqData.GetAliasesReqData.IndexNames)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsAddAliases(ctx *gin.Context) {
	var reqData dto2.AddAliasesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsAddAliases(ctx,
		reqData.AddAliasesReqData.IndexName,
		reqData.AddAliasesReqData.AliasName)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsRemoveAliases(ctx *gin.Context) {
	var reqData dto2.RemoveAliasesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsRemoveAliases(ctx,
		reqData.RemoveAliasesReqData.IndexName,
		reqData.RemoveAliasesReqData.AliasName)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsMoveToAnotherIndexAliases(ctx *gin.Context) {
	var reqData dto2.MoveToAnotherIndexAliasesReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsMoveToAnotherIndexAliases(ctx,
		reqData.MoveToAnotherIndexAliasesReqData.Body)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.OperateSuccess, res)
}

func (this *PluginUtilController) EsTaskList(ctx *gin.Context) {
	var reqData dto2.TaskListReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsTaskList(ctx)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) EsTasksCancel(ctx *gin.Context) {
	var reqData dto2.TasksCancelReq

	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	res, err := esI.EsTasksCancel(ctx,
		reqData.TasksCancelReqData.TaskId)
	if err != nil {
		this.ErrorProtobuf(ctx, errors.WithStack(err))
		return
	}

	this.SuccessProtobuf(ctx, response.SearchSuccess, res)
}

// 进行增删改等操作
func (this *PluginUtilController) MysqlExecSql(ctx *gin.Context) {
	var reqData dto2.MysqlExecReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	rowsAffected, err := esI.MysqlExecSql(ctx, reqData.DbName, reqData.Sql, reqData.Args...)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res := vo2.MysqlExecSqlRes{
		RowsAffected: rowsAffected,
	}

	this.Success(ctx, response.OperateSuccess, res)
}

// 进行查询操作
func (this *PluginUtilController) MysqlSelectSql(ctx *gin.Context) {
	var reqData dto2.MysqlSelectReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := esI.MysqlSelectSql(ctx, reqData.DbName, reqData.Sql, reqData.Args...)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, vo2.MysqlSelectSqlRes{Result: res})
}

// 进行查询操作
func (this *PluginUtilController) MysqlFirstSql(ctx *gin.Context) {
	var reqData dto2.MysqlSelectReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := esI.MysqlFirstSql(ctx, reqData.DbName, reqData.Sql, reqData.Args...)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, vo2.MysqlFirstSqlRes{Result: res})
}

func (this *PluginUtilController) RedisExecCommand(ctx *gin.Context) {
	var reqData dto2.RedisExecReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}

	res, err := esI.RedisExecCommand(ctx, reqData.DbName, reqData.Args...)
	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}

	this.SuccessProtobufByAny(ctx, response.OperateSuccess, map[string]interface{}{"data": res})
}

func (this *PluginUtilController) MongoExecCommand(ctx *gin.Context) {
	var reqData dto2.MongoExecReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}

	res, err := esI.ExecMongoCommand(ctx, reqData.DbName, reqData.Command, reqData.Timeout)
	if err != nil {
		this.ErrorProtobuf(ctx, err)
		return
	}

	this.SuccessProtobufByAny(ctx, response.SearchSuccess, res)
}

func (this *PluginUtilController) ShowMongoDbs(ctx *gin.Context) {
	var reqData dto2.ShowMongoDbsReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(ctx, reqData.EsConnectData.EsConnect, reqData.EsConnectData.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, reqData.EsConnectData.EsConnect))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := esI.ShowMongoDbs(ctx)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}
