package data_conversion

import (
	"fmt"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	request.DataxInfoTestLinkReq
}

func NewMysql(data request.DataxInfoTestLinkReq) Datasource {
	return &Mysql{
		data,
	}
}

func (this *Mysql) Ping() error {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s",
		this.Username,
		this.Pwd,
		this.IP,
		this.Port,
		this.DbName)
	db, err := sqlx.Open("mysql", dbSource)
	if err != nil {
		return err
	}

	return db.Ping()
}
