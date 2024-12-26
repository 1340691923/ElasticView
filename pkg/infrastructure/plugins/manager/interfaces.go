package manager

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
)

type Service interface {
	//插件通过ID查找插件。
	Plugin(ctx context.Context, id string) (*plugin.Plugin, bool)
	//Plugins返回所有插件。
	Plugins(ctx context.Context) []*plugin.Plugin
	//Add将提供的插件添加到注册表中。
	Add(ctx context.Context, plugin *plugin.Plugin) error
	//Remove从注册表中删除请求的插件。
	Remove(ctx context.Context, id string) error
}
