// 日志引擎层
package logger

import (
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"io"
	"log"
	"os"
	"path/filepath"

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

	logPath = filepath.Join(logPath, cfg.CmdName)

	storageDays := 3
	if cfg.Log.StorageDays != 0 {
		storageDays = cfg.Log.StorageDays
	}

	infoWriter, err := getWriter(filepath.Join(logPath, "info.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	errWriter, err := getWriter(filepath.Join(logPath, "err.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	debugWriter, err := getWriter(filepath.Join(logPath, "debug.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	warnWriter, err := getWriter(filepath.Join(logPath, "warn.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}

	allWriter, err := getWriter(filepath.Join(logPath, "all.log"), storageDays)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("日志启动异常:%s", err))
	}
	var core zapcore.Core
	if cfg.DeBug {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel),

			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), allLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel),
		)
	}
	log.Println("日志组件加载成功", logPath)
	return zap.New(core, zap.AddCaller(), zap.Development()), nil
}
