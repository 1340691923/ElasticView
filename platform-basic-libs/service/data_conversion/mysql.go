package data_conversion

import (
	"database/sql"
	"fmt"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	request.DataxInfoTestLinkReq
}

func (this *Mysql) Transfer(id int, transferReq request.TransferReq) {
	panic("implement me")
}

func (this *Mysql) GetTableColumns(tableName string) (interface{}, error) {
	conn, err := this.getConn()
	if err != nil {
		return nil, err
	}

	type Res struct {
		Field   string `db:"Field"`
		Type    string `db:"Type"`
		Comment string `db:"Comment"`
	}

	res := []Res{}

	var tmp []struct {
		Field      sql.NullString `db:"Field"`
		Type       sql.NullString `db:"Type"`
		Comment    sql.NullString `db:"Comment"`
		Collation  sql.NullString `db:"Collation"`
		Null       sql.NullString `db:"Null"`
		Key        sql.NullString `db:"Key"`
		Default    sql.NullString `db:"Default"`
		Privileges sql.NullString `db:"Privileges"`
		Extra      sql.NullString `db:"Extra"`
	}

	err = conn.Select(&tmp, fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName))

	for _, d := range tmp {
		t := Res{}
		if d.Field.Valid {
			t.Field = d.Field.String
		}
		if d.Type.Valid {
			t.Type = d.Type.String
		}
		if d.Comment.Valid {
			t.Comment = d.Comment.String
		}
		res = append(res, t)
	}

	return res, err
}

func (this *Mysql) GetTables() ([]string, error) {
	conn, err := this.getConn()
	if err != nil {
		return nil, err
	}
	var list []string
	err = conn.Select(&list, "show tables")
	if err != nil {
		return nil, err
	}
	return list, nil
}

func NewMysql(data request.DataxInfoTestLinkReq) Datasource {
	return &Mysql{
		data,
	}
}

func (this *Mysql) getConn() (*sqlx.DB, error) {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s",
		this.Username,
		this.Pwd,
		this.IP,
		this.Port,
		this.DbName)
	db, err := sqlx.Open("mysql", dbSource)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (this *Mysql) Ping() error {
	conn, err := this.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Ping()
}
