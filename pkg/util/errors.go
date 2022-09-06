package util

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
)

func FilterMysqlNilErr(err error) bool {
	if err != nil && err != sql.ErrNoRows {
		return true
	}
	return false
}

func FilterRedisNilErr(err error) bool {
	if err != nil && err != redis.ErrNil {
		return true
	}
	return false
}
