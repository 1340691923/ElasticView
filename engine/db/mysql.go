//MySql引擎层
package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// sqlx 全局变量
var Sqlx *sqlx.DB

// 用squirrel生成sql语句
var SqlBuilder = squirrel.StatementBuilder

type Eq = squirrel.Eq
type Or = squirrel.Or
type Like = squirrel.Like
type Gte = squirrel.GtOrEq
type Lte = squirrel.LtOrEq
type SelectBuilder = squirrel.SelectBuilder
type InsertBuilder = squirrel.InsertBuilder
type UpdateBuilder = squirrel.UpdateBuilder

// NewMySQL 创建一个连接MySQL的实体池
func NewSQLX(dbSource string, maxOpenConns, maxIdleConns int) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("mysql", dbSource)
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
