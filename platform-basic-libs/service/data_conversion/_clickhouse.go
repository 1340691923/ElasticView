package data_conversion

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type Clickhouse struct {
	request.DataxInfoTestLinkReq
}

func (this *Clickhouse) Transfer(id int, transferReq request.TransferReq) (err error) {
	var (
		page  = 1
		limit = 20
	)
	conn, err := this.getConn()
	if err != nil {
		return err
	}
	sql, args, err := db.SqlBuilder.Select(strings.Join(transferReq.Cols.TableCols, ",")).
		From(transferReq.SelectTable).
		Limit(uint64(limit)).
		Offset(db.CreatePage(page, limit)).ToSql()
	if err != nil {
		return err
	}

	rows, err := conn.Query(sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{}

	for rows.Next() {
		err := rows.Scan(cache...)
		if err != nil {
			return err
		}
		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{})
		}
		list = append(list, item)
	}
	log.Println(list)

	/*ll, _ := json.Marshal(list)
	log.Println("list", string(ll))*/
	return errors.New("stop")
}

func (this *Clickhouse) GetTableColumns(tableName string) (interface{}, error) {
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
		Name               sql.NullString `db:"name"`
		Type               sql.NullString `db:"type"`
		Comment            sql.NullString `db:"comment"`
		DefaultType        sql.NullString `db:"default_type"`
		Default_expression sql.NullString `db:"default_expression"`
		Codec_expression   sql.NullString `db:"codec_expression"`
		Ttl_expression     sql.NullString `db:"ttl_expression"`
	}

	err = conn.Select(&tmp, fmt.Sprintf("desc %s", tableName))

	for _, d := range tmp {
		t := Res{}
		if d.Name.Valid {
			t.Field = d.Name.String
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

func (this *Clickhouse) GetTables() ([]string, error) {
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

func (this *Clickhouse) getConn() (*sqlx.DB, error) {
	dbSource := fmt.Sprintf(
		"tcp://%s:%v?username=%s&password=%s&database=%s&compress=true",
		this.IP,
		this.Port,
		this.Username,
		this.Pwd,
		this.DbName)
	db, err := sqlx.Open("clickhouse", dbSource)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewClickhouse(data request.DataxInfoTestLinkReq) Datasource {
	return &Clickhouse{
		data,
	}
}

func (this *Clickhouse) Ping() error {
	conn, err := this.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Ping()
}
