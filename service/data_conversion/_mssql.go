package data_conversion

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/jmoiron/sqlx"
)

type Mssql struct {
	request.DataxInfoTestLinkReq
}

func NewMssql(data request.DataxInfoTestLinkReq) Datasource {
	return &Mssql{
		data,
	}
}

func (this *Mssql) Ping() error {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v", this.IP, this.Username, this.Pwd, this.Port)
	db, err := sqlx.Connect("mssql", connString)
	if err != nil {
		return err
	}
	_, err = db.Exec("use " + this.DbName)

	return err
}
