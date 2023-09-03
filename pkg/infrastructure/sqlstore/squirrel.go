package sqlstore

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"strings"
)

var SqlBuilder = squirrel.StatementBuilder

type (
	Eq  = squirrel.Eq
	Or  = squirrel.Or
	And = squirrel.And

	NotEq  = squirrel.NotEq
	Gt     = squirrel.Gt
	Lt     = squirrel.Lt
	GtOrEq = squirrel.GtOrEq
	LtOrEq = squirrel.LtOrEq

	Like          = squirrel.Like
	Gte           = squirrel.GtOrEq
	Lte           = squirrel.LtOrEq
	SelectBuilder = squirrel.SelectBuilder
	InsertBuilder = squirrel.InsertBuilder
	UpdateBuilder = squirrel.UpdateBuilder
)

// 创建分页查询
func CreatePage(page, limit int) uint64 {
	tmp := (page - 1) * limit
	return uint64(tmp)
}

// 创建模糊查询
func CreateLike(column string) string {
	return fmt.Sprint("%", column, "%")
}

func Placeholders(n int) string {
	var b strings.Builder
	for i := 0; i < n-1; i++ {
		b.WriteString("?,")
	}
	if n > 0 {
		b.WriteString("?")
	}
	return b.String()
}
