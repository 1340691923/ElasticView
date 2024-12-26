package util

import (
	"github.com/pkg/errors"
	"github.com/xwb1989/sqlparser"
)

func ExtractTableName(sql string) ([]string, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tableNames := make([]string, 0)
	err = sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node := node.(type) {
		case sqlparser.TableName:
			tableNames = append(tableNames, node.Name.CompliantName())
		}
		return true, nil
	}, stmt)

	return tableNames, err
}
