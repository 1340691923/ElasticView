package data_conversion

import (
	"context"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"sync"
)

type DataConversion struct {
	InputC              chan util.Map
	ctx                 context.Context
	wg                  *sync.WaitGroup
	realTimeWarehousing RealTimeWarehousing
	indexName           string
	typeName            string

	taskId int
}

func (this *DataConversion) Data() {
	for {

		select {
		case <-this.ctx.Done():
			logs.Logger.Sugar().Infof("退出协程")
			this.wg.Done()
			return
		case data := <-this.InputC:
			var err error

			err = this.realTimeWarehousing.Add(
				this.realTimeWarehousing.dealwithDataF.fn1(this.indexName, this.typeName, data),
			)

			if err != nil {
				logs.Logger.Sugar().Errorf("上报失败 重新上报err", err)
				continue
			}
		default:
		}
	}
}
