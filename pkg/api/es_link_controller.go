package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/services/cache_service"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"time"

	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
)

// Es 连接管理控制器
type EsLinkController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	sqlx            *orm.Gorm
	esLinkService   *es_link_service.EsLinkService
	jwtSvr          *jwt_svr.Jwt
	esCache         *cache_service.EsCache
}

func NewEsLinkController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, sqlx *orm.Gorm, esLinkService *es_link_service.EsLinkService, jwtSvr *jwt_svr.Jwt, esCache *cache_service.EsCache) *EsLinkController {
	return &EsLinkController{BaseController: baseController, log: log, esClientService: esClientService, sqlx: sqlx, esLinkService: esLinkService, jwtSvr: jwtSvr, esCache: esCache}
}

// @Summary 获取Es连接列表
// @Tags es连接信息
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Success 0 {object} []vo.EsLink
// @Router /api/es_link/ListAction [post]
func (this *EsLinkController) ListAction(ctx *gin.Context) {

	type ReqData struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	var reqData ReqData
	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	list, count, err := this.esLinkService.GetListAction(ctx, userInfo.UserID, this.GetRoleCache(ctx), reqData.Page, reqData.PageSize)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"count": count,
		"list":  list,
	})
}

func (this *EsLinkController) GetEsCfgList(ctx *gin.Context) {

	type ReqData struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	var reqData ReqData
	err := ctx.BindJSON(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	list, count, err := this.esLinkService.GetEsCfgList(ctx, userInfo.UserID, reqData.Page, reqData.PageSize)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"list":  list,
		"count": count,
	})
}

