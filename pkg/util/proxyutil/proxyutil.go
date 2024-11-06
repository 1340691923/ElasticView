package proxyutil

import (
	"net"
	"net/http"
)

func PrepareProxyRequest(req *http.Request) {
	req.Header.Del("X-Forwarded-Host")
	req.Header.Del("X-Forwarded-Port")
	req.Header.Del("X-Forwarded-Proto")

	if req.RemoteAddr != "" {
		remoteAddr, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			remoteAddr = req.RemoteAddr
		}
		if req.Header.Get("X-Forwarded-For") != "" {
			req.Header.Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For")+", "+remoteAddr)
		} else {
			req.Header.Set("X-Forwarded-For", remoteAddr)
		}
	}
}

func ClearCookieHeader(req *http.Request, keepCookiesNames []string) {
	var keepCookies []*http.Cookie
	for _, c := range req.Cookies() {
		for _, v := range keepCookiesNames {
			if c.Name == v {
				keepCookies = append(keepCookies, c)
			}
		}
	}

	req.Header.Del("Cookie")
	for _, c := range keepCookies {
		req.AddCookie(c)
	}
}

func SetProxyResponseHeaders(header http.Header) {
	header.Set("Content-Security-Policy", "sandbox")
}
