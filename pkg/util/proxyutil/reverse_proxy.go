package proxyutil

import (
	"context"
	"errors"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"

	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// StatusClientClosedRequest A non-standard status code introduced by nginx
// for the case when a client closes the connection while nginx is processing
// the request.
// https://httpstatus.in/499/
const StatusClientClosedRequest = 499

// ReverseProxyOption reverse proxy option to configure a httputil.ReverseProxy.
type ReverseProxyOption func(*httputil.ReverseProxy)

// NewReverseProxy creates a new httputil.ReverseProxy with sane default configuration.
func NewReverseProxy(logger *logger.AppLogger, director func(*http.Request), opts ...ReverseProxyOption) *httputil.ReverseProxy {
	if logger == nil {
		panic("logger cannot be nil")
	}

	if director == nil {
		panic("director cannot be nil")
	}

	p := &httputil.ReverseProxy{
		FlushInterval: time.Millisecond * 200,
		ErrorHandler:  errorHandler(logger),
		ErrorLog:      log.New(&logWrapper{logger: logger}, "", 0),
		Director:      director,
	}

	for _, opt := range opts {
		opt(p)
	}

	origDirector := p.Director
	p.Director = wrapDirector(origDirector)

	if p.ModifyResponse == nil {
		// nolint:bodyclose
		p.ModifyResponse = modifyResponse(logger)
	} else {
		modResponse := p.ModifyResponse
		p.ModifyResponse = func(resp *http.Response) error {
			if err := modResponse(resp); err != nil {
				return err
			}

			// nolint:bodyclose
			return modifyResponse(logger)(resp)
		}
	}

	return p
}

// wrapDirector wraps a director and adds additional functionality.
func wrapDirector(d func(*http.Request)) func(req *http.Request) {
	return func(req *http.Request) {
		d(req)
		PrepareProxyRequest(req)

		// Clear Origin and Referer to avoid CORS issues
		req.Header.Del("Origin")
		req.Header.Del("Referer")
	}
}

// modifyResponse enforces certain constraints on http.Response.
func modifyResponse(logger *logger.AppLogger) func(resp *http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Del("Set-Cookie")
		SetProxyResponseHeaders(resp.Header)
		return nil
	}
}

type timeoutError interface {
	error
	Timeout() bool
}

// errorHandler handles any errors happening while proxying a request and enforces
// certain HTTP status based on the kind of error.
// If client cancel/close the request we return 499 StatusClientClosedRequest.
// If timeout happens while communicating with upstream server we return http.StatusGatewayTimeout.
// If any other error we return http.StatusBadGateway.
func errorHandler(logger *logger.AppLogger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, context.Canceled) {
			logger.Debug("Proxy request cancelled by client")
			w.WriteHeader(StatusClientClosedRequest)
			return
		}

		// nolint:errorlint
		if timeoutErr, ok := err.(timeoutError); ok && timeoutErr.Timeout() {
			logger.Sugar().Errorf("Proxy request timed out", "err", err)
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		logger.Sugar().Errorf("Proxy request failed", "err", err)
		w.WriteHeader(http.StatusBadGateway)
	}
}

type logWrapper struct {
	logger *logger.AppLogger
}

// Write writes log messages as bytes from proxy.
func (lw *logWrapper) Write(p []byte) (n int, err error) {
	withoutNewline := strings.TrimSuffix(string(p), "\n")
	lw.logger.Sugar().Errorf("Proxy request error", "error", withoutNewline)
	return len(p), nil
}

func WithTransport(transport http.RoundTripper) ReverseProxyOption {
	if transport == nil {
		panic("transport cannot be nil")
	}

	return ReverseProxyOption(func(rp *httputil.ReverseProxy) {
		rp.Transport = transport
	})
}

func WithModifyResponse(fn func(*http.Response) error) ReverseProxyOption {
	if fn == nil {
		panic("fn cannot be nil")
	}

	return ReverseProxyOption(func(rp *httputil.ReverseProxy) {
		rp.ModifyResponse = fn
	})
}
