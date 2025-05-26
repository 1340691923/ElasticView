package oauth

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
)

type OAuthRegistry struct {
	providers map[string]OAuthProvider
	log       *logger.AppLogger
}

func NewOAuthRegistry(cfg *config.Config, log *logger.AppLogger) *OAuthRegistry {
	registry := &OAuthRegistry{
		providers: make(map[string]OAuthProvider),
		log:       log,
	}
	
	if cfg.OAuth != nil {
		if cfg.OAuth.Wechat != nil && cfg.OAuth.Wechat.AppID != "" {
			registry.providers["wechat"] = NewWechatService(cfg.OAuth.Wechat, log)
		}
		
		if cfg.OAuth.Dingtalk != nil && cfg.OAuth.Dingtalk.AppID != "" {
			registry.providers["dingtalk"] = NewDingtalkService(cfg.OAuth.Dingtalk, log)
		}
		
		if cfg.OAuth.Feishu != nil && cfg.OAuth.Feishu.AppID != "" {
			registry.providers["feishu"] = NewFeishuService(cfg.OAuth.Feishu, log)
		}
	}
	
	return registry
}

func (r *OAuthRegistry) GetProvider(name string) (OAuthProvider, bool) {
	provider, exists := r.providers[name]
	return provider, exists
}

func (r *OAuthRegistry) GetProviders() map[string]OAuthProvider {
	return r.providers
}

func (r *OAuthRegistry) GetProviderNames() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
