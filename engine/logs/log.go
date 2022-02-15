//日志引擎层
package logs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/1340691923/ElasticView/platform-basic-libs/util"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type Log struct {
	logPath     string
	storageDays int
}

//Options方法
type NewLogOptions func(log *Log)

//设置日志目录
func WithLogPath(logPath string) NewLogOptions {
	return func(log *Log) {
		log.logPath = logPath
	}
}

//设置日志存活天数
func WithStorageDays(storageDays int) NewLogOptions {
	return func(log *Log) {
		log.storageDays = storageDays
	}
}

//App 构造方法
func NewLog(opts ...NewLogOptions) *Log {
	log := &Log{
		logPath:     filepath.Join(util.GetCurrentDirectory(), "logs"),
		storageDays: 3,
	}
	for _, opt := range opts {
		opt(log)
	}
	return log
}

// 初始化日志 logger
func (this *Log) InitLog() (logger *zap.Logger, err error) {

	config := zapcore.EncoderConfig{
		MessageKey: "msg",
		TimeKey:    "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:      "file",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(config)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= zap.InfoLevel
	})

	infoWriter, err := this.getWriter(filepath.Join(this.logPath, "info.log"), this.storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	warnWriter, err := this.getWriter(filepath.Join(this.logPath, "err.log"), this.storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	var core zapcore.Core

	core = zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.Development()), nil
}

func (this *Log) getWriter(filename string, storageDays int) (io.Writer, error) {
	// 生成rotatelogs的Logger 实际生成的文件名 info.log.YYmmddHH
	hook, err := rotatelogs.New(
		filename+".%Y%m%d", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(storageDays)), // 保存3天
		rotatelogs.WithRotationTime(time.Hour*24),                      //切割频率 24小时
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	return hook, nil
}

func Debug(format string, v ...interface{}) {
	log.Println(format, v)
}
