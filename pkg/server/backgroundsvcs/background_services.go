package backgroundsvcs

import (
	"github.com/1340691923/ElasticView/pkg/registry"
	"github.com/1340691923/ElasticView/pkg/web"
)

func ProvideBackgroundServiceRegistry(
	httpServer *web.WebServer,
) *BackgroundServiceRegistry {
	return NewBackgroundServiceRegistry(
		httpServer,
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
