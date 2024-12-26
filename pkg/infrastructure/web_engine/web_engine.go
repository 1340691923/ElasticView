package web_engine

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

var authenticationPaths []string

type RouterConfigGroup struct {
	GroupRemark   string         `json:"label"`
	RouterConfigs []RouterConfig `json:"options"`
}

type RouterConfig struct {
	Url      string `json:"value"`
	Remark   string `json:"label"`
	NeedAuth bool   `json:"needAuth"`
}

type WebEngine struct {
	g *gin.Engine
}

func NewWebEngine() *WebEngine {
	gin.SetMode(gin.ReleaseMode)
	return &WebEngine{g: gin.Default()}
}

func (this *WebEngine) GetGinEngine() *gin.Engine {
	return this.g
}

func (this *WebEngine) Run() {
	this.g.Run()
}

func (this *WebEngine) Group(remark string, relativePath string, handlers ...gin.HandlerFunc) *MyRouterGroup {
	rg := this.g.Group(relativePath, handlers...)
	return NewMyRouterGroup(rg, relativePath, remark)
}

func (this *WebEngine) GetAuthenticationPaths() []string {
	return authenticationPaths
}

func (this *WebEngine) GetRouterConfigGroups() []RouterConfigGroup {
	routerConfigGroups := []RouterConfigGroup{}

	mockMap.Range(func(key string, val []RouterConfig) {
		routerConfigGroups = append(routerConfigGroups, RouterConfigGroup{
			GroupRemark:   key,
			RouterConfigs: val,
		})
	})

	return routerConfigGroups
}

type MyRouterGroup struct {
	rg                *gin.RouterGroup
	groupRelativePath string
	groupTag          string
	routerConfigs     []RouterConfig
}

func NewMyRouterGroup(rg *gin.RouterGroup, relativePath, remark string) *MyRouterGroup {
	return &MyRouterGroup{
		rg:                rg,
		groupTag:          remark,
		groupRelativePath: relativePath,
	}
}

func (this *MyRouterGroup) Use(middleware ...gin.HandlerFunc) *MyRouterGroup {
	this.rg.Handlers = append(this.rg.Handlers, middleware...)
	return this
}

func (this *MyRouterGroup) Group(remark, relativePath string, handlers ...gin.HandlerFunc) *MyRouterGroup {
	rg := this.rg.Group(relativePath, handlers...)
	return NewMyRouterGroup(rg, relativePath, remark)
}

func (this *MyRouterGroup) BasePath() string {
	return this.rg.BasePath()
}

func (this *MyRouterGroup) Handle(needAuth bool, httpMethod, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.Handle(httpMethod, relativePath, handlers...)
}

func (this *MyRouterGroup) POST(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.POST(relativePath, handlers...)
}

func (this *MyRouterGroup) GET(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.GET(relativePath, handlers...)
}

func (this *MyRouterGroup) DELETE(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.DELETE(relativePath, handlers...)
}

func (this *MyRouterGroup) PATCH(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.PATCH(relativePath, handlers...)
}

func (this *MyRouterGroup) PUT(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.PUT(relativePath, handlers...)
}

func (this *MyRouterGroup) OPTIONS(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.OPTIONS(relativePath, handlers...)
}

func (this *MyRouterGroup) HEAD(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.HEAD(relativePath, handlers...)
}

func (this *MyRouterGroup) Any(needAuth bool, remark, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	this.saveMockMap(remark, relativePath, needAuth)
	return this.rg.Any(relativePath, handlers...)
}

func (this *MyRouterGroup) StaticFile(relativePath, filepath string) gin.IRoutes {
	return this.rg.StaticFile(relativePath, filepath)
}

func (this *MyRouterGroup) Static(relativePath, root string) gin.IRoutes {
	return this.rg.Static(relativePath, root)
}

func (this *MyRouterGroup) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return this.rg.StaticFS(relativePath, fs)
}

func (this *MyRouterGroup) joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := this.lastChar(relativePath) == '/' && this.lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func (this *MyRouterGroup) lastChar(str string) uint8 {
	if str == "" {
		panic("路由匹配符不能为空")
	}
	return str[len(str)-1]
}

func (this *MyRouterGroup) saveMockMap(remark, relativePath string, needAuth bool) {
	path := this.joinPaths(this.groupRelativePath, relativePath)
	authenticationPaths = append(authenticationPaths, path)

	mockMap.Store(this.groupTag, RouterConfig{
		Url:      path,
		Remark:   remark,
		NeedAuth: needAuth,
	})
}
