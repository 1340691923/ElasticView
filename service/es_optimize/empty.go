package es_optimize

import (
	"context"
	"errors"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

// 关闭索引
type Empty struct {
	indexName []string
}

func (this *Empty) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this *Empty) DoV6(client *elasticV6.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}
	_, err = client.DeleteByQuery(this.indexName[0]).Body(`
{
  "query": { 
    "match_all": {
    }
  }
}`).Do(context.Background())
	return
	return
}

func (this *Empty) DoV7(client *elasticV7.Client) (err error) {
	if len(this.indexName) == 0 {
		return errors.New("索引名不能为空")
	}
	_, err = client.DeleteByQuery(this.indexName[0]).Body(`
{
  "query": { 
    "match_all": {
    }
  }
}`).Do(context.Background())
	return
}

func newEmpty() OptimizeInterface {
	return &Empty{}
}
