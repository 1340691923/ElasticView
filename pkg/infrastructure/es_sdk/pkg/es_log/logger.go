package es_log

import (
	"bufio"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	Enable             bool
	EnableRequestBody  bool
	EnableResponseBody bool
	Logger             *zap.Logger
}

func NewLogger(enable, enableRequestBody, enableResponseBody bool, logger *zap.Logger) *Logger {
	return &Logger{Enable: enable, EnableRequestBody: enableRequestBody, EnableResponseBody: enableResponseBody, Logger: logger}
}

func (l *Logger) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	if !l.Enable {
		return nil
	}
	bsize := 200
	var b = bytes.NewBuffer(make([]byte, 0, bsize))
	var v = make([]byte, 0, bsize)

	appendTime := func(t time.Time) {
		v = v[:0]
		v = t.AppendFormat(v, time.RFC3339)
		b.Write(v)
	}

	appendQuote := func(s string) {
		v = v[:0]
		v = strconv.AppendQuote(v, s)
		b.Write(v)
	}

	appendInt := func(i int64) {
		v = v[:0]
		v = strconv.AppendInt(v, i, 10)
		b.Write(v)
	}

	port := req.URL.Port()

	b.WriteRune('{')
	// -- Timestamp
	b.WriteString(`"@timestamp":"`)
	appendTime(start.UTC())
	b.WriteRune('"')
	// -- Event
	b.WriteString(`,"event":{`)
	b.WriteString(`"lose time":`)
	appendQuote(dur.Truncate(time.Millisecond).String())
	b.WriteString(`,"ev_user_id":`)
	appendQuote(req.Header.Get("ev_user_id"))
	b.WriteRune('}')
	// -- URL
	b.WriteString(`,"url":{`)
	b.WriteString(`"scheme":`)
	appendQuote(req.URL.Scheme)
	b.WriteString(`,"domain":`)
	appendQuote(req.URL.Hostname())
	if port != "" {
		b.WriteString(`,"port":`)
		b.WriteString(port)
	}
	b.WriteString(`,"path":`)
	appendQuote(req.URL.Path)
	b.WriteString(`,"query":`)
	appendQuote(req.URL.RawQuery)
	b.WriteRune('}') // Close "url"
	// -- HTTP
	b.WriteString(`,"http":`)
	// ---- Request
	b.WriteString(`{"request":{`)
	b.WriteString(`"method":`)
	appendQuote(req.Method)
	if l.RequestBodyEnabled() && req != nil && req.Body != nil && req.Body != http.NoBody {
		var buf bytes.Buffer
		if req.GetBody != nil {
			b, _ := req.GetBody()
			buf.ReadFrom(b)
		} else {
			buf.ReadFrom(req.Body)
		}

		b.Grow(buf.Len() + 8)
		b.WriteString(`,"body":`)
		appendQuote(buf.String())
	}
	b.WriteRune('}') // Close "http.request"
	// ---- Response
	b.WriteString(`,"response":{`)
	b.WriteString(`"status_code":`)
	appendInt(int64(resStatusCode(res)))
	if l.ResponseBodyEnabled() && res != nil && res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)

		b.Grow(buf.Len() + 8)
		b.WriteString(`,"body":`)
		appendQuote(buf.String())
	}
	b.WriteRune('}') // Close "http.response"
	b.WriteRune('}') // Close "http"
	// -- Error
	if err != nil {
		b.WriteString(`,"error":{"message":`)
		appendQuote(err.Error())
		b.WriteRune('}') // Close "error"
	}
	b.WriteRune('}')
	l.Logger.Sugar().Infof("es请求日志", b.String())
	return nil
}

func resStatusCode(res *http.Response) int {
	if res == nil {
		return -1
	}
	return res.StatusCode
}

func logBodyAsText(strBuilder *strings.Builder, body io.Reader, prefix string) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		s := scanner.Text()
		if s != "" {
			strBuilder.WriteString(fmt.Sprintf("%s %s\n", prefix, s))
		}
	}
}

// RequestBodyEnabled returns true when the request body should be logged.
func (l *Logger) RequestBodyEnabled() bool { return l.EnableRequestBody }

// ResponseBodyEnabled returns true when the response body should be logged.
func (l *Logger) ResponseBodyEnabled() bool { return l.EnableResponseBody }
