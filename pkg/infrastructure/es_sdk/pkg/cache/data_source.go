package cache

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	dataSourceCache *sync.Map
	once            sync.Once
)

func init() {
	once.Do(func() {
		dataSourceCache = new(sync.Map)
		go func() {
			ticker := time.NewTicker(5 * time.Minute)

			// 循环等待定时器触发
			for {
				select {
				case <-ticker.C:
					//进行检测连接是否存在ping异常，存在则移除
					CleanDataSourceCache(false)
					runtime.GC()
				}
			}
		}()
	})
}

func SaveDataSourceCache(connId interface{}, dataSource pkg.ClientInterface) {
	dataSourceCache.Store(connId, dataSource)
}

func GetDataSourceCache(connId interface{}) (pkg.ClientInterface, bool) {
	ds, ok := dataSourceCache.Load(connId)

	log.Println("获取连接", connId, ok)

	if ok {
		return ds.(pkg.ClientInterface), ok
	}

	return nil, ok
}

func DeleteDataSourceCache(connId interface{}) {
	dataSourceCache.Delete(connId)
}

func CleanDataSourceCache(isClearAll bool) {
	dataSourceCache.Range(func(connId, dataSource any) bool {

		if isClearAll {
			DeleteDataSourceCache(connId)
			return true
		}

		ds, ok := GetDataSourceCache(connId)

		if !ok {
			return true
		}

		res, err := ds.Ping(context.Background())

		if err != nil {
			DeleteDataSourceCache(connId)
			return true
		}

		if res.StatusErr() != nil {
			DeleteDataSourceCache(connId)
			return true
		}

		return true
	})
}
