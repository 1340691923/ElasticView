package dao

import "github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"

type EsConnectDao struct {
	orm *sqlstore.SqlStore
}
