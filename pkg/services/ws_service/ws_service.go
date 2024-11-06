package ws_service

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"runtime"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
	"reflect"
	"sync"
)

type Ctx struct {
	ReqData  interface{}
	Context  context.Context
	WsClient *websocket.Conn
	UserId   int `json:"user_id"`
	RoleId   int `json:"role_id"`
}

type EventCallback func(ctx *Ctx)
type ApiHandle func(ctx *Ctx) interface{}
type EventCallbackWithErr func(ctx *Ctx) (err error)

type WsService struct {
	log                      *logger.AppLogger
	cfg                      *config.Config
	orm                      *sqlstore.SqlStore
	jwtSvr                   *jwt_svr.Jwt
	protoMap                 map[int]reflect.Type
	svrRouter                map[string]ApiHandle
	onConnectCallback        EventCallback
	onSessionAccepted        EventCallback
	onSessionClosed          EventCallback
	onDestory                EventCallbackWithErr
	onMounted                EventCallbackWithErr
	onMountedExcludeProtoMap map[string]struct{}
	onDestoryExcludeProtoMap map[string]struct{}
	user2ConnMap             *sync.Map

	heartController *HeartController
}

func (this *WsService) init() {
	//添加请求协议
	this.addReqProto(C2S_PING, &dto.C2S_PING{})

	//添加控制器
	this.addStrcuts(
		this.heartController,
	)

	//钩子
	this.OnSessionAccepted(func(ctx *Ctx) {
		this.log.Sugar().Infof("新建连接", ctx.UserId)
	})

	this.OnSessionClosed(func(ctx *Ctx) {
		this.log.Sugar().Infof("连接断开", ctx.UserId)
	})

	this.OnDestoryCallback(func(ev *Ctx) error {
		this.log.Sugar().Infof("控制器销毁", ev.UserId)
		return nil
	})

	this.OnMountedCallback(func(ctx *Ctx) error {
		this.log.Sugar().Infof("控制器初始化", ctx.UserId)
		return nil
	})

}

func (this *WsService) addReqProto(code int, proto interface{}) {
	this.protoMap[code] = reflect.TypeOf(proto).Elem()
}

func NewWsService(
	log *logger.AppLogger,
	cfg *config.Config,
	orm *sqlstore.SqlStore,
	jwtSvr *jwt_svr.Jwt,
	heartController *HeartController,
) *WsService {
	wsSvr := &WsService{
		log:                      log,
		cfg:                      cfg,
		orm:                      orm,
		jwtSvr:                   jwtSvr,
		protoMap:                 map[int]reflect.Type{},
		svrRouter:                map[string]ApiHandle{},
		onMountedExcludeProtoMap: map[string]struct{}{},
		onDestoryExcludeProtoMap: map[string]struct{}{},
		heartController:          heartController,
		user2ConnMap:             new(sync.Map),
	}
	wsSvr.init()
	return wsSvr
}

func (this *WsService) InitConnect(conn *websocket.Conn, ginCtx *gin.Context, userId, roleId int) {

	ctx := context.WithValue(context.Background(), "request", ginCtx.Request)

	wsCtx := &Ctx{
		Context:  ctx,
		WsClient: conn,
		UserId:   userId,
		RoleId:   roleId,
	}

	this.SessionAccepted(wsCtx)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			this.SessionClosed(&Ctx{
				Context:  ctx,
				WsClient: conn,
				UserId:   userId,
				RoleId:   roleId,
			})
			log.Println("server read err:", err)
			break
		}
		reqData, err := this.getReqData(message)

		if err != nil {
			log.Println(" server getData err:", err)
			break
		}
		wsCtx.ReqData = reqData
		this.RunServer(wsCtx)
	}
}

