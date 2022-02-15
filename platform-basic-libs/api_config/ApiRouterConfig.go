package api_config

import (
	"path"
	"sync"

	fiber "github.com/gofiber/fiber/v2"
)

const (
	MethodAny     = "ANY"
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodOptions = "OPTIONS"
)

type RouterConfig struct {
	Url    string `json:"url"`
	Remark string `json:"remark"`
}

type ApiRouterConfig struct {
	routerConfigs               []RouterConfig
	noVerificationRouterConfigs []string
}

type MountApiBasePramas struct {
	Remark, Method, AbsolutePath, RelativePath string
}

func (this *MountApiBasePramas) IsAnyMethod() bool {
	return this.Method == MethodAny
}

func (this *ApiRouterConfig) GetRouterConfigs() (routerConfig []RouterConfig) {
	return this.routerConfigs
}

func (this *ApiRouterConfig) MountApi(mountApiBasePramas MountApiBasePramas, routerGroup *fiber.Group, authentication bool, handlers ...fiber.Handler) {
	if authentication {
		this.routerConfigs = append(this.routerConfigs, RouterConfig{
			Url:    this.joinPaths(mountApiBasePramas.AbsolutePath, mountApiBasePramas.RelativePath),
			Remark: mountApiBasePramas.Remark,
		})
	} else {
		this.noVerificationRouterConfigs = append(this.noVerificationRouterConfigs, this.joinPaths(mountApiBasePramas.AbsolutePath, mountApiBasePramas.RelativePath))
	}

	this.mountApi(mountApiBasePramas, routerGroup, handlers...)
}

func (this *ApiRouterConfig) mountApi(mountApiBasePramas MountApiBasePramas, routerGroup *fiber.Group, handlers ...fiber.Handler) {
	if mountApiBasePramas.IsAnyMethod() {
		routerGroup.All(mountApiBasePramas.RelativePath, handlers...)
	} else {
		routerGroup.Add(mountApiBasePramas.Method, mountApiBasePramas.RelativePath, handlers...)
	}
}

func (this *ApiRouterConfig) joinPaths(absolutePath, relativePath string) string {
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

func (this *ApiRouterConfig) lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

var apiRouterConfig *ApiRouterConfig
var once sync.Once

func NewApiRouterConfig() *ApiRouterConfig {
	once.Do(func() {
		apiRouterConfig = new(ApiRouterConfig)
	})
	return apiRouterConfig
}
