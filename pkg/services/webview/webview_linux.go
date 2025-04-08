//go:build linux
// +build linux

package webview

import (
	"context"
)

func (this *WebView) Run(ctx context.Context) (err error) {
	<-ctx.Done()
	return nil
}
