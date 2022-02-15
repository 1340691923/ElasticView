package es_optimize

import (
	"context"
	"errors"

	"github.com/olivere/elastic"
)

// 关闭索引
type Close struct {
	indexName []string
}

func (this *Close) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this *Close) Do(client *elastic.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}
	_, err = client.CloseIndex(this.indexName[0]).Do(context.Background())
	return
}

func newClose() OptimizeInterface {
	return &Close{}
}
