//ES引擎层
package escache

import (
	"sync"
)

// Es连接缓存
type EsCache struct {
	esConnectMap *sync.Map
}

var (
	once    sync.Once
	esCache *EsCache
)

// es
func NewEsCache() *EsCache {
	once.Do(func() {
		esCache = &EsCache{esConnectMap: new(sync.Map)}
	})
	return esCache
}

// 新增一个es实例
func (this *EsCache) Set(id int, esClient *EsConnect) {
	this.esConnectMap.Store(id, esClient)
}

// 通过es连接表 的id获取一个保存在内存的id
func (this *EsCache) Get(id int) *EsConnect {
	if v, getConnect := this.esConnectMap.Load(id); getConnect {
		return v.(*EsConnect)
	}
	return nil
}

// 通过id删除一个内存中的es实例
func (this *EsCache) Rem(id int) {
	this.esConnectMap.Range(func(key, value interface{}) bool {
		if id == key {
			this.esConnectMap.Delete(key)
		}
		return true
	})
}
