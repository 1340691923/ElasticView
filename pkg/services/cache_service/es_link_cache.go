package cache_service

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"sync"
)

// Es连接缓存
type EsCache struct {
	esConnectMap   *sync.Map
	esLinkListMap  *sync.Map
	esConnectLock  *sync.RWMutex
	esLinkListLock *sync.RWMutex
}

var (
	once    sync.Once
	esCache *EsCache
)

// es
func NewEsCache() *EsCache {
	once.Do(func() {
		esCache = &EsCache{
			esConnectMap:   new(sync.Map),
			esLinkListMap:  new(sync.Map),
			esConnectLock:  new(sync.RWMutex),
			esLinkListLock: new(sync.RWMutex),
		}
	})
	return esCache
}

// 新增一个es实例
func (this *EsCache) Set(id int, esClient *model.EsConnect) {
	this.esConnectLock.Lock()
	this.esConnectMap.Store(id, esClient)
	this.esConnectLock.Unlock()
}

// 通过es连接表 的id获取一个保存在内存的id
func (this *EsCache) Get(id int) *model.EsConnect {

	this.esConnectLock.RLock()
	defer this.esConnectLock.RUnlock()
	if v, getConnect := this.esConnectMap.Load(id); getConnect {
		return v.(*model.EsConnect)
	}
	return nil
}

// 新增一个es实例
func (this *EsCache) EsLinkSet(roleId string, opt []vo.EsLinkOpt) {
	this.esLinkListLock.Lock()
	this.esLinkListMap.Store(roleId, opt)
	this.esLinkListLock.Unlock()
}

func (this *EsCache) EsLinkGet(roleId string) (bool, []vo.EsLinkOpt) {
	this.esLinkListLock.RLock()
	defer this.esLinkListLock.RUnlock()

	v, isGetConnect := this.esLinkListMap.Load(roleId);

	if v == nil {
		return isGetConnect, []vo.EsLinkOpt{}
	}

	return isGetConnect, v.([]vo.EsLinkOpt)
}

func (this *EsCache) Truncate() {
	this.esConnectLock.Lock()
	this.esConnectMap.Range(func(key, value interface{}) bool {
		this.esConnectMap.Delete(key)
		return true
	})
	this.esConnectLock.Unlock()
	this.esLinkListLock.Lock()
	this.esLinkListMap.Range(func(key, value interface{}) bool {
		this.esLinkListMap.Delete(key)
		return true
	})
	this.esLinkListLock.Unlock()
}
