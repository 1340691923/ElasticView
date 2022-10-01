package es_optimize

import (
	"context"
	"errors"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

// 关闭索引
type Close struct {
	indexName []string
}

func (this *Close) CleanIndexName() {
	this.indexName = []string{}
}

func (this *Close) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this *Close) DoV6(client *elasticV6.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}
	_, err = client.CloseIndex(this.indexName[0]).Do(context.Background())
	return
}

func (this *Close) DoV7(client *elasticV7.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}
	_, err = client.CloseIndex(this.indexName[0]).Do(context.Background())
	return
}

func newClose() OptimizeInterface {
	return &Close{}
}
