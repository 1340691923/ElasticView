package grpcplugin

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/process"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
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
}

func newPlugin(descriptor PluginDescriptor, env []string, log hclog.Logger, closeLogWriteCallback func() error) *GrpcPlugin {
	return &GrpcPlugin{
		descriptor: descriptor,
		logger:     log,
		clientFactory: func() *plugin.Client {
			var clientConfig *goplugin.ClientConfig
			if descriptor.isDebug {
				clientConfig = newTestClientConfig(
					descriptor.pluginAddr,
					descriptor.pid,
					log,
					descriptor.versionedPlugins,
				)
			} else {
				clientConfig = newProductClientConfig(
					descriptor.executablePath,
					descriptor.executableArgs,
					env,
					log,
					descriptor.versionedPlugins,
				)
			}

			return plugin.NewClient(clientConfig)
		},
		closeLogWriteCallback: closeLogWriteCallback,
	}
}

func (p *GrpcPlugin) PluginID() string {
	return p.descriptor.pluginID
}

func (p *GrpcPlugin) Logger() hclog.Logger {
	return p.logger
}

func (p *GrpcPlugin) Start(_ context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.client = p.clientFactory()

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
