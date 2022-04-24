package data_conversion

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
	"strings"
	"time"
)

type Mysql struct {
	request.DataxInfoTestLinkReq
}

func (this *Mysql) Transfer(id int, transferReq *request.TransferReq) (err error) {
	var page = 1

	ctx, cancel := context.WithCancel(context.Background())

	conn, err := this.getConn()
	if err != nil {
		updateDataXListStatus(id, 0, 0, Error, err.Error())
		return err
	}
	maxOpenConns := transferReq.MaxOpenConns
	maxIdleConns := transferReq.MaxIdleConns

	if maxOpenConns > 0 {
		conn.SetMaxOpenConns(maxOpenConns) //最大打开的连接数
	}

	if maxIdleConns > 0 {
		conn.SetMaxIdleConns(maxIdleConns) //最大空闲连接数
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("关闭连接了")
				conn.Close()
				return
			default:
				err = conn.Ping()
				if err != nil {
					log.Println(fmt.Sprintf(`"mysql db can't connect! 数据抽取任务id :%v`, id))
				}
				time.Sleep(time.Minute)
			}
		}
	}()

	count := 0
	err = db.SqlBuilder.Select("count(1)").From(transferReq.SelectTable).RunWith(conn).Scan(&count)
	if err != nil {
		updateDataXListStatus(id, 0, 0, Error, err.Error())
		return err
	}

	if count == 0 {
		updateDataXListStatus(id, 0, 0, Error, fmt.Sprintf("该表【%s】没有数据", transferReq.SelectTable))
		return errors.New(fmt.Sprintf("该表【%s】没有数据", transferReq.SelectTable))
	}

	if transferReq.GoNum == 0 {
		transferReq.GoNum = 30
	}

	limit := transferReq.BufferSize

	ts := GetTaskInstance()

	ts.SetCancelFunc(id, cancel)

	length := int(math.Ceil(float64(float64(count) / float64(limit))))
	lastLimit := count % limit

	esConnect, err := es.GetEsClientByID(transferReq.EsConnect)

	if err != nil {
		updateDataXListStatus(id, 0, 0, Error, err.Error())
		return err
	}


	go func() {

		var sqlFn  func(offset uint64, limit int) string
		if transferReq.AutoIncrementId != ""{
			sqlFn = func(offset uint64, limit int) string {
				sql := fmt.Sprintf(`SELECT %s FROM %s WHERE %s >= (select %s from %s limit %v, 1) limit %v`,
					strings.Join(transferReq.Cols.TableCols, ","),
					transferReq.SelectTable,
					transferReq.AutoIncrementId,
					transferReq.AutoIncrementId,
					transferReq.SelectTable,
					offset,
					limit,
				)
				return sql
			}
		}else{
			sqlFn = func(offset uint64, limit int) string {
				sql := fmt.Sprintf(`SELECT %s FROM %s limit %v,%v`,
					strings.Join(transferReq.Cols.TableCols, ","),
					transferReq.SelectTable,
					offset,
					limit,
				)
				return sql
			}
		}



		switch esConnect.Version {
		case 6:
			esConn, err := es.NewEsClientV6(esConnect)
			if err != nil {
				updateDataXListStatus(id, 0, 0, Error, err.Error())
				return
			}

			err = transferEsV6(
				id, transferReq, page, limit, lastLimit,
				length, count, sqlFn, ctx, conn, esConn,
			)
			if err != nil {
				updateDataXListStatus(id, 0, 0, Error, err.Error())
				return
			}
		case 7:
			fallthrough
		case 8:
			esConn, err := es.NewEsClientV7(esConnect)
			if err != nil {
				updateDataXListStatus(id, 0, 0, Error, err.Error())
				return
			}

			err = transferEsV7(
				id, transferReq, page, limit, lastLimit,
				length, count, sqlFn, ctx, conn, esConn,
			)
			if err != nil {
				updateDataXListStatus(id, 0, 0, Error, err.Error())
				return
			}

		}
	}()

	return nil
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
