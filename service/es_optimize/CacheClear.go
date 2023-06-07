// 索引基本操作层
package es_optimize

import (
	"context"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

// 清除缓存
type CacheClear struct {
	indexName []string
}

func (this *CacheClear) CleanIndexName() {
	this.indexName = []string{}
}

func (this *CacheClear) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this CacheClear) DoV6(client *elasticV6.Client) (err error) {
	_, err = client.ClearCache(this.indexName...).Do(context.Background())
	return
}

func (this CacheClear) DoV7(client *elasticV7.Client) (err error) {
	_, err = client.ClearCache(this.indexName...).Do(context.Background())
	return
}

func newCacheClear() OptimizeInterface {
	return &CacheClear{}
}
