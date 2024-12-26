package orm

import (
	"fmt"
)

// 创建分页查询
func CreatePage(page, limit int) int {
	tmp := (page - 1) * limit
	return int(tmp)
}

// 创建模糊查询
func CreateLike(column string) string {
	return fmt.Sprint("%", column, "%")
}
