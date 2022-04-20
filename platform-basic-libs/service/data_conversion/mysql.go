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
	"time"
)

type Mysql struct {
	request.DataxInfoTestLinkReq
}

func (this *Mysql) Transfer(id int, transferReq request.TransferReq) (err error) {
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
	log.Println(sql, args)
	if err != nil {
		return err
	}

	list, err := queryRows(conn, sql, args...)
	if err != nil {
		panic(err)
	}
	log.Println(list)

	/*ll, _ := json.Marshal(list)
	log.Println("list", string(ll))*/
	return errors.New("stop")
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

func queryRows(db *sqlx.DB, sqlStr string, val ...interface{}) (list []map[string]interface{}, err error) {
	var rows *sql.Rows
	rows, err = db.Query(sqlStr, val...)
	if err != nil {
		return
	}
	defer rows.Close()
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 这里需要初始化为空数组，否则在查询结果为空的时候，返回的会是一个未初始化的指针
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}

		ret := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				ret[columns[i]] = nil
			} else {
				switch val := (*scanArgs[i].(*interface{})).(type) {
				case byte:
					ret[columns[i]] = val
					break
				case []byte:
					v := string(val)
					switch v {
					case "\x00": // 处理数据类型为bit的情况
						ret[columns[i]] = 0
					case "\x01": // 处理数据类型为bit的情况
						ret[columns[i]] = 1
					default:
						ret[columns[i]] = v
						break
					}
					break
				case time.Time:
					if val.IsZero() {
						ret[columns[i]] = nil
					} else {
						ret[columns[i]] = val.Format("2006-01-02 15:04:05")
					}
					break
				default:
					ret[columns[i]] = val
				}
			}
		}
		list = append(list, ret)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (this *Mysql) Ping() error {
	conn, err := this.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Ping()
}
