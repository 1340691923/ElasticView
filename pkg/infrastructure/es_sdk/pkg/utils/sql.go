package utils

import (
	"fmt"
	"gorm.io/gorm"
)

// QuerySQL 通用的 SQL 查询函数
func QuerySQL(db *gorm.DB, sql string, args ...interface{}) (columns []string, results []map[string]interface{}, err error) {
	// 执行 SQL 查询
	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		err = fmt.Errorf("failed to execute SQL: %v", err)
		return
	}
	defer rows.Close()

	// 获取字段名
	columns, err = rows.Columns()
	if err != nil {
		err = fmt.Errorf("failed to get columns: %v", err)
		return
	}

	// 准备结果容器
	results = []map[string]interface{}{}

	// 遍历每一行数据
	for rows.Next() {
		// 创建值的容器
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描数据到容器
		if err = rows.Scan(valuePtrs...); err != nil {
			err = fmt.Errorf("failed to scan row: %v", err)
			return
		}

		// 将每一行数据转换为 map
		rowData := make(map[string]interface{})
		for i, col := range columns {
			rowData[col] = values[i]
		}

		// 添加到结果集
		results = append(results, rowData)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("row iteration error: %v", err)
		return
	}

	// 返回结果
	return
}
