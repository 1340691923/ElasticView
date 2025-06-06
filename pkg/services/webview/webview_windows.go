//go:build windows
// +build windows

package webview

import (
	"context"
	"fmt"
	"github.com/inkeliz/gowebview"
	"github.com/pkg/errors"
	"runtime"
)

func (this *WebView) Run(ctx context.Context) (err error) {
	var w gowebview.WebView
	if !this.cfg.DeBug && runtime.GOOS == "windows" {
		openAddr := fmt.Sprintf("http://localhost:%d/#/", this.cfg.Port)
		w, err = gowebview.New(&gowebview.Config{
			Debug: this.cfg.DeBug,
			URL:   openAddr,
			WindowConfig: &gowebview.WindowConfig{
				Title: "ElasticView",
				Size:  &gowebview.Point{X: 1280, Y: 720},
			}})
		if err != nil {
			return errors.WithStack(err)
		}
		defer w.Destroy()

		go func() {
			w.Run()
		}()
	}
	<-ctx.Done()
	return nil
}
