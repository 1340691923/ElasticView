package data_conversion

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
	"time"

	"strings"
)

type Clickhouse struct {
	request.DataxInfoTestLinkReq
}

func (this *Clickhouse) Transfer(id int, transferReq *request.TransferReq) (err error) {
	var page  = 1

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

	err = conn.Ping()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
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

	esConnect, err := escache.GetEsClientByID(transferReq.EsConnect)

	if err != nil {
		updateDataXListStatus(id, 0, 0, Error, err.Error())
		return err
	}
	if transferReq.EsDocID != "" {
		transferReq.Cols.TableCols = append(transferReq.Cols.TableCols, transferReq.EsDocID)
	}
	_ = func(offset uint64, limit int) string {
		sql := fmt.Sprintf(`SELECT %s FROM %s WHERE id >= (select id from %s limit %v, 1) limit %v`,
			strings.Join(transferReq.Cols.TableCols, ","),
			transferReq.SelectTable,
			transferReq.SelectTable,
			offset,
			limit,
		)
		return sql
	}

	go func() {
		sqlFn := func(offset uint64, limit int) string {
			sql := fmt.Sprintf(`SELECT %s FROM %s limit %v,%v`,
				strings.Join(transferReq.Cols.TableCols, ","),
				transferReq.SelectTable,
				offset,
				limit,
			)
			return sql
		}

		switch esConnect.Version {
		case 6:
			esConn, err := escache.NewEsClientV6(esConnect)
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
			esConn, err := escache.NewEsClientV7(esConnect)
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
