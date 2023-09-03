//go:build wireinject
// +build wireinject

package server

import (
	"github.com/1340691923/ElasticView/pkg/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	log2 "github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/middleware"
	"github.com/1340691923/ElasticView/pkg/registry"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/server/backgroundsvcs"
	"github.com/1340691923/ElasticView/pkg/services/alias_service"
	"github.com/1340691923/ElasticView/pkg/services/cat_service"
	"github.com/1340691923/ElasticView/pkg/services/cluser_settings_service"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_backup"
	"github.com/1340691923/ElasticView/pkg/services/es_doc_service"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/1340691923/ElasticView/pkg/services/es_service"
	"github.com/1340691923/ElasticView/pkg/services/es_task_service"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/1340691923/ElasticView/pkg/services/gm_role"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/services/index_service"
	"github.com/1340691923/ElasticView/pkg/services/navicat_service"
	"github.com/1340691923/ElasticView/pkg/web"
	"github.com/google/wire"
)

var wireSet = wire.NewSet(
	wire.Bind(new(registry.BackgroundServiceRegistry), new(*backgroundsvcs.BackgroundServiceRegistry)),
	config.InitConfig,
	log2.InitLog,
	sqlstore.NewSqlStore,

	request.NewRequest,
	response.NewResponse,
	api.NewBaseController,

	api.NewDslHistoryController,
	api.NewEsBackUpController,
	api.NewEsController,
	api.NewEsIndexController,
	api.NewEsLinkController,
	api.NewEsMappingController,
	api.NewEsTaskController,
	api.NewGuidController,
	es.NewEsClientService,
	es.NewEsCache,
	alias_service.NewAliasService,
	cat_service.NewCatService,
	cluser_settings_service.NewClusterSettingsService,
	es_backup.NewEsBackUpService,
	api.NewEsDocController,
	es_task_service.NewEsTaskService,
	index_service.NewIndexService,
	navicat_service.NewNavicatService,
	es_doc_service.NewEsDocService,
	api.NewEsCrudController,
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
	NewServer,
)

func Initialize(args *config.CommandLineArgs) (*Server, error) {
	wire.Build(wireSet)
	return &Server{}, nil
}
