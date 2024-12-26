package middleware

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/my_error"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MiddleWareService struct {
	cfg              *config.Config
	orm              *orm.Gorm
	log              *logger.AppLogger
	jwtSvr           *jwt_svr.Jwt
	res              *response.Response
	gmUserService    *gm_user.GmUserService
	rbac             *access_control.Rbac
	routerEngine     *web_engine.WebEngine
	pluginRegistry   manager.Service
	gmOperaterLogSvr *gm_operater_log.GmOperaterLogService
}

func NewMiddleWareService(cfg *config.Config, orm *orm.Gorm, log *logger.AppLogger, jwtSvr *jwt_svr.Jwt, res *response.Response, gmUserService *gm_user.GmUserService, rbac *access_control.Rbac, routerEngine *web_engine.WebEngine, pluginRegistry manager.Service, gmOperaterLogSvr *gm_operater_log.GmOperaterLogService) *MiddleWareService {
	return &MiddleWareService{cfg: cfg, orm: orm, log: log, jwtSvr: jwtSvr, res: res, gmUserService: gmUserService, rbac: rbac, routerEngine: routerEngine, pluginRegistry: pluginRegistry, gmOperaterLogSvr: gmOperaterLogSvr}
}

func (this *MiddleWareService) JwtMiddleware(c *gin.Context) {

	var err error
	var claims *jwt_svr.Claims
	token := util.GetToken(c)

	if _, logoff := util.TokenBucket.Load(token); logoff {
		err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_CHECK_TOKEN_FAIL)
		this.res.Error(c, err)
		c.Abort()
		return
	}
	if util.GetToken(c) == "" {
		err = my_error.NewBusiness(TOKEN_ERROR, INVALID_PARAMS)
		this.res.Error(c, err)
		c.Abort()
		return
	}

	service := this.gmUserService
	claims, err = this.jwtSvr.ParseToken(token)
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_CHECK_TOKEN_TIMEOUT)

			newToken, _ := this.jwtSvr.CreateTokenByCliams(*claims)

			this.res.Output(c.Writer, map[string]interface{}{
				"code":     ERROR_AUTH_CHECK_TOKEN_TIMEOUT,
				"msg":      err.Error(),
				"newToken": newToken,
			})
			c.Abort()
			return
		} else {
			err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_CHECK_TOKEN_FAIL)
			this.res.Error(c, err)
			c.Abort()
			return
		}
	}

	isExsitUser, err := service.IsExsitUser(c, claims)

	if err != nil {
		this.res.Error(c, errors.WithStack(err))
		c.Abort()
		return
	}

	if !isExsitUser {
		this.res.Error(c, my_error.NewBusiness(TOKEN_ERROR, ERROR_CEHCK_USER_EXITS))
		c.Abort()
		return
	}

	c.Set("userName", claims.Username)
	c.Set("userId", claims.UserID)
	//c.Set("roleId", claims.RoleId)

	c.Next()
	return
}

func (this *MiddleWareService) OperaterLog(ctx *gin.Context) {

	var err error
	token := util.GetToken(ctx)
	var claims *jwt_svr.Claims
	claims, err = this.jwtSvr.ParseToken(token)
	if err != nil {
		this.log.Error("OperaterLog jwt err", zap.Error(err))
		ctx.Abort()
		return
	}
	contentType := ctx.GetHeader("Content-Type")

	var body []byte

	if strings.Contains(contentType, "application/json") {
		var b []byte
		b, _ = ctx.GetRawData()

		body, err = util.GzipCompress(util.Bytes2str(b))
		if err != nil {
			return
		}

		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}
	startT := time.Now()

	ctx.Next()

	if strings.Contains(contentType, "application/json") {
		costT := time.Now().Sub(startT).String()
		resStatus := "失败"
		if ctx.Writer.Status() == 200 {
			resStatus = "正常"
		}
		gmOperaterLog := &model.GmOperaterLog{
			OperaterName:   cast.ToString(claims.Username),
			OperaterId:     cast.ToInt(claims.UserID),
			OperaterAction: ctx.Request.URL.Path,
			Method:         ctx.Request.Method,
			Body:           body,
			Created:        time.Now(),
			CostTime:       costT,
			Status:         resStatus,
		}
		this.gmOperaterLogSvr.Save(gmOperaterLog)
	}

	return
}

