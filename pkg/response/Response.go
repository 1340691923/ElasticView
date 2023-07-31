// 自定义响应 辅助方法层
package response

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/util"
	fiber "github.com/gofiber/fiber/v2"

	. "github.com/1340691923/ElasticView/pkg/my_error"

	"go.uber.org/zap"
)

// 自定义响应方法
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SUCCESS = 0
	ERROR   = 500
)

const (
	SearchSuccess  = "SearchSuccess"
	DeleteSuccess  = "DeleteSuccess"
	OperateSuccess = "OperateSuccess"
	LogoutSuccess  = "LogoutSuccess"
	LinkSuccess    = "LinkSuccess"
	LoginSuccess   = "LoginSuccess"
)

var resMap = map[string]map[string]string{
	"en": {
		SearchSuccess:  "query was successful",
		DeleteSuccess:  "Deleted successfully",
		OperateSuccess: "Operation successful",
		LogoutSuccess:  "Logout successful",
		LinkSuccess:    "Connection successful",
		LoginSuccess:   "Login successful",
	},
	"zh": {
		SearchSuccess:  "查询成功",
		DeleteSuccess:  "删除成功",
		OperateSuccess: "操作成功",
		LogoutSuccess:  "注销成功",
		LinkSuccess:    "连接成功",
		LoginSuccess:   "登录成功",
	},
}

func (this *Response) JsonDealErr(err error) string {

	b, _ := json.Marshal(this.DealErr(err))
	return util.BytesToStr(b)
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
func (this *Response) Success(ctx *fiber.Ctx, msg string, data interface{}) error {
	this.send(ctx, msg, SUCCESS, data)
	return nil
}

// 错误信息
func (this *Response) Error(ctx *fiber.Ctx, err error) error {
	errorTrace := this.getTrace(err)

	myErr := ErrorToErrorCode(err)

	logs.Logger.Error("Error", zap.Strings("err", this.DealErr(myErr)))

	this.send(ctx, myErr.Error(), myErr.Code(), errorTrace)
	return nil
}

// 输出
func (this *Response) send(ctx *fiber.Ctx, msg string, code int, data interface{}) error {
	var res Response
	res.Code = code
	_, ok := resMap[ctx.Get("Current-Language", "zh")][msg]

	if ok {
		res.Msg = resMap[ctx.Get("Current-Language", "zh")][msg]
	} else {
		res.Msg = msg
	}

	res.Data = data
	ctx.Status(http.StatusOK).JSON(res)
	return nil
}

// 输出
func (this *Response) Output(ctx *fiber.Ctx, data interface{}) error {
	ctx.Status(http.StatusOK).JSON(data)
	return nil
}

// 得到trace信息
func (this *Response) getTrace(err error) []string {
	goEnv := os.Getenv("GO_ENV")
	errorTrace := []string{}
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

func (this *Response) DownloadExcel(downloadFileName string, titleList []string, data [][]string, ctx *fiber.Ctx) (err error) {

	var downloadUrl = fmt.Sprintf("data/%v.csv", time.Now().Format("20060102150405"))

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
				logs.Logger.Sugar().Errorf("err", err)
			}
		}()
	}()
	f, err := os.Open(downloadUrl)
	if err != nil {
		logs.Logger.Sugar().Errorf("os.Open failed:", err)
		return
	}
	defer f.Close()

	// 将文件读取出来
	filedata, err := io.ReadAll(f)
	if err != nil {
		logs.Logger.Sugar().Errorf("io.ReadAll failed:", err)
		return
	}
	ctx.Response().Header.Set("Content-Disposition", `attachment; filename="`+downloadFileName+`.xlsx"`)
	ctx.Write(filedata)
	return
}

func (this *Response) DownloadExcel2(downloadFileName string, titleList []interface{}, data [][]interface{}, ctx *fiber.Ctx) (err error) {
	log.Println("download")
	xlsx := excelize.NewFile()

	for index, _ := range data {
		if index == 0 {
			// 如果为0写入新的excel 第一行为字段名称
			xlsx.SetSheetRow("Sheet1", "A1", &titleList)
		}

		//因为index是从0开始，第一行被字段占用，从第二行开始写入整行数据
		var lint = strconv.Itoa(index + 2)
		log.Println("index", index)
		xlsx.SetSheetRow("Sheet1", "A"+lint, &data[index])
	}
	log.Println("download2")
	var downloadUrl = fmt.Sprintf("data/%v.xlsx", time.Now().Format("20060102150405"))

	xlsx.SaveAs(downloadUrl)
	log.Println("download3")
	//defer os.Remove(downloadUrl)
	f, err := os.Open(downloadUrl)
	if err != nil {
		logs.Logger.Sugar().Errorf("os.Open failed:", err)
		return
	}
	defer f.Close()
	log.Println("download4")
	// 将文件读取出来
	filedata, err := io.ReadAll(f)
	if err != nil {
		logs.Logger.Sugar().Errorf("io.ReadAll failed:", err)
		return
	}
	log.Println("download5")
	ctx.Response().Header.Set("Content-Disposition", `attachment; filename="`+downloadFileName+`.xlsx"`)
	ctx.Write(filedata)
	return
}
