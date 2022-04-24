package data_conversion

import (
	"context"
	"database/sql"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"github.com/jmoiron/sqlx"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	Success = "数据导入成功"
	Error   = "数据导入失败"
	Running = "数据正在导入中..."
	Cancel  = "数据导入任务取消"
)

func queryRows(table2EsMap map[string]string, db *sqlx.DB, sqlStr string, val ...interface{}) (list []map[string]interface{}, err error) {

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

	for index := range columns {
		columns[index] = table2EsMap[columns[index]]
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

func transferEsV6(
	id int, transferReq *request.TransferReq,
	page, limit, lastLimit, length, count int,
	createSqlFn func(offset uint64, limit int) string, ctx context.Context, conn *sqlx.DB, esConn *elasticV6.Client) (err error) {

	table2EsColMap := map[string]string{}
	for _, t := range transferReq.Cols.EsCols {
		table2EsColMap[t.TbCol] = t.Col
	}

	if transferReq.Reset {

		_, err = esConn.DeleteByQuery().
			Index(transferReq.IndexName).
			Type(transferReq.TypeName).
			Body(`
					{
					  "query": {
						"match_all": {}
					  }
					}`).
			Slices(5).
			Do(ctx)
		if err != nil {
			updateDataXListStatus(id, 0, 0, Error, err.Error())
			return err
		}
	}

	limitTmp := limit

	InputC := make(chan map[string]interface{}, 10000)

	goNumControl := make(chan struct{}, transferReq.GoNum)

	realTimeWarehousing := NewRealTimeWarehousingV6(transferReq.EsBufferSize, transferReq.EsFlushInterval, esConn, ctx, id, count)
	realTimeWarehousing.RegularFlushing()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logs.Logger.Sugar().Infof("退出协程")
				return
			case data := <-InputC:
				var err error
				err = realTimeWarehousing.Add(elasticV6.NewBulkIndexRequest().Index(transferReq.IndexName).Type(transferReq.TypeName).Doc(data))
				if err != nil {
					updateDataXListStatus(id, 0, 0, Error, err.Error())
					logs.Logger.Sugar().Errorf("上报失败 重新上报err", err)
				}
			default:

			}
		}
	}()

	go func(page, length, limitTmp, lastLimit int) {
		for ; page <= length; page++ {

			goNumControl <- struct{}{}

			if page == length {
				limitTmp = lastLimit
			}

			go func(page, limit, lastLimit, limitTmp int) {

				select {
				case <-ctx.Done():
					logs.Logger.Sugar().Infof("任务结束", id)
					return
				default:

				}

				sql := createSqlFn(db.CreatePage(page, limit), limitTmp)

				list, err := queryRows(table2EsColMap, conn, sql)
				if err != nil {
					updateDataXListStatus(id, 0, 0, Error, err.Error())
					<-goNumControl
					logs.Logger.Error("err", zap.String("sql", sql), zap.Error(err))
					return
				}

				for _, data := range list {
					InputC <- data
				}

				<-goNumControl
			}(page, limit, lastLimit, limitTmp)
		}
	}(page, length, limitTmp, lastLimit)

	return nil
}

func transferEsV7(
	id int, transferReq *request.TransferReq,
	page, limit, lastLimit, length, count int,
	createSqlFn func(offset uint64, limit int) string, ctx context.Context, conn *sqlx.DB, esConn *elasticV7.Client) (err error) {

	table2EsColMap := map[string]string{}
	for _, t := range transferReq.Cols.EsCols {
		table2EsColMap[t.TbCol] = t.Col
	}

	if transferReq.Reset {

		_, err = esConn.DeleteByQuery().
			Index(transferReq.IndexName).
			Body(`
					{
					  "query": {
						"match_all": {}
					  }
					}`).
			Slices(5).
			Do(ctx)
		if err != nil {
			updateDataXListStatus(id, 0, 0, Error, err.Error())
			return err
		}
	}

	limitTmp := limit

	InputC := make(chan map[string]interface{}, 10000)

	goNumControl := make(chan struct{}, transferReq.GoNum)

	realTimeWarehousing := NewRealTimeWarehousingV7(transferReq.EsBufferSize, transferReq.EsFlushInterval, esConn, ctx, id, count)
	realTimeWarehousing.RegularFlushing()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logs.Logger.Sugar().Infof("退出协程")
				return
			case data := <-InputC:
				var err error
				err = realTimeWarehousing.Add(elasticV7.NewBulkIndexRequest().Index(transferReq.IndexName).Doc(data))
				if err != nil {
					updateDataXListStatus(id, 0, 0, Error, err.Error())
					logs.Logger.Sugar().Errorf("上报失败 重新上报err", err)
				}
			default:

			}
		}
	}()

	go func(page, length, limitTmp, lastLimit int) {
		for ; page <= length; page++ {

			goNumControl <- struct{}{}

			if page == length {
				limitTmp = lastLimit
			}

			go func(page, limit, lastLimit, limitTmp int) {

				select {
				case <-ctx.Done():
					logs.Logger.Sugar().Infof("任务结束", id)
					return
				default:

				}

				sql := createSqlFn(db.CreatePage(page, limit), limitTmp)

				list, err := queryRows(table2EsColMap, conn, sql)
				if err != nil {
					updateDataXListStatus(id, 0, 0, Error, err.Error())
					<-goNumControl
					logs.Logger.Error("err", zap.String("sql", sql), zap.Error(err))
					return
				}

				for _, data := range list {
					InputC <- data
				}

				<-goNumControl
			}(page, limit, lastLimit, limitTmp)
		}
	}(page, length, limitTmp, lastLimit)

	return nil
}

var lock *sync.RWMutex

func init() {
	lock = new(sync.RWMutex)
}

func updateDataXListStatus(id, dbcount, escount int, status, msg string) (err error) {
	lock.Lock()
	defer lock.Unlock()
	if status == Error {
		ts := GetTaskInstance()
		ts.CancelById(id)
	}
	_, err = db.SqlBuilder.Update("datax_transfer_list").
		SetMap(map[string]interface{}{
			"status":    status,
			"error_msg": msg,
			"dbcount":   dbcount,
			"escount":   escount,
			"updated":   time.Now().Format(util.TimeFormat),
		}).Where(db.Eq{"id": id}).RunWith(db.Sqlx).Exec()

	return
}
