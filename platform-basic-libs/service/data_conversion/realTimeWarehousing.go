package data_conversion

import (
	"github.com/1340691923/ElasticView/engine/logs"
	"sync"
	"time"
)

type RealTimeWarehousing struct {
	buffer        []interface{}
	dealwithDataF dealwithData

	bufferMutex   *sync.RWMutex
	batchSize     int
	flushInterval int
}

func NewRealTimeWarehousing(batchSize, flushInterval int, esClient interface{}, dealwithDataF dealwithData) *RealTimeWarehousing {
	logs.Logger.Sugar().Infof("NewRealTimeWarehousing batchSize:%v,flushInterval:%v", batchSize, flushInterval)
	return &RealTimeWarehousing{
		buffer:        make([]interface{}, 0, batchSize),
		bufferMutex:   new(sync.RWMutex),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		dealwithDataF: dealwithDataF,
	}
}

func (this *RealTimeWarehousing) Flush() (err error) {
	this.bufferMutex.Lock()

	logs.Logger.Sugar().Infof("RealTimeWarehousing入库")
	if len(this.buffer) > 0 {

		err := this.dealwithDataF.fn2(this.buffer)

		if err != nil {
			this.buffer = make([]interface{}, 0, this.batchSize)
		}

	}
	this.bufferMutex.Unlock()
	return nil
}

func (this *RealTimeWarehousing) Add(data interface{}) (err error) {
	this.bufferMutex.Lock()
	this.buffer = append(this.buffer, data)
	this.bufferMutex.Unlock()
	logs.Logger.Sugar().Infof("this.getBufferLength", this.getBufferLength(), this.batchSize)
	if this.getBufferLength() >= this.batchSize {
		err := this.Flush()
		return err
	}

	return nil
}

func (this *RealTimeWarehousing) getBufferLength() int {
	this.bufferMutex.RLock()
	defer this.bufferMutex.RUnlock()
	return len(this.buffer)
}

func (this *RealTimeWarehousing) FlushAll() error {
	for this.getBufferLength() > 0 {
		if err := this.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (this *RealTimeWarehousing) RegularFlushing() {
	go func() {
		ticker := time.NewTicker(time.Duration(this.flushInterval) * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			if err := this.Flush(); err != nil {
				logs.Logger.Sugar().Errorf("err", err)
			}
		}
	}()
}