func (this *MiddleWareService) Rbac(ctx *gin.Context) {
	var err error
	token := util.GetToken(ctx)
	var claims *jwt_svr.Claims
	claims, err = this.jwtSvr.ParseToken(token)

	if err != nil {
		this.log.Error("Rbac ", zap.Error(err))
		return
	}
	obj := ctx.Request.RequestURI

	roleids, err := this.gmUserService.GetRolesByUserID(claims.UserID)
	if err != nil {
		this.log.Error("Rbac ", zap.Error(err))
		return
	}

	ctx.Set("ev_roles", roleids)

	sort.Ints(roleids)

	for _, rg := range this.routerEngine.GetRouterConfigGroups() {
		for _, routerConfig := range rg.RouterConfigs {
			if obj == routerConfig.Url && routerConfig.NeedAuth {
				for _, roleId := range roleids {
					ok, err := this.rbac.Enforce(strconv.Itoa(roleId), obj, "*")
					if err != nil {
						this.res.Error(ctx, my_error.NewBusiness(TOKEN_ERROR, ERROR_RBAC_LOAD))
						ctx.Abort()
						return
					}
					if !ok {
						this.res.Error(ctx, errors.New(fmt.Sprintf("您没有操作该资源的权限:%s", routerConfig.Remark)))
						ctx.Abort()
						return
					}
					break
				}
			}
		}
	}

	ctx.Next()
	return
}

func (this *MiddleWareService) ValidatePluginSign(ctx *gin.Context) {

	pluginID := ctx.GetHeader("X-Plugin-ID")
	pluginSign := ctx.GetHeader("X-Plugin-Signature")
	p, hasPlugin := this.pluginRegistry.Plugin(context.Background(), pluginID)
	if !hasPlugin {
		this.res.Error(ctx, errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID)))
		ctx.Abort()
		return
	}
	if !p.PluginData().PluginJsonData.BackendDebug {
		var b []byte
		b, _ = ctx.GetRawData()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

		if !verifySignature(pluginSign, p.SignKey, string(b)) {
			this.res.Error(ctx, errors.New(fmt.Sprintf("插件签名校验失败:%s", pluginID)))
			ctx.Abort()
			return
		}
	}

	ctx.Next()
	return
}

func verifySignature(signature, signatureKey, jsonString string) bool {
	// 计算 HMAC-SHA256 签名
	mac := hmac.New(sha256.New, []byte(signatureKey))
	mac.Write([]byte(jsonString))
	expectedSignatureBytes := mac.Sum(nil)

	// 将期望的签名转换为 Base64 编码的字符串
	expectedSignature := base64.StdEncoding.EncodeToString(expectedSignatureBytes)

	// 比较期望的签名与实际签名是否相同
	return signature == expectedSignature
}

func (this *MiddleWareService) CheckVersion(c *gin.Context) {
	if !this.cfg.DeBug && config.GetVersion() != c.GetHeader("X-Version") && c.GetHeader("X-Version") != "test" {
		err := my_error.NewError(fmt.Sprintf("后台已更新版本（新版本：%s,您的版本：%s），请刷新页面", config.GetVersion(), c.GetHeader("X-Version")), ERROR_CEHCK_VERSION_FAIL)
		this.res.Error(c, err)
		c.Abort()
		return
	}

	// 处理请求
	c.Next()
}

/*func (this *MiddleWareService) Timer(ctx *gin.Context) {

	// start timer
	start := time.Now()
	// next routes
	err := ctx.Next()
	// stop timer
	stop := time.Now()

	this.log.Info("时间拦截器",
		zap.String("访问资源", ctx.Path()),
		zap.Reflect("body", string(this.getPostBody(ctx))),
		zap.String("消耗时间：", stop.Sub(start).String()))
	return err

}
*/
