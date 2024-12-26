package orm

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogI struct {
	level gormLogger.LogLevel
	log   hclog.Logger
}

func NewGormLogI(log hclog.Logger) *GormLogI {
	return &GormLogI{log: log}
}

func (g *GormLogI) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	logI := NewGormLogI(g.log)
	logI.level = level
	return logI
}

func (g *GormLogI) Info(ctx context.Context, s string, i ...interface{}) {
	if g.level >= gormLogger.Info {
		g.log.Info(s, i...)
	}
}

func (g *GormLogI) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.level >= gormLogger.Warn {
		g.log.Warn(s, i...)
	}
}

func (g *GormLogI) Error(ctx context.Context, s string, i ...interface{}) {
	if g.level >= gormLogger.Error {
		g.log.Error(s, i...)
	}
}

func (g *GormLogI) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.level >= gormLogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			g.log.Error("%s %s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.log.Error("%s %s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > 200*time.Millisecond && g.level >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", 200*time.Millisecond)
		if rows == -1 {
			g.log.Warn("%s %s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.log.Warn("%s %s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.level <= gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			g.log.Info("%s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.log.Info("%s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