func (this *EsLinkController) GetEsCfgOpt(ctx *gin.Context) {

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	list, err := this.esLinkService.GetEsCfgOpt(ctx, userInfo.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
}

func (this *EsLinkController) DeleteEsCfgRelation(ctx *gin.Context) {

	var req dto.GetEsCfgRelation
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.esLinkService.DeleteEsCfgRelation(ctx, this.sqlx.DB, req.ID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
}

// @Summary 查看ES连接配置下拉选
// @Tags es连接信息
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Success 0 {object} []vo.EsLinkOpt
// @Router /api/es_link/OptAction [post]
func (this *EsLinkController) OptAction(ctx *gin.Context) {

	roles := this.GetRoleCache(ctx)

	optList, err := this.esClientService.GetEsLinkOptions(ctx, roles)

	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}

	this.Success(ctx, response.SearchSuccess, optList)
}

// @Summary 新增连接信息
// @Tags es连接信息
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body dto.InsertEsLink false "查询参数"
// @Success 0 {object} response.ResponseData
// @Router /api/es_index/InsertAction [post]
func (this *EsLinkController) InsertAction(ctx *gin.Context) {

	var reqData dto.InsertEsLink
	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Ip == "" {
		this.Error(ctx, errors.New("连接信息不能为空"))
		return
	}
	if reqData.Remark == "" {
		reqData.Remark = reqData.Ip
	}
	_, ok := factory.EsServiceMap[reqData.Version]
	if !ok {
		this.Error(ctx, factory.VersionErr())
		return
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	tx := this.sqlx.Begin()

	es_link_id, err := this.esLinkService.SaveEsLink(ctx, tx, &model.EsLinkV2{
		Ip:       reqData.Ip,
		Created:  time.Now(),
		Updated:  time.Now(),
		Remark:   reqData.Remark,
		Version:  reqData.Version,
		CreateBy: userInfo.UserID,
	})

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	for _, cfgId := range reqData.CfgIds {
		err = this.esLinkService.SaveEslinkRoleCfgByEsLinkId(ctx, tx, cfgId, es_link_id)
		if err != nil {
			tx.Rollback()
			this.Error(ctx, err)
			return
		}
	}
	tx.Commit()
	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// @Summary 修改连接信息
// @Tags es连接信息
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body dto.UpdateEsLink false "查询参数"
// @Success 0 {object} response.ResponseData
// @Router /api/es_index/UpdateAction [post]
func (this *EsLinkController) UpdateAction(ctx *gin.Context) {
	var reqData dto.UpdateEsLink
	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	_, ok := factory.EsServiceMap[reqData.Version]
	if !ok {
		this.Error(ctx, factory.VersionErr())
		return
	}

	if reqData.Ip == "" {
		this.Error(ctx, errors.New("连接信息不能为空"))
		return
	}

	if reqData.Remark == "" {
		reqData.Remark = reqData.Ip
	}

	if reqData.Id <= 0 {
		this.Error(ctx, errors.New("ID不能为空"))
		return
	}

	updateMap := map[string]interface{}{}

	updateMap["ip"] = reqData.Ip
	updateMap["remark"] = reqData.Remark
	updateMap["version"] = reqData.Version
	updateMap["updated"] = time.Now().Format(util.TimeFormat)

	tx := this.sqlx.Begin()

	err = this.esLinkService.UpdateEsLink(ctx, tx, updateMap, reqData.Id)

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	err = this.esLinkService.DeleteEsCfgRelationByEsLinkId(ctx, tx, reqData.Id)

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	for _, cfgId := range reqData.CfgIds {

		err = this.esLinkService.SaveEslinkRoleCfgByEsLinkId(ctx, tx, cfgId, reqData.Id)

		if err != nil {
			tx.Rollback()
			this.Error(ctx, err)
			return
		}
	}
	tx.Commit()
	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// @Summary 删除es连接
// @Tags es连接信息
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body dto.DeleteEsLink false "查询参数"
// @Success 0 {object} response.ResponseData
// @Router /api/es_index/DeleteAction [post]
func (this *EsLinkController) DeleteAction(ctx *gin.Context) {

	var req dto.DeleteEsLink

	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if req.Id <= 0 {
		this.Error(ctx, errors.New("id不能为空"))
		return
	}

	tx := this.sqlx.Begin()

	err = this.esLinkService.DeleteById(ctx, tx, req.Id)

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	err = this.esLinkService.DeleteEsCfgRelationByEsLinkId(ctx, tx, req.Id)
	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}
	tx.Commit()
	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.DeleteSuccess, nil)
}

func (this *EsLinkController) InsertEsCfgAction(ctx *gin.Context) {

	var reqData dto.InsertEsLinkCfg
	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Pwd != "" {
		reqData.Pwd, err = this.esClientService.EsPwdESBEncrypt(ctx, reqData.Pwd)
		if err != nil {
			this.Error(ctx, err)
			return
		}
	}

	userInfo, _ := this.jwtSvr.ParseToken(this.GetToken(ctx))

	for _, header := range reqData.Header {
		if header.Key == "" {
			this.Error(ctx, errors.New("请求头中key不能为空"))
			return
		}
		if header.Value == "" {
			this.Error(ctx, errors.New("请求头中value不能为空"))
			return
		}
	}
	header, err := json.Marshal(reqData.Header)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	es_link_cfg_id, err := this.esLinkService.SaveEsLinkCfgV2(ctx, &model.EslinkCfgV2{
		User:     reqData.User,
		Pwd:      reqData.Pwd,
		Rootpem:  &reqData.RootPEM,
		Certpem:  &reqData.CertPEM,
		Keypem:   &reqData.KeyPEM,
		Created:  time.Now(),
		Updated:  time.Now(),
		CreateBy: userInfo.UserID,
		Header:   string(header),
		Remark:   reqData.Remark,
	})
	if err != nil {
		this.Error(ctx, err)
		return
	}

	for _, shareRole := range reqData.ShareRoles {

		err = this.sqlx.Exec("insert into gm_role_eslink_cfg_v2(role_id,es_link_cfg_id)"+
			" values (?,?)", shareRole, es_link_cfg_id).Error

		if err != nil {
			this.Error(ctx, err)
			return
		}

	}

	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

func (this *EsLinkController) UpdateEsCfgAction(ctx *gin.Context) {

	var reqData dto.UpdateEsLinkCfg
	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Id <= 0 {
		this.Error(ctx, errors.New("id不能为空"))
		return
	}

	for _, header := range reqData.Header {
		if header.Key == "" {
			this.Error(ctx, errors.New("请求头中key不能为空"))
			return
		}
		if header.Value == "" {
			this.Error(ctx, errors.New("请求头中value不能为空"))
			return
		}
	}
	header, err := json.Marshal(reqData.Header)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	updateMap := map[string]interface{}{}
	updateMap["user"] = reqData.User
	if reqData.Pwd != "" {
		updateMap["pwd"], err = this.esClientService.EsPwdESBEncrypt(ctx, reqData.Pwd)
		if err != nil {
			this.Error(ctx, err)
			return
		}
	} else {
		updateMap["pwd"] = reqData.Pwd
	}
	updateMap["remark"] = reqData.Remark
	updateMap["rootpem"] = reqData.RootPEM
	updateMap["certpem"] = reqData.CertPEM
	updateMap["keypem"] = reqData.KeyPEM
	updateMap["created"] = time.Now().Format(util.TimeFormat)
	updateMap["updated"] = time.Now().Format(util.TimeFormat)
	updateMap["header"] = string(header)

	tx := this.sqlx.Begin()

	err = this.esLinkService.UpdateEsLinkCfgById(ctx, tx, updateMap, reqData.Id)

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	err = this.esLinkService.DeleteRoleEslinkCfgByEsLinkCfgId(ctx, tx, reqData.Id)

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	for _, shareRole := range reqData.ShareRoles {
		err = tx.Exec("insert into gm_role_eslink_cfg_v2(role_id,es_link_cfg_id)"+
			" values (?,?)", shareRole, reqData.Id).Error
		if err != nil {
			tx.Rollback()
			this.Error(ctx, err)
			return
		}
	}
	tx.Commit()
	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

func (this *EsLinkController) DeleteEsCfgAction(ctx *gin.Context) {

	var reqData dto.DeleteEsLinkCfg
	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Id <= 0 {
		this.Error(ctx, errors.New("id不能为空"))
		return
	}

	err = this.esLinkService.DeleteEsCfg(ctx, reqData.Id)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	cache.CleanDataSourceCache(true)
	this.esCache.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
	return
}
