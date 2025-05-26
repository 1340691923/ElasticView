package backgroundsvcs

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugin_rpc"
	"github.com/1340691923/ElasticView/pkg/infrastructure/pluginstore"
	"github.com/1340691923/ElasticView/pkg/registry"
	"github.com/1340691923/ElasticView/pkg/services/eve_service"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/1340691923/ElasticView/pkg/services/print_logo"
	"github.com/1340691923/ElasticView/pkg/services/updatechecker"
	"github.com/1340691923/ElasticView/pkg/services/web"
	"github.com/1340691923/ElasticView/pkg/services/webview"
)

func ProvideBackgroundServiceRegistry(
	httpServer *web.WebServer,
	evUpdate *updatechecker.EvUpdate,
	pluginsUpdate *updatechecker.PluginsService,
	pluginStoreService *pluginstore.PluginStoreService,
	pluginRpcServer *plugin_rpc.PluginRpcServer,
	eveService *eve_service.EvEService,
	gmOperaterLogService *gm_operater_log.GmOperaterLogService,
	webView *webview.WebView,
	printLogo *print_logo.PrintLogo,
) *BackgroundServiceRegistry {
	return NewBackgroundServiceRegistry(
		httpServer,
		pluginStoreService,
		pluginRpcServer,
		evUpdate,
		pluginsUpdate,
		eveService,
		gmOperaterLogService,
		webView,
		printLogo,
	)
}

// BackgroundServiceRegistry provides background services.
type BackgroundServiceRegistry struct {
	Services []registry.BackgroundService
}

func NewBackgroundServiceRegistry(services ...registry.BackgroundService) *BackgroundServiceRegistry {
	return &BackgroundServiceRegistry{services}
}

func (r *BackgroundServiceRegistry) GetServices() []registry.BackgroundService {
	return r.Services
}
