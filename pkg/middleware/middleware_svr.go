package middleware

import (
	"bytes"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"time"
)

type MiddleWareService struct {
	cfg           *config.Config
	sqlx          *sqlstore.SqlStore
	log           *logger.AppLogger
	jwtSvr        *jwt_svr.Jwt
	res           *response.Response
	gmUserService *gm_user.GmUserService
	rbac          *access_control.Rbac
	routerEngine  *web_engine.WebEngine
}

func NewMiddleWareService(cfg *config.Config, sqlx *sqlstore.SqlStore,
	log *logger.AppLogger, rbac *access_control.Rbac,
	jwtSvr *jwt_svr.Jwt, res *response.Response, gmUserService *gm_user.GmUserService,
	routerEngine *web_engine.WebEngine) *MiddleWareService {
	return &MiddleWareService{cfg: cfg, routerEngine: routerEngine, sqlx: sqlx, log: log, rbac: rbac, jwtSvr: jwtSvr, res: res, gmUserService: gmUserService}
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
		err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_CHECK_TOKEN_FAIL)
		this.res.Error(c, err)
		c.Abort()
		return
	}
	if time.Now().Unix() > cast.ToInt64(claims.ExpiresAt.Unix()) {
		err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_CHECK_TOKEN_TIMEOUT)
		this.res.Error(c, err)
		return
	}
	if !service.IsExitUser(claims) {
		err = my_error.NewBusiness(TOKEN_ERROR, ERROR_AUTH_CHECK_TOKEN_TIMEOUT)
		this.res.Error(c, err)
		c.Abort()
		return
	}

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

	var b []byte
	b, _ = ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	gmOperaterLog := &model.GmOperaterLog{
		OperaterName:   cast.ToString(claims.Username),
		OperaterId:     cast.ToInt(claims.UserID),
		OperaterAction: ctx.Request.URL.Path,
		Method:         ctx.Request.Method,
		Body:           b,
		OperaterRoleId: cast.ToInt(claims.RoleId),
		Sqlx:           this.sqlx,
	}

	err = gmOperaterLog.Insert()

	if err != nil {
		this.log.Error("OperaterLog", zap.Error(err))
	}
	ctx.Next()
	return
}

const ADMIN_ROLE = 1

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

	sub := int(claims.RoleId)

	//最高权限用户可免接口鉴权
	if sub == ADMIN_ROLE {
		ctx.Next()
		return
	}

	for _, rg := range this.routerEngine.GetRouterConfigGroups() {
		for _, routerConfig := range rg.RouterConfigs {
			if obj == routerConfig.Url && routerConfig.NeedAuth {
				ok, err := this.rbac.Enforce(strconv.Itoa(sub), obj, "*")
				if err != nil {
					this.res.Error(ctx, my_error.NewBusiness(TOKEN_ERROR, ERROR_RBAC_LOAD))
					ctx.Abort()
					return
				}
				if !ok {
					this.res.Error(ctx, my_error.NewBusiness(TOKEN_ERROR, ERROR_RBAC_AUTH))
					ctx.Abort()
					return
				}
			}
		}
	}

	ctx.Next()
	return
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
