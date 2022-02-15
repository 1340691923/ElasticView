package es_optimize

import (
	"context"
	"errors"
	"log"

	"github.com/olivere/elastic"
)

// 开启索引
type Open struct {
	indexName []string
}

func (this *Open) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this *Open) Do(client *elastic.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}
	log.Println("open")

	_, err = client.OpenIndex(this.indexName[0]).Do(context.Background())

	return
}

func newOpen() OptimizeInterface {
	return &Open{}
}
