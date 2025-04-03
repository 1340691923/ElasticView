package manager

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
	"sync"
)

type PluginManager struct {
	store map[string]*plugin.Plugin
	mu    sync.RWMutex
}

func NewPluginManager() (pluginManager *PluginManager) {
	return &PluginManager{store: map[string]*plugin.Plugin{}}
}

func (this *PluginManager) Plugin(_ context.Context, pluginID string) (*plugin.Plugin, bool) {
	return this.plugin(pluginID)
}

func (this *PluginManager) Plugins(_ context.Context) []*plugin.Plugin {
	this.mu.RLock()
	defer this.mu.RUnlock()

	res := make([]*plugin.Plugin, 0, len(this.store))
	for _, p := range this.store {
		res = append(res, p)
	}

	return res
}

func (this *PluginManager) AddPlugin(_ context.Context, p *plugin.Plugin) error {

	if this.isRegistered(p.ID) {
		return fmt.Errorf("plugin %s is already registered", p.ID)
	}

	this.mu.Lock()
	defer this.mu.Unlock()
	this.store[p.ID] = p

	return nil
}

func (this *PluginManager) Reload(ctx context.Context, newPlugin *plugin.Plugin) error {

	this.mu.Lock()
	defer this.mu.Unlock()

	oldPlugin := this.store[newPlugin.ID]
	err := oldPlugin.Decommission()

	if err != nil {
		return err
	}

	err = newPlugin.Start(ctx)

	if err != nil {
		oldPlugin.Start(ctx)
		return err
	}

	this.store[newPlugin.ID] = newPlugin

	err = oldPlugin.Stop(ctx)
	if err != nil {
		return err
	}

	newPlugin.DisDecommission()

	return nil
}

func (this *PluginManager) Remove(_ context.Context, pluginID string) error {
	_, ok := this.plugin(pluginID)
	if !ok {
		return fmt.Errorf("plugin %s is not registered", pluginID)
	}

	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.store, pluginID)

	return nil
}

func (this *PluginManager) plugin(pluginID string) (*plugin.Plugin, bool) {
	this.mu.RLock()
	defer this.mu.RUnlock()
	p, exists := this.store[pluginID]

	if !exists {
		return nil, false
	}

	return p, true
}

func (this *PluginManager) isRegistered(pluginID string) bool {
	p, exists := this.plugin(pluginID)

	// This may have matched based on an alias
	return exists && p.ID == pluginID
}
