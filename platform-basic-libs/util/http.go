package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	. "github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

//获取真实的IP  1.1.1.1, 2.2.2.2, 3.3.3.3
func CtxClientIP(ctx *fasthttp.RequestCtx) string {
	clientIP := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		//获取最开始的一个 即 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ctx.RemoteIP().String()
}

func GetIp(r *http.Request) string {
	// var r *http.Request
	ip := ClientPublicIP(r)
	if ip == "" {
		ip = ClientIP(r)
	}
	return ip
}

// DoURL 请求URL并且解析JSON格式的返回数据
func DoURL(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetURL 请求URL
func GetURL(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetURL 请求URL
func CtxGetURL(URL string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源

	// 默认是application/x-www-form-urlencoded
	req.Header.SetMethod("GET")

	req.SetRequestURI(URL)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}
	b := resp.Body()
	return b, nil
}

// GetValueURL 请求URL 附带参数
func GetValueURL(URL string, params url.Values) ([]byte, error) {
	if params == nil {
		return GetURL(URL)
	}
	resp, err := http.Get(fmt.Sprint(URL, "?", params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func CtxGetValueURL(URL string, params url.Values) ([]byte, error) {
	if params == nil {
		return CtxGetURL(URL)
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源

	// 默认是application/x-www-form-urlencoded
	req.Header.SetMethod("GET")

	req.SetRequestURI(fmt.Sprint(URL, "?", params.Encode()))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}
	b := resp.Body()
	return b, nil
}

// GetURLReceiveJSON GET请求 自动解析JSON
func GetURLReceiveJSON(URL string, params url.Values, receive interface{}) error {
	body, err := GetValueURL(URL, params)
	if err != nil {
		return err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %s, %v", body, err)
	}
	return nil
}

func CtxGetURLReceiveJSON(URL string, params url.Values, receive interface{}) error {
	body, err := CtxGetValueURL(URL, params)
	if err != nil {
		return err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %s, %v", body, err)
	}
	return nil
}

// PostURL 请求URL
func PostURL(URL string, params url.Values) ([]byte, error) {
	resp, err := http.PostForm(URL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// 检查http请求中是否包含所需参数
// Deprecated: 使用CheckNotNil代替
func CheckParam(hr *http.Request, args ...string) string {
	if strings.ToUpper(hr.Method) == "GET" {
		for _, val := range args {
			rs := hr.FormValue(val)
			if StringIsEmpty(rs) {
				return val
			}
		}
		return ""
	} else if strings.ToUpper(hr.Method) == "POST" { //post
		for _, val := range args {

			rs := hr.PostFormValue(val)
			if StringIsEmpty(rs) {
				return val
			}
		}
		return ""
	} else {
		return hr.Method
	}
}

// PostURLReceiveJSON POST请求  自动解析JSON
func PostURLReceiveJSON(URL string, params url.Values, receive interface{}) error {
	body, err := PostURL(URL, params)
	if err != nil {
		return err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("body:%v,err:%v", string(body), err)
	}
	return nil
}

// PostURLReceiveJSON POST请求  自动解析JSON
func PostMapReceiveJSON(URL string, maps map[string]string, receive interface{}) error {
	params := url.Values{}
	for k, v := range maps {
		params.Set(k, v)
	}
	body, err := PostURL(URL, params)
	if err != nil {
		return err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("body:%v,err:%v", string(body), err)
	}
	return nil
}

// PostJSON POST请求 BODY为JSON格式 ContentType=application/json
func PostJSON(URL string, v interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(URL, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// PostJSON POST请求 BODY为JSON格式 ContentType=application/json
func GetJSON(URL string, v interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", URL, bytes.NewReader(b))
	reqest.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(reqest)
	//resp, err := http.Post(URL, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// PostJSONReceiveJSON POST请求 BODY为JSON格式 ContentType=application/json 自动解析JSON
func PostJSONReceiveJSON(URL string, send, receive interface{}) error {
	body, err := PostJSON(URL, send)
	if err != nil {
		return err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("error:%v,body{%s}", err, body)
	}
	return nil
}

// PostToJSON POST请求 BODY为json格式
// Deprecated: Please use PostJSON to replace
func PostToJSON(URL string, v interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// CheckNotNil 检查HTTP参数是否为空
func CheckNotNil(r *http.Request, args ...string) error {
	if args == nil || r == nil {
		return nil
	}

	switch r.Method {
	case "GET":
		query := r.URL.Query()
		for _, v := range args {
			if strings.TrimSpace(query.Get(v)) == "" {
				return fmt.Errorf("param(%s) is invalid", v)
			}
		}
	case "POST":
		for _, v := range args {
			if strings.TrimSpace(r.PostFormValue(v)) == "" {
				return fmt.Errorf("param(%s) is invalid", v)
			}
		}
	default:
		return errors.New("r.Method is not GET or POST")
	}
	return nil
}

// StringIsEmpty 判断是否有值为空或null或(null)
func StringIsEmpty(s ...string) bool {
	var str string
	for _, v := range s {
		str = strings.TrimSpace(v)
		if v == "" || strings.EqualFold(str, "(null)") || strings.EqualFold(str, "null") {
			return true
		}
	}
	return false
}

// WriteJSON 写入json字符串
func WriteJSON(w io.Writer, v interface{}) (int, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	return w.Write(b)
}

// GetRemoteIP 获取IP
func GetRemoteIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	var ip = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip == "" {
		ip, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	}
	return ip
}

// CheckRemoteIP 验证IP
// in ips return true
func CheckRemoteIP(r *http.Request, ips ...string) bool {
	if r == nil {
		return false
	}
	var ip = GetRemoteIP(r)
	for _, v := range ips {
		if ip == v {
			return true
		}
	}
	return false
}

var regIPv4 = regexp.MustCompile(
	`^(((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))\.){3}((\d{1,2})|(1\d{1,2})|(2[0-4]\d)|(25[0-5]))$`,
)

// IsValidIPv4 验证是否为合法的ipv4
func IsValidIPv4(ip string) bool {
	return regIPv4.MatchString(ip)
}

// FormIntDefault 获取Form参数 如果出错则返回默认值
func FormIntDefault(r *http.Request, key string, def int) int {
	i, err := strconv.Atoi(r.FormValue(key))
	if err != nil {
		return def
	}
	return i
}

// FormIntDefault 获取Form参数 如果出错则返回默认值
func CtxFormIntDefault(ctx *Ctx, key string, def int) int {
	i, err := strconv.Atoi(ctx.FormValue(key))
	if err != nil {
		return def
	}
	return i
}

// FormIntSliceDefault 获取Form参数 如果出错则返回默认值
func FormIntSliceDefault(r *http.Request, key, sep string, def []int) []int {
	var i int
	var err error
	var rlt []int
	for _, v := range strings.Split(r.FormValue(key), sep) {
		i, err = strconv.Atoi(v)
		if err != nil {
			continue
		}
		rlt = append(rlt, i)
	}
	if rlt == nil {
		return def
	}
	return rlt
}

// FormFileValue 快速获取表单提交的文件
// 也用于处理同表单一起提交的信息
func FormFileValue(r *http.Request, key string) (string, error) {
	f, _, err := r.FormFile(key)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FormFileValues 快速获取表单提交的文件
// 也用于处理同表单一起提交的信息
func FormFileValues(r *http.Request, key string) ([]string, error) {
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return nil, err
		}
	}
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if fhs := r.MultipartForm.File[key]; len(fhs) > 0 {
			var rlt = make([]string, 0, len(fhs))
			for i := range fhs {
				f, err := fhs[i].Open()
				if err != nil {
					return nil, err
				}

				b, err := ioutil.ReadAll(f)
				f.Close()
				if err != nil {
					return nil, err
				}
				rlt = append(rlt, string(b))
			}
			return rlt, nil
		}
	}
	return nil, http.ErrMissingFile
}

func GetToken(ctx *Ctx) (token string) {
	return ctx.Get("X-Token")
}
