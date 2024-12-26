package web_engine

import "sync"

func init() {
	mockMap = NewMockMap()
}

var mockMap *MockMap //RouterConfig

type MockMap struct {
	keys   []string
	values [][]RouterConfig
	m      map[string][]RouterConfig
	lock   *sync.RWMutex
}

func NewMockMap() *MockMap {
	return &MockMap{
		keys:   []string{},
		values: [][]RouterConfig{},
		m:      map[string][]RouterConfig{},
		lock:   new(sync.RWMutex),
	}
}

func (this *MockMap) Store(data string, val RouterConfig) *MockMap {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.m[data]; ok {
		i := 0
		for index, v := range this.keys {
			if v == data {
				i = index
				break
			}
		}
		this.values[i] = append(this.values[i], val)
		this.m[data] = this.values[i]
		return this
	}
	this.keys = append(this.keys, data)
	value := []RouterConfig{val}
	this.values = append(this.values, value)
	this.m[data] = value
	return this
}

func (this *MockMap) Load(key string) ([]RouterConfig, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	v, ok := this.m[key]
	return v, ok
}

func (this *MockMap) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.m == nil {
		return 0
	}
	return len(this.m)
}

func (this *MockMap) Keys() []string {
	return this.keys
}

func (this *MockMap) Values() [][]RouterConfig {
	return this.values
}

func (this *MockMap) Range(fn func(key string, val []RouterConfig)) {
	for _, k := range this.keys {
		v, _ := this.Load(k)
		fn(k, v)
	}
}
