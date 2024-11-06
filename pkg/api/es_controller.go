package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto/common"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_service"

	"github.com/gin-gonic/gin"
)

// Es 基本操作
type EsController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	esService       *es_service.EsService
	jwtSvr          *jwt_svr.Jwt
	cfg             *config.Config
}

func NewEsController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, esService *es_service.EsService, jwtSvr *jwt_svr.Jwt, cfg *config.Config) *EsController {
	return &EsController{BaseController: baseController, log: log, esClientService: esClientService, esService: esService, jwtSvr: jwtSvr, cfg: cfg}
}

// @Summary 测试es连接
// @Tags ES
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body model.EsConnect false "查询参数"
// @Success 0 {object} response.ResponseData
// @Router /api/es/PingAction [post]
func (this *EsController) PingAction(ctx *gin.Context) {
	esConnect := new(dto.EsConnect)
	err := ctx.Bind(esConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	esConnect.Header = append(esConnect.Header, dto.HeaderKv{
		Key:   "ev_user_id",
		Value: cast.ToString(userInfo.UserID),
	})

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, esConnect.Id))
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

	this.Success(ctx, response.OperateSuccess, res)
}

// @Summary 得到所有的索引数量
// @Tags es索引
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body dto.EsConnectID false "查询参数"
// @Success 0 {object} int
// @Router /api/es_index/IndexsCountAction [post]
func (this *EsController) IndexsCountAction(ctx *gin.Context) {
	esConnectID := new(common.EsConnectID)
	err := ctx.Bind(&esConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	esConnect, err := this.esClientService.GetEsClientByID(ctx,
		esConnectID.EsConnectID, userInfo.UserID)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, esConnectID.EsConnectID))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.esService.EsIndexCount(ctx, esI)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, res)
}

func (this *EsController) CatAction(ctx *gin.Context) {

	esCat := new(dto.EsCat)
	err := ctx.Bind(&esCat)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	esConnect, err := this.esClientService.GetEsClientByID(ctx, esCat.EsConnect, userInfo.UserID)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg(this.cfg, esCat.EsConnect))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	catSvr := this.esService

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
	default:
		err = errors.New("未知类型")
	}

	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}

	if data.StatusErr() != nil {
		this.Error(ctx, errors.WithStack(data.StatusErr()))
		return
	}

	this.Success(ctx, response.SearchSuccess, data.JsonRawMessage())
}
