package plugins

import (
	"context"
	"errors"
	"fmt"
	logger2 "github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/util/proxyutil"
	"github.com/1340691923/ElasticView/pkg/util/response"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type DataSourcePlugin struct {
	ginCtx    *gin.Context
	rpcPlugin backend.CallResourceHandler
	log       *logger2.AppLogger
}

func NewDataSourcePlugin(ginCtx *gin.Context, rpcPlugin backend.CallResourceHandler, log *logger2.AppLogger) *DataSourcePlugin {
	return &DataSourcePlugin{ginCtx: ginCtx, rpcPlugin: rpcPlugin, log: log}
}

func (this *DataSourcePlugin) CallPluginResource() {
	req, err := this.pluginResourceRequest()
	if err != nil {
		this.ginCtx.String(http.StatusBadRequest, "Failed for create plugin resource request", err)
		return
	}

	if err = this.makePluginResourceRequest(req); err != nil {
		handleCallResourceError(err, this.ginCtx)
	}

}

func (this *DataSourcePlugin) makePluginResourceRequest(req *http.Request) error {
	keepCookieModel := struct {
		KeepCookies []string `json:"keepCookies"`
	}{}

	proxyutil.ClearCookieHeader(req, keepCookieModel.KeepCookies)
	proxyutil.PrepareProxyRequest(req)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	crReq := &backend.CallResourceRequest{
		Path:    req.URL.Path,
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: req.Header,
		Body:    body,
	}

	childCtx, cancel := context.WithCancel(req.Context())
	defer cancel()
	stream := newCallResourceResponseStream(childCtx)

	var wg sync.WaitGroup
	wg.Add(1)

	defer func() {
		if err := stream.Close(); err != nil {
			this.log.Sugar().Errorf("Failed to close plugin resource stream", "err", err)
		}
		wg.Wait()
	}()

	var flushStreamErr error
	go func() {
		flushStreamErr = this.flushStream(stream, this.ginCtx.Writer)
		wg.Done()
	}()

	if err := this.CallResource(crReq, stream); err != nil {
		return err
	}

	return flushStreamErr
}

type callResourceClientResponseStream interface {
	Recv() (*backend.CallResourceResponse, error)
	Close() error
}

func (this *DataSourcePlugin) flushStream(stream callResourceClientResponseStream, w http.ResponseWriter) error {
	processedStreams := 0

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			if processedStreams == 0 {
				return errors.New("received empty resource response")
			}
			return nil
		}
		if err != nil {
			if processedStreams == 0 {
				return fmt.Errorf("%v: %w", "failed to receive response from resource call", err)
			}

			this.log.Sugar().Errorf("Failed to receive response from resource call", "err", err)
			return stream.Close()
		}

		// Expected that headers and status are only part of first stream
		if processedStreams == 0 && resp.Headers != nil {
			// Make sure a content type always is returned in response
			if _, exists := resp.Headers["Content-Type"]; !exists {
				resp.Headers["Content-Type"] = []string{"application/json"}
			}

			for k, values := range resp.Headers {
				// Due to security reasons we don't want to forward
				// cookies from a backend plugin to clients/browsers.
				if k == "Set-Cookie" {
					continue
				}

				for _, v := range values {
					// TODO: Figure out if we should use Set here instead
					// nolint:gocritic
					w.Header().Add(k, v)
				}
			}

			proxyutil.SetProxyResponseHeaders(w.Header())

			w.WriteHeader(resp.Status)
		}

		if _, err := w.Write(resp.Body); err != nil {
			this.log.Sugar().Errorf("Failed to write resource response", "err", err)
		}

		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		processedStreams++
	}
}

func (this *DataSourcePlugin) pluginResourceRequest() (*http.Request, error) {
	clonedReq := this.ginCtx.Request.Clone(this.ginCtx.Request.Context())
	rawURL := this.ginCtx.Param("action")

	if clonedReq.URL.RawQuery != "" {
		rawURL += "?" + clonedReq.URL.RawQuery
	}
	urlPath, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	clonedReq.URL = urlPath
	return clonedReq, nil
}

type callResourceResponseStream struct {
	ctx    context.Context
	stream chan *backend.CallResourceResponse
	closed bool
}

func newCallResourceResponseStream(ctx context.Context) *callResourceResponseStream {
	return &callResourceResponseStream{
		ctx:    ctx,
		stream: make(chan *backend.CallResourceResponse),
	}
}

func (s *callResourceResponseStream) Send(res *backend.CallResourceResponse) error {
	if s.closed {
		return errors.New("cannot send to a closed stream")
	}

	select {
	case <-s.ctx.Done():
		return errors.New("cancelled")
	case s.stream <- res:
		return nil
	}
}

func (s *callResourceResponseStream) Recv() (*backend.CallResourceResponse, error) {
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case res, ok := <-s.stream:
		if !ok {
			return nil, io.EOF
		}
		return res, nil
	}
}

func (s *callResourceResponseStream) Close() error {
	if s.closed {
		return errors.New("cannot close a closed stream")
	}

	close(s.stream)
	s.closed = true
	return nil
}

func handleCallResourceError(err error, reqCtx *gin.Context) {
	response.JsonApiErr(reqCtx, 500, "Failed to call resource", err)
}

func (this DataSourcePlugin) CallResource(req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	err := this.rpcPlugin.CallResource(this.ginCtx, req, sender)
	if err != nil {
		return fmt.Errorf("%v: %w", "Failed to call resource", err)
	}

	return nil
}
