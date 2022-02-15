//ES引擎层
package es

import (
	"sync"
)

// Es连接缓存
type EsCache struct {
	esConnectMap map[int]EsClient
}

var (
	once    sync.Once
	esCache *EsCache
)

// es
func NewEsCache() *EsCache {
	once.Do(func() {
		esCache = &EsCache{esConnectMap: map[int]EsClient{}}
	})
	return esCache
}

// 新增一个es实例
func (this *EsCache) Set(id int, esClient EsClient) {
	this.esConnectMap[id] = esClient
}

// 通过es连接表 的id获取一个保存在内存的id
func (this *EsCache) Get(id int) EsClient {
	if _, getConnect := this.esConnectMap[id]; getConnect {
		return this.esConnectMap[id]
	}
	return nil
}

// 通过id删除一个内存中的es实例
func (this *EsCache) Rem(id int) {
	if _, getConnect := this.esConnectMap[id]; getConnect {
		delete(this.esConnectMap, id)
	}
}
