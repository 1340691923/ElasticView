package api

import (
	"github.com/1340691923/ElasticView/platform-basic-libs/service/ws"
	"github.com/gofiber/websocket/v2"
)

//长链接
func Ws(c *websocket.Conn) {
	ws.Ws(c)
}
