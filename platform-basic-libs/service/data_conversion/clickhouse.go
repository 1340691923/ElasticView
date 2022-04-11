package data_conversion

import (
	"fmt"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/jmoiron/sqlx"
)

type Clickhouse struct {
	request.DataxInfoTestLinkReq
}

func NewClickhouse(data request.DataxInfoTestLinkReq) Datasource {
	return &Clickhouse{
		data,
	}
}

func (this *Clickhouse) Ping() error {
	dbSource := fmt.Sprintf(
		"tcp://%s:%v?username=%s&password=%s&database=%s&compress=true",
		this.IP,
		this.Port,
		this.Username,
		this.Pwd,
		this.DbName)
	db, err := sqlx.Open("clickhouse", dbSource)
	if err != nil {
		return err
	}

	return db.Ping()
}
