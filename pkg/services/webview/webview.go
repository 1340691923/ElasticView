package webview

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"

	"go.uber.org/zap"
)

type WebView struct {
	log *zap.Logger
	cfg *config.Config
}

func ProvideWebView(log *logger.AppLogger, cfg *config.Config) (*WebView, error) {
	log = log.Named("webview")
	return &WebView{
		log: logger.ZapLog2AppLog(log),
		cfg: cfg,
	}, nil
}
