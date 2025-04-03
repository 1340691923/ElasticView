package backendplugin

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/hashicorp/go-hclog"
	process2 "github.com/shirou/gopsutil/v3/process"
)

type Plugin interface {
	PluginID() string
	Logger() hclog.Logger
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Exited() bool
	Decommission() error
	DisDecommission() error
	IsDecommissioned() bool
	backend.CheckHealthHandler
	backend.CallResourceHandler
	backend.LiveHandler
	GetPid() int
	GetProcessUtil() (*process2.Process, error)
}
