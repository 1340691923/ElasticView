//MySql引擎层
package db

import (
	"fmt"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

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
		db.SetMaxOpenConns(maxOpenConns)
	}

	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
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
