package es_optimize

import (
	"context"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

//刷新索引
type Refresh struct {
	indexName []string
}

func (this *Refresh) CleanIndexName() {
	this.indexName = []string{}
}

func (this *Refresh) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this Refresh) DoV6(client *elasticV6.Client) (err error) {
	_, err = client.Refresh(this.indexName...).Do(context.Background())
	return
}

func (this Refresh) DoV7(client *elasticV7.Client) (err error) {
	_, err = client.Refresh(this.indexName...).Do(context.Background())
	return
}

func newRefresh() OptimizeInterface {
	return &Refresh{}
}
