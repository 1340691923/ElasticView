package api

import (
	"fmt"
	dto2 "github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/services/eve_service"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/services/plugin_install_service"
	"github.com/1340691923/ElasticView/pkg/services/plugin_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
)

type PluginController struct {
	*BaseController
	log             *logger.AppLogger
	orm             *orm.Gorm
	pluginService   *plugin_service.PluginService
	eveService      *eve_service.EvEService
	pluginInstaller *plugin_install_service.PluginInstaller
	gmUserService   *gm_user.GmUserService
	jwtSvr          *jwt_svr.Jwt
}

func NewPluginController(baseController *BaseController, log *logger.AppLogger, orm *orm.Gorm, pluginService *plugin_service.PluginService, eveService *eve_service.EvEService, pluginInstaller *plugin_install_service.PluginInstaller, gmUserService *gm_user.GmUserService, jwtSvr *jwt_svr.Jwt) *PluginController {
	return &PluginController{BaseController: baseController, log: log, orm: orm, pluginService: pluginService, eveService: eveService, pluginInstaller: pluginInstaller, gmUserService: gmUserService, jwtSvr: jwtSvr}
}

func (this *PluginController) CallPlugin(ctx *gin.Context) {
	err := this.pluginService.CallPlugin(ctx, ctx.Param("plugin_id"))
	if err != nil {
		this.Error(ctx, err)
		return
	}
}

func (this *PluginController) CallPluginViews(ctx *gin.Context) {
	err := this.pluginService.CallPluginViews(ctx, ctx.Param("plugin_id"))
	if err != nil {
		this.Error(ctx, err)
		return
	}
}

func (this *PluginController) PluginMarket(ctx *gin.Context) {

	var req dto.FromEvPluginReq

	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	pluginList, err := this.eveService.GetRemotePlugins(ctx, &req)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, pluginList)

}

func (this *PluginController) GetPluginInfo(ctx *gin.Context) {
	var req dto.FormEvPluginInfoReq

	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.eveService.GetRemotePluginInfo(ctx, &req)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)

}

func (this *PluginController) InstallPlugin(ctx *gin.Context) {
	var req dto2.InstallPlugin
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginInstaller.Add(ctx, req.PluginID, req.Version)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	userInfo, err := this.jwtSvr.ParseToken(this.GetToken(ctx))

	if err != nil {
		this.Error(ctx, err)
		return
	}

	roles, err := this.gmUserService.GetRolesByUserID(userInfo.UserID)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if util.InArr(roles, 1) {
		this.Success(ctx, "安装成功,请刷新", nil)
		return
	}

	this.Success(ctx, "安装成功", nil)
}

func (this *PluginController) StarPlugin(ctx *gin.Context) {
	var req dto.StarPlugin
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginService.StarPlugin(ctx, req.PluginId)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "操作成功", nil)
}

func (this *PluginController) AddComment(ctx *gin.Context) {
	var req dto.AddCommentRequest
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginService.AddComment(ctx, req.PluginID, req.Content, req.ParentID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "操作成功", nil)
}

func (this *PluginController) LikeComment(ctx *gin.Context) {
	var req dto.LikeCommentRequest
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginService.LikeComment(ctx, req.CommentID, req.State)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "操作成功", nil)
}

func (this *PluginController) ListComments(ctx *gin.Context) {
	var req dto.ListCommentsRequest
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	list, err := this.pluginService.ListComments(ctx, req.PluginID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
}

func (this *PluginController) UnInstallPlugin(ctx *gin.Context) {
	var req dto2.InstallPlugin
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginInstaller.Remove(ctx, req.PluginID, "")
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "卸载成功", nil)
}

func (this *PluginController) GetWxArticleList(ctx *gin.Context) {
	list, err := this.eveService.GetWxArticleList(ctx)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
}

func (this *PluginController) GetLocalPluginList(ctx *gin.Context) {
	plugins := this.pluginService.PluginList(ctx)

	this.Success(ctx, response.SearchSuccess, plugins)
}

func (this *PluginController) ImportEvKey(ctx *gin.Context) {
	var req dto2.ImportEvKey
	err := ctx.Bind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.eveService.SaveEvKey(req.EvKey)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.eveService.FlushAccessToken(ctx)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
}

func (this *PluginController) UploadPlugin(ctx *gin.Context) {
	f, err := ctx.FormFile("file")
	if err != nil {
		this.Error(ctx, err)
		return
	}

	var pluginId string
	/*os := runtime.GOOS
	switch os {
	case "windows":
		pluginId, err = this.pluginInstaller.AddUploadPlugin(ctx, f)
	default:
		pluginId, err = this.pluginInstaller.AddUploadPluginByNoWindows(ctx, f)
	}*/
	pluginId, err = this.pluginInstaller.AddUploadPlugin(ctx, f)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, fmt.Sprintf("%s安装成功", pluginId), nil)
}
