package util

import (
	"database/sql"
	"errors"
	"strings"
)

func FilterMysqlNilErr(err error) bool {
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return true
	}
	return false
}

func IsMysqlRepeatError(err error) bool {
	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return true
	}
	return false
}
