package es_optimize

import (
	"context"
	"errors"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

// 开启索引
type Open struct {
	indexName []string
}

func (this *Open) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this *Open) DoV6(client *elasticV6.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}

	_, err = client.OpenIndex(this.indexName[0]).Do(context.Background())

	return
}

func (this *Open) DoV7(client *elasticV7.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}

	_, err = client.OpenIndex(this.indexName[0]).Do(context.Background())

	return
}

func newOpen() OptimizeInterface {
	return &Open{}
}
