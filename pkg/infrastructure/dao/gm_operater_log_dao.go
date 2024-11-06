package dao

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
)

type GmOperaterLogDao struct {
	orm *sqlstore.SqlStore
}

func NewGmOperaterLogDao(orm *sqlstore.SqlStore) *GmOperaterLogDao {
	return &GmOperaterLogDao{orm: orm}
}
