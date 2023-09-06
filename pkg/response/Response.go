// 自定义响应 辅助方法层
package response

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	. "github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// 自定义响应方法
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	log  *logger.AppLogger
}

func NewResponse(log *logger.AppLogger) *Response {
	return &Response{log: log.Named("response")}
}

const (
	SUCCESS = 0
	ERROR   = 500
)

const (
	SearchSuccess       = "查询成功"
	DeleteSuccess       = "删除成功"
	OperateSuccess      = "操作成功"
	LogoutSuccess       = "注销成功"
	ChangeLayoutSuccess = "修改布局成功"
)

func (this *Response) JsonDealErr(err error) string {
	b, _ := json.Marshal(this.DealErr(err))
	return string(b)
}

// trace
func (this *Response) DealErr(err error) (errorTrace []string) {
	errorTrace = append(errorTrace, err.Error())
	if err != nil {
		for i := 1; ; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			f := runtime.FuncForPC(pc)
			if f.Name() != "runtime.main" && f.Name() != "runtime.goexit" && !strings.Contains(file, "Response.go") {
				errStrings := "文件名:" + file + ",行数:" + strconv.Itoa(line) + ",函数名:" + f.Name()
				errorTrace = append(errorTrace, errStrings)
			}
		}
	}
	return errorTrace
}

// 正确信息
func (this *Response) Success(ctx *gin.Context, msg string, data interface{}) error {
	this.Msg = msg
	this.Data = data
	this.send(ctx, SUCCESS)
	return nil
}

// 错误信息
func (this *Response) FastError(write io.Writer, err error) error {
	myErr := ErrorToErrorCode(err)

	this.Output(write, map[string]interface{}{
		"code": myErr.Code(),
		"msg":  myErr.Error(),
	})
	return nil
}

// 错误信息
func (this *Response) Error(ctx *gin.Context, err error) error {
	errorTrace := this.getTrace(err)

	myErr := ErrorToErrorCode(err)

	this.log.Error("Error", zap.Strings("err", this.DealErr(myErr)))

	this.Msg = myErr.Error()
	this.Data = errorTrace
	this.send(ctx, myErr.Code())
	return nil
}

// 输出
func (this *Response) send(ctx *gin.Context, code int) error {
	this.Code = code
	var err error
	if this.Code != 0 {
		ctx.JSON(http.StatusAccepted, this)
	} else {
		ctx.JSON(http.StatusOK, this)
	}

	if err != nil {
		ctx.JSON(http.StatusAccepted, map[string]interface{}{"msg": err, "code": 500})
	}
	return nil
}

// 输出
func (this *Response) Output(write io.Writer, data map[string]interface{}) error {
	b, _ := json.Marshal(data)
	write.Write(b)
	return nil
}

// 得到trace信息
func (this *Response) getTrace(err error) []string {
	goEnv := os.Getenv("GO_ENV")
	var errorTrace []string
	if goEnv == "product" {
		errorTrace = this.DealErr(err)
	}
	return errorTrace
}

// 处理异常（业务异常和默认异常）
func ErrorToErrorCode(err error) *MyError {
	if err == nil {
		return nil
	}
	errorCode, ok := err.(*MyError)

	if ok {
		return errorCode
	}
	return NewError(err.Error(), ERROR).(*MyError)
}

func (this *Response) ReturnValOrNull(value, empty interface{}) interface{} {
	var vValue = reflect.ValueOf(value)
	if value == nil || (vValue.Kind() == reflect.Slice && vValue.Len() == 0) {
		return empty
	}
	return value
}

func (this *Response) SliceReturnValOrNull(value []string, empty interface{}) interface{} {
	if value == nil || len(value) == 0 {
		return empty
	}
	return value
}

func (this *Response) DownloadExcel(downloadFileName string, titleList []string, data [][]string, ctx *gin.Context, log *logger.AppLogger) (err error) {

	var downloadUrl = fmt.Sprintf("data/%v.csv", time.Now().Format("20060102150405"))

	if !util.CheckFileIsExist(downloadUrl) {
		os.Create(downloadUrl)
	}

	file, err := os.OpenFile(downloadUrl, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
		return
	}
	defer file.Close()
	// 写入UTF-8 BOM，防止中文乱码
	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	w.Write(titleList)

	for _, d := range data {
		w.Write(d)
	}
	// 写文件需要flush，不然缓存满了，后面的就写不进去了，只会写一部分
	w.Flush()

	defer func() {
		go func() {
			time.Sleep(5 * time.Second)
			err := os.Remove(downloadUrl)
			if err != nil {
				log.Sugar().Errorf("err", err)
			}
		}()
	}()
	f, err := os.Open(downloadUrl)
	if err != nil {
		log.Sugar().Errorf("os.Open failed:", err)
		return
	}
	defer f.Close()

	// 将文件读取出来
	filedata, err := io.ReadAll(f)
	if err != nil {
		log.Sugar().Errorf("io.ReadAll failed:", err)
		return
	}
	ctx.Header("Content-Disposition", `attachment; filename="`+downloadFileName+`.xlsx"`)
	ctx.Writer.Write(filedata)
	return
}
