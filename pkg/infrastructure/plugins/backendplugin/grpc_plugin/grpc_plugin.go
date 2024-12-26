package grpcplugin

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/process"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
	process2 "github.com/shirou/gopsutil/v3/process"
	"os"
	"os/exec"

	"sync"
)

type pluginClient interface {
	backend.CheckHealthHandler
	backend.CallResourceHandler
}

type GrpcPlugin struct {
	descriptor            PluginDescriptor
	clientFactory         func() *plugin.Client
	client                *plugin.Client
	pluginClient          pluginClient
	logger                hclog.Logger
	mutex                 sync.RWMutex
	closeLogWriteCallback func() error
	decommissioned        bool
	pid                   int
	cmd                   *exec.Cmd
}

func (p *GrpcPlugin) GetPid() int {

	if p.cmd != nil {
		return p.cmd.Process.Pid
	}

	if p.descriptor.isDebug {
		return p.pid
	}

	return 0
}

func (p *GrpcPlugin) GetProcessUtil() (*process2.Process, error) {
	return process2.NewProcess(int32(p.GetPid()))
}

func newPlugin(descriptor PluginDescriptor, env []string, log hclog.Logger, closeLogWriteCallback func() error) *GrpcPlugin {
	p := &GrpcPlugin{
		descriptor: descriptor,
		logger:     log,

		closeLogWriteCallback: closeLogWriteCallback,
	}

	if !descriptor.isDebug {
		os.Chmod(descriptor.executablePath, 0755)
		cmd := exec.Command(descriptor.executablePath, descriptor.executableArgs...)
		cmd.Env = env
		p.cmd = cmd
	} else {
		p.pid = descriptor.pid
	}

	p.clientFactory = func() *plugin.Client {
		var clientConfig *plugin.ClientConfig
		if descriptor.isDebug {
			clientConfig = newTestClientConfig(
				descriptor.pluginAddr,
				descriptor.pid,
				log,
				descriptor.versionedPlugins,
			)
		} else {
			clientConfig = newProductClientConfig(
				p.cmd,
				log,
				descriptor.versionedPlugins,
			)
		}
		return plugin.NewClient(clientConfig)
	}

	return p
}

func (p *GrpcPlugin) PluginID() string {
	return p.descriptor.pluginID
}

func (p *GrpcPlugin) IsDebug() bool {
	return p.descriptor.isDebug
}

func (p *GrpcPlugin) Logger() hclog.Logger {
	return p.logger
}

func (p *GrpcPlugin) Start(_ context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.client = p.clientFactory()

	p.client.Kill()

	rpcClient, err := p.client.Client()
	if err != nil {
		return err
	}

	if p.client.NegotiatedVersion() < 2 {
		return errors.New("plugin protocol version not supported")
	}
	p.pluginClient, err = newClientV2(rpcClient)
	if err != nil {
		return err
	}

	if p.pluginClient == nil {
		return errors.New("no compatible plugin implementation found")
	}

	elevated, err := process.IsRunningWithElevatedPrivileges()
	if err != nil {
		p.logger.Error("Error checking plugin process execution privilege", "error", err)
	}
	if elevated {
		p.logger.Warn("Plugin process is running with elevated privileges. This is not recommended")
	}

	return nil
}

func (p *GrpcPlugin) Stop(_ context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.closeLogWriteCallback()
	if p.client != nil {
		p.client.Kill()
	}
	return nil
}

func (p *GrpcPlugin) Exited() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if p.client != nil {
		return p.client.Exited()
	}
	return true
}

func (p *GrpcPlugin) Decommission() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.decommissioned = true

	return nil
}

func (p *GrpcPlugin) IsDecommissioned() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.decommissioned
}

func (p *GrpcPlugin) getPluginClient() (pluginClient, bool) {
	p.mutex.RLock()
	if p.client == nil || p.client.Exited() || p.pluginClient == nil {
		p.mutex.RUnlock()
		return nil, false
	}
	pluginClient := p.pluginClient
	p.mutex.RUnlock()
	return pluginClient, true
}

func (p *GrpcPlugin) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	pluginClient, ok := p.getPluginClient()
	if !ok {
		return nil, errors.New("该插件没有实现CheckHealth接口")
	}
	return pluginClient.CheckHealth(ctx, req)
}

func (p *GrpcPlugin) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	pluginClient, ok := p.getPluginClient()
	if !ok {
		return errors.New("该插件没有实现CallResource接口")
	}
	return pluginClient.CallResource(ctx, req, sender)
}
