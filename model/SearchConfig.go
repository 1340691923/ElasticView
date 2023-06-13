// 数据库实体层
package model

import (
	"github.com/1340691923/ElasticView/pkg/engine/db"
)

/*
	http://sql2struct.atotoa.com 根据表结构生成 go结构体

GmDslHistoryModel DSL历史记录
*/
type SearchConfig struct {
	ID         int    `db:"id" json:"id"`
	IndexName  string `db:"index_name" json:"indexName"`
	Remark     string `db:"remark" json:"remark"`
	InputCols  string `json:"input_cols" db:"input_cols"`
	OutputCols string `json:"output_cols" db:"output_cols"`

	EsConnect int    `db:"es_connect" json:"-"`
	Updated   string `db:"updated" json:"updated"`
	Created   string `db:"created" json:"created"`
	Limit     int    `json:"-" db:"-"`
	Page      int    `json:"-" db:"-"`
}

// 表名
func (this *SearchConfig) TableName() string {
	return "search_index_config"
}

// Insert
func (this *SearchConfig) Insert() (err error) {
	return
}

// Clean
func (this *SearchConfig) Clean() (err error) {
	return
}

// List
func (this *SearchConfig) List() (res []SearchConfig, err error) {
	builder := db.SqlBuilder.
		Select("*").
		From(this.TableName()).
		Where(db.Eq{"es_connect": this.EsConnect}).
		OrderBy("id desc").
		Limit(uint64(this.Limit)).Offset(db.CreatePage(this.Page, this.Limit))

	sql, args, err := builder.ToSql()

	if err != nil {
		return
	}
	err = db.Sqlx.Select(&res, sql, args...)

	if err != nil {
		return
	}

	return
}

func (this *SearchConfig) All() (res []SearchConfig, err error) {
	builder := db.SqlBuilder.
		Select("*").
		From(this.TableName()).
		Where(db.Eq{"es_connect": this.EsConnect}).
		OrderBy("id desc")
	sql, args, err := builder.ToSql()

	if err != nil {
		return
	}
	err = db.Sqlx.Select(&res, sql, args...)

	if err != nil {
		return
	}

	return
}

// Count
func (this *SearchConfig) Count() (count int, err error) {
	builder := db.SqlBuilder.
		Select("count(*)").
		From(this.TableName()).
		Where(db.Eq{"es_connect": this.EsConnect})

	sql, args, err := builder.ToSql()
	if err != nil {
		return
	}
	err = db.Sqlx.Get(&count, sql, args...)

	if err != nil {
		return
	}

	return
}
