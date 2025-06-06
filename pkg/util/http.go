package util

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

// GetURLReceiveJSON GET请求 自动解析JSON
func GetURLReceiveJSON(URL string, params url.Values, receive interface{}) error {
	body, err := GetValueURL(URL, params)
	if err != nil {
		return err
	}

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
// Deprecated: 使用Checks.logNil代替
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

	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("body:%v,err:%v", string(body), err)
	}
	return nil
}

// PostJSON POST请求 BODY为JSON格式 ContentType=application/json
func PostJSON(URL string, v interface{}) ([]byte, error) {

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

func GetURLWithHeader(URL string, header map[string]string) ([]byte, error) {

	client := &http.Client{}
	client.Timeout = 3 * time.Minute
	reqest, err := http.NewRequest("GET", URL, nil)

	for k, v := range header {
		reqest.Header.Add(k, v)
	}
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

	err = json.Unmarshal(body, receive)
	if err != nil {
		return fmt.Errorf("error:%v,body{%s}", err, body)
	}
	return nil
}

// PostToJSON POST请求 BODY为json格式
// Deprecated: Please use PostJSON to replace
func PostToJSON(URL string, v interface{}) ([]byte, error) {

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

func GetToken(ctx *gin.Context) (token string) {
	return ctx.GetHeader("X-Token")
}

func DownloadFile(fileURL, dir string) (filename string, err error) {
	// 解析 URL 获取文件名
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	filename = path.Base(parsedURL.Path)

	// 创建目录，如果不存在的话
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 创建目标文件的完整路径
	filePath := filepath.Join(dir, filename)

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	// 获取远程文件
	resp, err := http.Get(fileURL)
	if err != nil {
		out.Close()
		os.Remove(filePath)
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		out.Close()
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// 将内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	out.Close()
	return filename, nil
}

func DownloadFileV2(fileURL, dir string) (filename string, err error) {
	// 解析 URL 获取文件名
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	filename = path.Base(parsedURL.Path)

	// 创建目录，如果不存在的话
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 创建目标文件的完整路径
	filePath := filepath.Join(dir, filename)

	// 获取远程文件
	err = NewDownloader(10, true).Download(fileURL, filename)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to download file: %v", err)
	}

	return filename, nil
}

// GetURLWithHeaders 发起 GET 请求并附带自定义 headers，返回响应内容（字节）
func GetURLWithHeaders(url string, headers map[string]string) ([]byte, error) {
	// 创建客户端（可设置超时）
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	// 添加 Headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func Cors(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET, POST, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Accept-Encoding,Accept-Language,Access-Control-Request-Headers,Access-Control-Request-Method,Connection,Referer,Sec-Fetch-Dest,User-Agent, Origin,Authorization,Content-Type,X-Token,x-token,X-Version,Current-Language")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})

}
