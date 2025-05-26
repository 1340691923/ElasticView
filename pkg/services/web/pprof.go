package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
)

func (this *WebServer) runpprof() {
	pprofGroup := this.engine.GetGinEngine().Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(http.HandlerFunc(pprof.Handler("allocs").ServeHTTP)))
		pprofGroup.GET("/block", gin.WrapH(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
		pprofGroup.GET("/goroutine", gin.WrapH(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
		pprofGroup.GET("/heap", gin.WrapH(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
		pprofGroup.GET("/mutex", gin.WrapH(http.HandlerFunc(pprof.Handler("mutex").ServeHTTP)))
		pprofGroup.GET("/threadcreate", gin.WrapH(http.HandlerFunc(pprof.Handler("threadcreate").ServeHTTP)))
	}
}
