package middleware

import (
	"bytes"
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

func (this *MiddleWareService) GinZapLoggerfunc(c *gin.Context) {
	start := time.Now()

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	c.Next()

	end := time.Now()
	latency := end.Sub(start).String()

	method := c.Request.Method
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

	if raw != "" {
		path = path + "?" + raw
	}

	this.log.Info("HTTP Request",
		zap.Int("status", statusCode),
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", clientIP),
		zap.String("latency", latency),
		zap.String("error", errorMessage),
	)
}

func (this *MiddleWareService) GinZapRecovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 打印错误信息

			this.log.Error("panic recovered",
				zap.Any("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.Stack("stacktrace"),
			)
			this.res.Error(c, errors.New(cast.ToString(err)))

		}
	}()
	c.Next()
}
