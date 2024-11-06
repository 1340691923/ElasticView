// 日志引擎层
package logger

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger = zap.Logger

// 初始化日志 logger
func InitLog(cfg *config.Config) (logger *AppLogger, err error) {

	getWriter := func(filename string, storageDays int) (io.Writer, error) {
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
	zapConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		TimeKey:    "ts",
		NameKey:    "logger",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(util.TimeFormat))
		},
		CallerKey:      "file",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(zapConfig)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	logPath := filepath.Join(cfg.HomePath, "logs")

	if cfg.Log.LogDir != "" {
		logPath = cfg.Log.LogDir
	}

	evlogPath := filepath.Join(logPath, cfg.CmdName)
	esReqLogPath := filepath.Join(logPath, "es_req")
	dbSqlLogPath := filepath.Join(logPath, "db_sql")

	storageDays := 3
	if cfg.Log.StorageDays != 0 {
		storageDays = cfg.Log.StorageDays
	}

	infoWriter, err := getWriter(filepath.Join(evlogPath, "info.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	errWriter, err := getWriter(filepath.Join(evlogPath, "err.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	debugWriter, err := getWriter(filepath.Join(evlogPath, "debug.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	warnWriter, err := getWriter(filepath.Join(evlogPath, "warn.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	esInfoWriter, err := getWriter(filepath.Join(esReqLogPath, "info.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	esErrWriter, err := getWriter(filepath.Join(esReqLogPath, "err.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	esWarnWriter, err := getWriter(filepath.Join(esReqLogPath, "warn.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	dbSqlInfoWriter, err := getWriter(filepath.Join(dbSqlLogPath, "info.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	dbSqlErrWriter, err := getWriter(filepath.Join(dbSqlLogPath, "err.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	dbSqlWarnWriter, err := getWriter(filepath.Join(dbSqlLogPath, "warn.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	var core zapcore.Core
	var esReqCore zapcore.Core
	var dbSqlCore zapcore.Core

	if cfg.DeBug {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), allLevel),
		)
		esReqCore = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(esInfoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(esErrWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(esWarnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), allLevel),
		)

		dbSqlCore = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(dbSqlInfoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(dbSqlErrWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(dbSqlWarnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), allLevel),
		)

	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		)
		esReqCore = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(esInfoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(esErrWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(esWarnWriter), warnLevel),
		)
		dbSqlCore = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(dbSqlInfoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(dbSqlErrWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(dbSqlWarnWriter), warnLevel),
		)

	}

	once.Do(func() {
		EsReqLog = zap.New(esReqCore, zap.AddCaller(), zap.Development())
		DbSqlLog = zap.New(dbSqlCore, zap.AddCaller(), zap.Development())
	})

	log.Println("日志组件加载成功", logPath)
	return zap.New(core, zap.AddCaller(), zap.Development()), nil
}

func InitPluginLog(cfg *config.Config, pluginName string) (logger hclog.Logger, closeWriteCallback func() error, err error) {

	logPath := filepath.Join(cfg.HomePath, "plugin_logs")

	if cfg.Log.PluginLogDir != "" {
		logPath = cfg.Log.PluginLogDir
	}

	pluginlogPath := filepath.Join(logPath, fmt.Sprintf("%s.log", pluginName))

	err = os.MkdirAll(logPath, os.ModePerm) // 0755 是目录权限
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if _, err = os.Stat(pluginlogPath); os.IsNotExist(err) {
		// 文件不存在，创建新文件
		_, err = os.Create(pluginlogPath)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
	} else if err != nil {
		// os.Stat() 出错
		return nil, nil, errors.WithStack(err)
	}

	writer, err := os.OpenFile(pluginlogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	log.Println("插件日志组件加载成功", pluginlogPath)

	return hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: writer,
			Level:  hclog.Debug,
		}), func() error {
			return writer.Close()
		}, nil
}

func InitDebugLog() (logger *AppLogger) {
	return zap.NewExample()
}

func ZapLog2AppLog(logger *zap.Logger) *AppLogger {
	return logger
}

var (
	EsReqLog *AppLogger
	DbSqlLog *AppLogger
	once     *sync.Once
)

func init() {
	once = new(sync.Once)
}