package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type IndexController struct {
	cfg *config.Config
	log *logger.AppLogger
}

func NewIndexController(cfg *config.Config, log *logger.AppLogger) *IndexController {
	return &IndexController{cfg: cfg, log: log}
}

func (this *IndexController) IndexHtml(c *gin.Context) {

	appUrl := ""
	appBaseUrl := ""
	frontEndCfg := new(vo.FrontEndCfg)
	frontEndCfg.AppUrl = appUrl
	frontEndCfg.AppSubUrl = appBaseUrl
	frontEndCfg.Version = config.GetVersion()
	frontEndCfg.Lang = this.cfg.GetLang()
	frontEndCfg.WatermarkContent = this.cfg.WatermarkContent

	c.HTML(http.StatusOK, "index.html", frontEndCfg)
}

func (this *IndexController) CallBack(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	login_callback := gjson.Get(state, "login_callback").String()

	parmas := url.Values{}
	parmas.Set("code", code)
	parmas.Set("state", state)
	URL := fmt.Sprintf("%s/#/login?%s", login_callback, parmas.Encode())

	c.Redirect(http.StatusFound, URL)
}

func (this *IndexController) GetI18nCfg(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "succ",
		"data": config.TranslationCfg,
	})
}

func (this *IndexController) Health(c *gin.Context) {
	c.Writer.WriteString("I am OK!")
}
