package data_conversion

import (
	"context"
	"encoding/json"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	elasticV7 "github.com/olivere/elastic/v7"
	"sync"
	"time"
)

type RealTimeWarehousingV7 struct {
	buffer        []*elasticV7.BulkIndexRequest
	bufferMutex   *sync.RWMutex
	conn          *elasticV7.Client
	taskId        int
	batchSize     int
	flushInterval int
	ctx           context.Context
	expectLen     int
	completeLen   int
}

func NewRealTimeWarehousingV7(batchSize, flushInterval int, conn *elasticV7.Client, ctx context.Context, taskId int, expectLen int) *RealTimeWarehousingV7 {
	logs.Logger.Sugar().Infof("NewRealTimeWarehousing batchSize:%v,flushInterval:%v", batchSize, flushInterval)

	return &RealTimeWarehousingV7{
		buffer:        make([]*elasticV7.BulkIndexRequest, 0, batchSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		conn:          conn,
		ctx:           ctx,
		taskId:        taskId,
		expectLen:     expectLen,
	}
}

func (this *RealTimeWarehousingV7) Flush() (err error) {

	this.bufferMutex.Lock()

	select {
	case <-this.ctx.Done():
		return
	default:

	}

	if len(this.buffer) > 0 {

		bulkRequest := this.conn.Bulk()

		for _, buffer := range this.buffer {
			bulkRequest.Add(buffer)
			this.completeLen++
		}
		res, err := bulkRequest.Do(this.ctx)

		if err != nil {
			updateDataXListStatus(this.taskId, this.expectLen, this.completeLen, Error, err.Error())
			logs.Logger.Sugar().Infof("插入失败！", err)
		} else {
			if res.Errors {
				resStr, _ := json.Marshal(res.Items[0])
				updateDataXListStatus(this.taskId, this.expectLen, this.completeLen, Error, string(resStr))
				logs.Logger.Sugar().Errorf("res", string(resStr))
			} else {
				if this.expectLen == this.completeLen {
					ts := GetTaskInstance()
					ts.CancelById(this.taskId)
					updateDataXListStatus(this.taskId, this.expectLen, this.completeLen, Success, "数据已全部导入完毕！")
					logs.Logger.Sugar().Infof("所有数据都插入完成！")
				} else {

					updateDataXListStatus(this.taskId, this.expectLen, this.completeLen, Running, "正在导入...")
				}
				logs.Logger.Sugar().Infof("插入成功，继续插入！")
			}
		}
		this.buffer = make([]*elasticV7.BulkIndexRequest, 0, this.batchSize)

	}
	this.bufferMutex.Unlock()
	return nil
}

func (this *RealTimeWarehousingV7) Add(data *elasticV7.BulkIndexRequest) (err error) {
	this.bufferMutex.Lock()
	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()

	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}
	return nil
}

func (this *RealTimeWarehousingV7) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *RealTimeWarehousingV7) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *RealTimeWarehousingV7) RegularFlushing() {
	go func() {
		ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-this.ctx.Done():
				return
			case <-ticker.C:
				if err := this.Flush(); err != nil {
					logs.Logger.Sugar().Errorf("err", err)
				}
			default:

			}

		}
	}()
}
