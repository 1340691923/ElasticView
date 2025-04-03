package grpcplugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/grpcplugin"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type ClientV2 struct {
	grpcplugin.PluginInfoClient
	grpcplugin.ResourceClient
	grpcplugin.LiveClient
}

func newClientV2(rpcClient plugin.ClientProtocol) (pluginClient, error) {
	rawBasic, err := rpcClient.Dispense("basic")
	if err != nil {
		return nil, err
	}

	rawResource, err := rpcClient.Dispense("resource")
	if err != nil {
		return nil, err
	}

	rawLive, err := rpcClient.Dispense("live")
	if err != nil {
		return nil, err
	}

	c := ClientV2{}
	if rawBasic != nil {
		if basicClient, ok := rawBasic.(grpcplugin.PluginInfoClient); ok {
			c.PluginInfoClient = basicClient
		}
	}

	if rawResource != nil {
		if resourceClient, ok := rawResource.(grpcplugin.ResourceClient); ok {
			c.ResourceClient = resourceClient
		}
	}

	if rawLive != nil {
		if liveClient, ok := rawLive.(grpcplugin.LiveClient); ok {
			c.LiveClient = liveClient
		}
	}

	return &c, nil
}

func (c *ClientV2) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if c.PluginInfoClient == nil {
		return nil, errors.New("该插件没有实现CheckHealth接口")
	}

	protoContext := backend.ToProto().PluginContext(req.PluginContext)
	protoResp, err := c.PluginInfoClient.CheckHealth(ctx, &pluginv2.CheckHealthRequest{PluginContext: protoContext, Headers: req.Headers})

	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusUnknown,
				Message: "该插件没有实现CheckHealth接口",
			}, nil
		}
		return nil, err
	}

	return backend.FromProto().CheckHealthResponse(protoResp), nil
}

func (c *ClientV2) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if c.ResourceClient == nil {
		return errors.New("该插件没有实现CallResource接口")
	}

	protoReq := backend.ToProto().CallResourceRequest(req)
	protoStream, err := c.ResourceClient.CallResource(ctx, protoReq)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return errors.New("该插件没有实现CallResource接口")
		}

		return fmt.Errorf("%v: %w", "Failed to call resource", err)
	}

	for {
		protoResp, err := protoStream.Recv()
		if err != nil {
			if status.Code(err) == codes.Unimplemented {
				return errors.New("该插件没有实现CallResource接口")
			}

			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("%v: %w", "failed to receive call resource response", err)
		}

		if err := sender.Send(backend.FromProto().CallResourceResponse(protoResp)); err != nil {
			return err
		}
	}
}

func (c *ClientV2) Pub2Channel(ctx context.Context, req *backend.Pub2ChannelRequest) (*backend.Pub2ChannelResponse, error) {
	if c.LiveClient == nil {
		return nil, errors.New("该插件没有实现Live接口1")
	}

	protoContext := backend.ToProto().PluginContext(req.PluginContext)
	protoResp, err := c.LiveClient.Pub2Channel(ctx, &pluginv2.Pub2ChannelRequest{PluginContext: protoContext, Channel: req.Channel, JsonDetails: req.Data})

	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return &backend.Pub2ChannelResponse{
				Status:  backend.PubStatusUnknown,
				Message: "该插件没有实现Live接口2",
			}, nil
		}
		return nil, err
	}

	return backend.FromProto().Pub2ChannelResponse(protoResp), nil
}
