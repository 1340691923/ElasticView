package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/services/live_svr"
	"net/http"
)

type WsController struct {
	live *live_svr.Live
	log  *logger.AppLogger
}

func NewWsController(live *live_svr.Live, log *logger.AppLogger) *WsController {
	return &WsController{live: live, log: log}
}

func (this *WsController) HttpHandle() http.Handler {
	return this.Cors(this.live.Handler)
}

func (this *WsController) Cors(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET, POST, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Accept-Encoding,Accept-Language,Access-Control-Request-Headers,Access-Control-Request-Method,Connection,Referer,Sec-Fetch-Dest,User-Agent, Origin,Authorization,Content-Type,X-Token,x-token,X-Version,Current-Language")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})

}
