package backendplugin

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/hashicorp/go-hclog"
)

type Plugin interface {
	PluginID() string
	Logger() hclog.Logger
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Exited() bool
	Decommission() error
	IsDecommissioned() bool
	backend.CheckHealthHandler
	backend.CallResourceHandler
}
