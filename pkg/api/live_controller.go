package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/services/live_svr"
	"github.com/1340691923/ElasticView/pkg/util"
	"net/http"
)

type LiveController struct {
	live *live_svr.Live
	log  *logger.AppLogger
}

func NewLiveController(live *live_svr.Live, log *logger.AppLogger) *LiveController {
	return &LiveController{live: live, log: log}
}

func (this *LiveController) HttpHandle() http.Handler {
	return util.Cors(this.live.Handler)
}
