//go:build windows
// +build windows

package server

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/resources/views"
	"github.com/getlantern/systray"
	"github.com/inkeliz/gowebview"
	"log"
	"runtime"
)

func (this *Server) runSystray() {
	go systray.Run(this.onReady, this.onExit)
}

func (this *Server) onReady() {
	iconData, _ := views.GetFavicon()

	systray.SetIcon(iconData)

	openItem := systray.AddMenuItem("打开", "打开浏览器访问")
	exitItem := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			select {
			case <-openItem.ClickedCh:
				if runtime.GOOS != "windows" {
					openAddr := fmt.Sprintf("http://localhost:%d", this.cfg.Port)
					util.OpenWinBrowser(openAddr)
				} else {
					openAddr := fmt.Sprintf("http://localhost:%d/#/", this.cfg.Port)
					w, err := gowebview.New(&gowebview.Config{
						Debug: this.cfg.DeBug,
						URL:   openAddr,
						WindowConfig: &gowebview.WindowConfig{
							Title: "ElasticView",
							Size:  &gowebview.Point{X: 1280, Y: 720},
						}})
					if err == nil {
						go func() {
							w.Run()
						}()
					} else {
						log.Println("open webview err", err)
					}
				}
			case <-exitItem.ClickedCh:
				this.Shutdown(this.context)
				return
			}
		}
	}()
}

func (this *Server) onExit() {
	this.Shutdown(this.context)
}
