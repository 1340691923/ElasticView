//MySql引擎层
package db

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/core"
	"github.com/1340691923/ElasticView/pkg/engine/config"
	"github.com/1340691923/ElasticView/pkg/util"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	_ "github.com/logoove/sqlite"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	core.Register(core.MinLevel, "mysql", InitSqlx)
}

// 初始化mysql连接
func InitSqlx() (fn func(), err error) {
	fn = func() {}
	mysql := config.GlobConfig.Mysql
	driverType := config.GlobConfig.DbType
	var dbSource string
	if driverType == config.SqliteDbTyp {
		currDir := util.GetCurrentDirectory()
		dataDir := filepath.Join(currDir, "data")
		if !util.CheckFileIsExist(dataDir) {
			os.MkdirAll(dataDir, os.ModePerm)
		}
		dbSource = filepath.Join(dataDir, config.GlobConfig.Sqlite.DbPath) + "?_loc=Local&_busy_timeout=9999999"
	} else {
		dbSource = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			mysql.Username,
			mysql.Pwd,
			mysql.IP,
			mysql.Port,
			mysql.DbName)
	}

	Sqlx, err = NewSQLX(
		driverType,
		dbSource,
		mysql.MaxOpenConns,
		mysql.MaxIdleConns,
	)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("%v组件初始化成功！连接：%v，最大打开连接数：%v，最大等待连接数:%v",
		driverType,
		dbSource,
		mysql.MaxOpenConns,
		mysql.MaxIdleConns,
	))
	return
}

// sqlx 全局变量
var Sqlx *sqlx.DB

// 用squirrel生成sql语句
var SqlBuilder = squirrel.StatementBuilder

type And = squirrel.And

type NotEq = squirrel.NotEq
type Gt = squirrel.Gt
type Lt = squirrel.Lt
type GtOrEq = squirrel.GtOrEq
type LtOrEq = squirrel.LtOrEq

type Eq = squirrel.Eq
type Or = squirrel.Or
type Like = squirrel.Like
type Gte = squirrel.GtOrEq
type Lte = squirrel.LtOrEq
type SelectBuilder = squirrel.SelectBuilder
type InsertBuilder = squirrel.InsertBuilder
type UpdateBuilder = squirrel.UpdateBuilder

// NewMySQL 创建一个连接MySQL的实体池
func NewSQLX(driverName, dbSource string, maxOpenConns, maxIdleConns int) (db *sqlx.DB, err error) {
	db, err = sqlx.Open(driverName, dbSource)
	if err != nil {
		return
	}
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns) //最大打开的连接数
	}

	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns) //最大空闲连接数
	}
	err = db.Ping()
	if err != nil {
		return
	}
	go func() {
		for {
			err = db.Ping()
			if err != nil {
				log.Println("mysql db can't connect!")
			}
			time.Sleep(time.Minute)
		}
	}()
	return
}

// 创建分页查询
func CreatePage(page, limit int) uint64 {
	tmp := (page - 1) * limit
	return uint64(tmp)
}

// 创建模糊查询
func CreateLike(column string) string {
	return fmt.Sprint("%", column, "%")
}

type Model2MapParmas struct {
	M                interface{}
	NeedZeroByInt    bool
	NeedZeroByString bool
	CreateTimeCol    string
	UpdateTimeCol    string
}

func Model2Map(m Model2MapParmas) (res map[string]interface{}) {
	res = map[string]interface{}{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := json.Marshal(m.M)
	json.Unmarshal(b, &res)
	for k, v := range res {
		switch v.(type) {
		case float64:
			if v.(float64) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case float32:
			if v.(float32) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case int64:
			if v.(int64) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case int16:
			if v.(int16) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case int8:
			if v.(int8) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case int32:
			if v.(int32) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case int:
			if v.(int) == 0 && !m.NeedZeroByInt {
				delete(res, k)
			}
		case string:
			if v.(string) == "" && !m.NeedZeroByString {
				delete(res, k)
			}
		}
	}
	if m.CreateTimeCol != "" {
		res[m.CreateTimeCol] = time.Now().Format(util.TimeFormat)
	}
	if m.UpdateTimeCol != "" {
		res[m.UpdateTimeCol] = time.Now().Format(util.TimeFormat)
	}
	return
}
