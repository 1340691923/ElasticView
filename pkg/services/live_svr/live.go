package live_svr

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/consts"
	"github.com/1340691923/eve-plugin-sdk-go/live"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/centrifugal/centrifuge"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Live struct {
	Handler        http.Handler
	log            *logger.AppLogger
	cfg            *config.Config
	node           *centrifuge.Node
	jwtSvr         *jwt_svr.Jwt
	pluginRegistry manager.Service
}

func (this *Live) Node() *centrifuge.Node {
	return this.node
}

func (this *Live) handleLog(msg centrifuge.LogEntry) {
	arr := make([]interface{}, 0)
	for k, v := range msg.Fields {
		switch v {
		case nil:
			v = "<nil>"
		case "":
			v = "<empty>"
		}
		arr = append(arr, k, v)
	}

	switch msg.Level {
	case centrifuge.LogLevelDebug:
		this.log.Sugar().Debugf(msg.Message, arr...)
	case centrifuge.LogLevelError:
		this.log.Sugar().Errorf(msg.Message, arr...)
	case centrifuge.LogLevelInfo:
		this.log.Sugar().Infof(msg.Message, arr...)
	default:
		this.log.Sugar().Debugf(msg.Message, arr...)
	}
}

const clientConcurrency = 12

func NewLive(log *logger.AppLogger, cfg *config.Config, jwtSvr *jwt_svr.Jwt, pluginRegistry manager.Service) *Live {
	live := &Live{
		log:            log,
		cfg:            cfg,
		jwtSvr:         jwtSvr,
		pluginRegistry: pluginRegistry,
	}
	// 初始化 Centrifuge Node，添加必要的配置
	centrifugeCfg := centrifuge.Config{
		LogLevel:           centrifuge.LogLevelError,
		LogHandler:         live.handleLog,
		ClientQueueMaxSize: 4194304, // 4MB

		HistoryMetaTTL: 7 * 24 * time.Hour,
	}

	node, err := centrifuge.New(centrifugeCfg)
	if err != nil {
		panic(err)
	}

	// 添加连接令牌处理
	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {

		c, err := jwtSvr.ParseToken(e.Token)

		if err != nil {
			log.Sugar().Errorf("ParseToken err token:%v err:%v", e.Token, err)
			return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
		}

		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{
				UserID:   cast.ToString(c.UserID),
				ExpireAt: c.ExpiresAt.Unix(),
			},
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		numConnections := node.Hub().NumClients()
		if cfg.GetLiveMaxConnections() >= 0 && numConnections > cfg.GetLiveMaxConnections() {
			log.Sugar().Warnf(
				"Max number of Live connections reached, increase max_connections in [live] configuration section",
				client.UserID(), client.ID(), cfg.GetLiveMaxConnections(),
			)
			client.Disconnect(centrifuge.DisconnectConnectionLimit)
			return
		}

		var semaphore chan struct{}
		if clientConcurrency > 1 {
			semaphore = make(chan struct{}, clientConcurrency)
		}
		log.Sugar().Debugf("Client connected", "user", client.UserID(), "client", client.ID())
		connectedAt := time.Now()

		// Called when client subscribes to the channel.
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			err := runConcurrentlyIfNeeded(client.Context(), semaphore, func() {
				cb(live.handleOnSubscribe(context.Background(), client, e))
			})
			if err != nil {
				cb(centrifuge.SubscribeReply{}, err)
			}
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			err := runConcurrentlyIfNeeded(client.Context(), semaphore, func() {
				cb(live.handleOnPublish(context.Background(), client, e))
			})
			if err != nil {
				cb(centrifuge.PublishReply{}, err)
			}
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			reason := e.Disconnect.Reason
			if e.Disconnect.Code == 3001 { // Shutdown
				return
			}

			log.Sugar().Debugf("Client disconnected", "user", client.UserID(), "client", client.ID(), "reason", reason, "elapsed", time.Since(connectedAt).String())
		})
	})

	// 启动 node
	if err := node.Run(); err != nil {
		panic(err)
	}

	// WebSocket 处理，添加 CORS 支持
	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 开发环境允许所有来源
		},
	})

	live.Handler = wsHandler
	live.node = node

	/*go func() {
		i := 0
		for {
			time.Sleep(3 * time.Second)

			err := live.LiveBroadcastEvMsg(false, &CommonMsg{
				Id:          i,
				PublishTime: time.Now().Format(util.TimeFormat),
				Title:       "test",
				Msg:         "test",
				BtnDesc:     "test",
				BtnJumpUrl:  "test",
				Type:        "success",
			})
			log.Sugar().Errorf("LiveBroadcastEvMsg err", err)
		}
	}()*/

	return live
}

