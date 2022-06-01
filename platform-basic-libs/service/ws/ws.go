package ws

import (
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
)

var ConnUUidMap sync.Map

type ManagerConnMap struct {
	Conns *sync.Map
}

func NewManagerConnMap() *ManagerConnMap {
	return &ManagerConnMap{
		Conns: new(sync.Map),
	}
}

func (this *ManagerConnMap) AddManagerConn(manager string, conn *websocket.Conn) {
	this.Conns.Store(manager, conn)
}

func (this *ManagerConnMap) DeleteConn(manager string) {
	this.Conns.Delete(manager)
}

func removeConns(deleteConn *websocket.Conn) {
	ConnUUidMap.Range(func(distinctId, managerMap interface{}) bool {
		managerMap.(*ManagerConnMap).Conns.Range(func(key, conn interface{}) bool {
			if conn == deleteConn {
				managerMap.(*ManagerConnMap).DeleteConn(key.(string))
				ConnUUidMap.Store(distinctId, managerMap)
			}
			return true
		})
		return true
	})
	deleteConn.Close()
}

//长链接
func Ws(c *websocket.Conn) {

	type ReqData struct {
		UUid     string `json:"uuid"`
		Token    string `json:"token"`
		Ping     string `json:"ping"`
		SendType string `json:"send_type"`
	}

	for {
		var reqData ReqData
		reqData.UUid = ""
		err := c.ReadJSON(&reqData)

		if err != nil {
			logs.Logger.Error("ws ReadMessage", zap.Error(err))
			removeConns(c)
			break
		}

		if reqData.UUid == "" {
			logs.Logger.Sugar().Errorf("reqData.UUid = %v", reqData.UUid)
			continue
		}
		cliams, err := jwt.ParseToken(reqData.Token)
		if err != nil {
			logs.Logger.Error("jwt.ParseToken", zap.Error(err))
			break
		}

		if cliams.ID == 0 {
			logs.Logger.Error("cliams.UserID = 0")
			break
		}

		if data, found := ConnUUidMap.Load(reqData.UUid); !found {
			conn := NewManagerConnMap()
			conn.AddManagerConn(reqData.UUid, c)
			ConnUUidMap.Store(reqData.UUid, conn)
		} else {
			data.(*ManagerConnMap).AddManagerConn(reqData.UUid, c)
			ConnUUidMap.Store(reqData.UUid, data)
		}

		err = c.WriteJSON(map[string]interface{}{"code": 0})
		if err != nil {
			if err == io.EOF {
				logs.Logger.Error("客户端已经断开WsSocket!", zap.Error(err))
			} else if err.(*net.OpError).Err.Error() == "use of closed network connection" {
				logs.Logger.Error("服务端已经断开WsSocket!", zap.Error(err))
			}
			break
		}
	}
}

func SendWs(distinctId string, typ string, data interface{}) {
	var err error
	managerMap, ok := ConnUUidMap.Load(distinctId)

	if ok {
		managerMap.(*ManagerConnMap).Conns.Range(func(key, value interface{}) bool {
			if err := value.(*websocket.Conn).WriteJSON(map[string]interface{}{
				"code": 1,
				"data": data,
				"typ":  typ,
			}); err != nil {
				if err == io.EOF {
					logs.Logger.Error("客户端已经断开WsSocket!", zap.Error(err))
				} else if err.Error() == "use of closed network connection" {
					logs.Logger.Error("服务端已经断开WsSocket!", zap.Error(err))
				} else {
					logs.Logger.Error("socket err!", zap.Error(err))
				}
				managerMap.(*ManagerConnMap).DeleteConn(key.(string))
				ConnUUidMap.Store(distinctId, managerMap)
			}
			return true
		})
	}

	if err != nil {
		return
	}

	return
}
