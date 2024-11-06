package process

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
)

type Manager interface {
	Start(ctx context.Context, p *plugin.Plugin) error
	Stop(ctx context.Context, p *plugin.Plugin) error
}