func (this *Live) handleOnSubscribe(ctx context.Context, client *centrifuge.Client, e centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	// 记录客户端订阅请求的日志
	this.log.Sugar().Debugf("Client wants to subscribe", "user", client.UserID(), "client", client.ID(), "channel", e.Channel)
	// 直接返回成功的订阅响应，无需验证
	return centrifuge.SubscribeReply{
		Options: centrifuge.SubscribeOptions{
			EmitPresence:   true, // 启用 presence（用户在线状态）
			EmitJoinLeave:  true, // 启用用户加入和离开事件
			PushJoinLeave:  true, // 推送用户加入和离开的消息
			EnableRecovery: true, // 启用消息重发
		},
	}, nil
}

func (this *Live) clientSend(client *centrifuge.Client, channel string, data interface{}) (err error) {

	b, err := json.Marshal(vo.LivePrivateData{
		Channel: channel,
		Data:    data,
	})
	if err != nil {
		return
	}

	err = client.Send(b)
	return
}

func (this *Live) handleOnPublish(ctx context.Context, client *centrifuge.Client, e centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
	// 记录客户端发布请求的日志
	this.log.Sugar().Debugf("Client wants to publish", "user", client.UserID(), "client", client.ID(), "channel", e.Channel)

	var respMap map[string]interface{}
	if e.Channel == consts.EvAllMsgChannel ||
		strings.Contains(e.Channel, consts.EvRoleMsgChannel) ||
		strings.Contains(e.Channel, consts.EvUserMsgChannel) {
		this.log.Info("common msg", zap.String("e.Channel", e.Channel))
	} else {
		pluginId, channel, err := this.ParseChannel(e.Channel)

		if err != nil {
			this.log.Sugar().Errorf("ParseChannel err", e.Channel, err)
			return centrifuge.PublishReply{}, centrifuge.ErrorInternal
		}

		plugin, ok := this.pluginRegistry.Plugin(ctx, pluginId)

		if !ok {
			this.log.Sugar().Errorf("该插件不存在 err", pluginId, e.Channel)
			return centrifuge.PublishReply{}, centrifuge.ErrorInternal
		}

		resp, err := plugin.Pub2Channel(ctx, &backend.Pub2ChannelRequest{
			Channel:       channel,
			Data:          e.Data,
			PluginContext: backend.PluginContext{},
		})

		if err != nil {
			this.log.Sugar().Errorf("Pub2Channel err", e.Channel, string(e.Data), err)
			return centrifuge.PublishReply{}, centrifuge.ErrorInternal
		}

		if resp.Status != backend.PubStatusOk {
			this.log.Sugar().Errorf("Pub2Channel err", e.Channel, string(e.Data), resp.Status)
			return centrifuge.PublishReply{}, centrifuge.ErrorInternal
		}

		err = json.Unmarshal(resp.JsonDetails, &respMap)
		if err != nil {
			this.log.Sugar().Errorf("Pub2Channel Unmarshal res err", e.Channel, string(resp.JsonDetails), err)
			return centrifuge.PublishReply{}, centrifuge.ErrorInternal
		}

		err = this.clientSend(client, e.Channel, respMap)
		if err != nil {
			this.log.Sugar().Errorf("Failed to send message to client", "error", err)
			return centrifuge.PublishReply{}, centrifuge.ErrorInternal
		}

	}

	// 由于消息已经通过 `client.Send()` 发送，我们返回一个空的 `PublishReply` 来避免广播
	return centrifuge.PublishReply{
		Options: centrifuge.PublishOptions{
			HistorySize: 1000,
			HistoryTTL:  7 * 24 * time.Hour,
		},
		Result: &centrifuge.PublishResult{},
	}, nil
}

func runConcurrentlyIfNeeded(ctx context.Context, semaphore chan struct{}, fn func()) error {
	if cap(semaphore) > 1 {
		select {
		case semaphore <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}
		go func() {
			defer func() { <-semaphore }()
			fn()
		}()
	} else {
		// No need in separate goroutines.
		fn()
	}
	return nil
}

func (this *Live) ParseChannel(channel string) (pluginId string, parseChannel string, err error) {

	arr := strings.Split(channel, consts.ChannelSplit)

	if len(arr) != 2 {
		err = errors.New("channel parse err:" + channel)
		return
	}

	pluginId = arr[0]
	parseChannel = arr[1]
	return
}

func (this *Live) LiveBroadcast(channel string, data interface{}) (err error) {
	numSubscribers := this.Node().Hub().NumSubscribers(channel)

	if numSubscribers <= 0 {
		return live.NoSubscriberErr
	}

	b := []byte{}
	if data != nil {
		b, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	_, err = this.Node().Publish(channel, b)

	if err != nil {
		return err
	}

	return nil
}