func (this *WsService) getReqData(reqByte []byte) (reqData interface{}, err error) {
	code := gjson.GetBytes(reqByte, "code").Int()
	reqDataString := gjson.GetBytes(reqByte, "data").String()
	_, ok := this.protoMap[int(code)]

	if !ok {
		errString := fmt.Sprintf("server err code:%d", code)
		err = errors.New(errString)
		return
	}
	reqData = reflect.New(this.protoMap[int(code)]).Interface()

	err = json.Unmarshal([]byte(reqDataString), &reqData)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (this *WsService) Mounted(ev *Ctx) error {
	var err error
	if this.onMounted != nil {
		key := reflect.TypeOf(ev.ReqData).String()
		if _, ok := this.onMountedExcludeProtoMap[key]; !ok {
			err = this.onMounted(ev)
		}
	}
	return err
}

func (this *WsService) OnMountedCallback(fn func(ev *Ctx) error, excludeProtoMsgs ...interface{}) {
	this.onMounted = fn
	for _, protoMsg := range excludeProtoMsgs {
		key := reflect.TypeOf(protoMsg).String()
		this.onMountedExcludeProtoMap[key] = struct{}{}
	}
}

func (this *WsService) OnConnectCallback(fn EventCallback) {
	this.onConnectCallback = fn
}

func (this *WsService) ConnectCallback(ev *Ctx) {
	if this.onConnectCallback != nil {
		this.onConnectCallback(ev)
	}
}

func (this *WsService) OnSessionClosed(fn EventCallback) {
	this.onSessionClosed = fn
}

func (this *WsService) OnDestoryCallback(fn func(ev *Ctx) error, excludeProtoMsgs ...interface{}) {
	this.onDestory = fn
	for _, protoMsg := range excludeProtoMsgs {
		key := reflect.TypeOf(protoMsg).String()
		this.onDestoryExcludeProtoMap[key] = struct{}{}
	}

}

func (this *WsService) Destory(ev *Ctx) error {
	var err error
	if this.onDestory != nil {
		key := reflect.TypeOf(ev.ReqData).String()
		if _, ok := this.onDestoryExcludeProtoMap[key]; !ok {
			err = this.onDestory(ev)
		}
	}
	return err
}

func (this *WsService) SessionClosed(ev *Ctx) {
	if this.onSessionClosed != nil {
		this.onSessionClosed(ev)
	}
}

func (this *WsService) OnSessionAccepted(fn EventCallback) {
	this.onSessionAccepted = fn
}

func (this *WsService) SessionAccepted(ev *Ctx) {
	if this.onSessionAccepted != nil {
		this.onSessionAccepted(ev)
	}
}

func (this *WsService) AddRoute(protoMsg interface{}, handler ApiHandle) {
	key := reflect.TypeOf(protoMsg).String()
	this.svrRouter[key] = handler
}

func (this *WsService) GetHandle(protoName string) (ApiHandle, bool) {
	application, b := this.svrRouter[protoName]
	return application, b
}

func (this *WsService) addStrcuts(data ...interface{}) {
	for _, dataV := range data {
		dataValue := reflect.ValueOf(dataV)
		for i := 0; i < dataValue.NumMethod(); i++ {
			handlerType := reflect.TypeOf(dataValue.Method(i).Interface())
			key := "*" + handlerType.In(0).Elem().String()
			tmp := i
			this.svrRouter[key] = func(ctx *Ctx) interface{} {
				values := []reflect.Value{reflect.ValueOf(ctx.ReqData), reflect.ValueOf(ctx)}
				resV := dataValue.Method(tmp).Call(values)
				if len(resV) != 1 {
					return nil
				}
				return resV[0].Interface()
			}
		}
	}

}

func (this *WsService) RunServer(ev *Ctx) {

	defer func() {
		if r := recover(); r != nil {
			//打印调用栈信息
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			this.log.Sugar().Errorf("panic stack info %s", stackInfo)
			this.log.Sugar().Errorf("--->Server Error:", r)
		}
	}()

	msgType := reflect.TypeOf(ev.ReqData).String()
	handle, foundHandle := this.GetHandle(msgType)
	if !foundHandle {
		this.log.Error("WS 没有找到该路由", zap.String("msgType", msgType))
		return
	}
	this.log.Debug("WS 收到请求", zap.Reflect(msgType, ev.ReqData))
	err := this.Mounted(ev)
	if err != nil {
		this.log.Error("WS 处理消息失败", zap.Reflect("msgType", ev.ReqData), zap.Error(err))
		return
	}
	res := handle(ev)

	if res != nil {
		this.log.Error("WS 处理请求并返回", zap.Reflect(msgType, ev.ReqData), zap.Reflect(reflect.TypeOf(res).String(), res))
		ev.WsClient.WriteJSON(res)
	} else {
		this.log.Error("WS 处理消息失败 res = nil")
	}

	this.Destory(ev)

}
