//go:build wireinject
// +build wireinject

package server

import (
	"github.com/1340691923/ElasticView/pkg/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	log2 "github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/middleware"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm/migrator"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugin_rpc"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager/process"
	"github.com/1340691923/ElasticView/pkg/infrastructure/pluginstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/request"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/registry"
	"github.com/1340691923/ElasticView/pkg/server/backgroundsvcs"
	"github.com/1340691923/ElasticView/pkg/services/big_mode_service"
	"github.com/1340691923/ElasticView/pkg/services/cache_service"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/1340691923/ElasticView/pkg/services/es_service"
	"github.com/1340691923/ElasticView/pkg/services/eve_service"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/1340691923/ElasticView/pkg/services/gm_role"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/services/live_svr"
	"github.com/1340691923/ElasticView/pkg/services/oauth"
	"github.com/1340691923/ElasticView/pkg/services/plugin_install_service"
	"github.com/1340691923/ElasticView/pkg/services/plugin_service"
	"github.com/1340691923/ElasticView/pkg/services/updatechecker"
	"github.com/1340691923/ElasticView/pkg/web"

	"github.com/google/wire"
)

var wireSet = wire.NewSet(

	wire.Bind(new(registry.BackgroundServiceRegistry), new(*backgroundsvcs.BackgroundServiceRegistry)),
	live_svr.NewLive,
	api.NewWsController,
	oauth.ProvideOAuthServiceRegistry,
	oauth.NewWorkWechat,
	big_mode_service.NewBigMode,
	api.NewAiController,
	plugin_install_service.ProvideInstaller,
	wire.Bind(new(manager.Service), new(*manager.PluginManager)),
	api.NewPluginController,
	process.ProvideService,
	migrator.NewMigrator,
	config.InitConfig,
	manager.NewPluginManager,
	pluginstore.NewPluginStoreService,
	eve_api.NewEvApi,
	eve_service.NewEvEService,
	dao.NewEvBackDao,
	log2.InitLog,
	dao.NewEslinkCfgV2Dao,
	dao.NewGmRoleDao,
	dao.NewGmUserDao,
	dao.NewEsLinkV2Dao,
	dao.NewGmRoleEslinkCfgV2Dao,
	dao.NewEslinkRoleCfgReletion,

	updatechecker.ProvidePluginsService,
	orm.NewGorm,
	cache_service.NewEsCache,
	request.NewRequest,
	response.NewResponse,
	api.NewBaseController,

	api.NewIndexController,
	updatechecker.ProvideEvUpdate,
	api.NewPluginUtilController,
	plugin_service.NewPluginService,
	api.NewEsController,
	api.NewEsLinkController,

	es.NewEsClientService,
	es_link_service.NewEsLinkService,
	es_service.NewEsService,
	api.NewGmOperaterController,
	gm_operater_log.NewGmOperaterLogService,
	api.NewManagerRoleController,
	gm_role.NewGmRoleService,
	api.NewManagerUserController,
	gm_user.NewGmUserService,
	jwt_svr.NewJwt,
	middleware.NewMiddleWareService,
	backgroundsvcs.ProvideBackgroundServiceRegistry,
	access_control.NewRbac,
	web_engine.NewWebEngine,
	web.NewWebServer,
	plugin_rpc.NewPluginRpcServer,
	NewServer,
)

func Initialize(args *config.CommandLineArgs) (*Server, error) {
	wire.Build(wireSet)
	return &Server{}, nil
}

func InitializeOrm(args *config.CommandLineArgs) (*orm.Gorm, error) {
	wire.Build(wireSet)
	return &orm.Gorm{}, nil
}

func InitializeEvApiDao(args *config.CommandLineArgs) (*dao.EvBackDao, error) {
	wire.Build(wireSet)
	return &dao.EvBackDao{}, nil
}

func InitializeGmRoleEslinkCfgV2Dao(args *config.CommandLineArgs) (*dao.GmRoleEslinkCfgV2Dao, error) {
	wire.Build(wireSet)
	return &dao.GmRoleEslinkCfgV2Dao{}, nil
}

func InitializeProvideInstaller(args *config.CommandLineArgs) (*plugin_install_service.PluginInstaller, error) {
	wire.Build(wireSet)
	return &plugin_install_service.PluginInstaller{}, nil
}
