package es_optimize

import (
	"context"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

// Flush
type Flush struct {
	indexName []string
}
func (this *Flush) CleanIndexName() {
	this.indexName = []string{}
}

func (this *Flush) SetIndexName(indexName string) {
	this.indexName = append(this.indexName, indexName)
}

func (this Flush) DoV6(client *elasticV6.Client) (err error) {
	if len(this.indexName) == 0 {
		_, err = client.Flush().Do(context.Background())
		return
	}
	_, err = client.Flush(this.indexName...).Do(context.Background())
	return
}

func (this Flush) DoV7(client *elasticV7.Client) (err error) {
	if len(this.indexName) == 0 {
		_, err = client.Flush().Do(context.Background())
		return
	}
	_, err = client.Flush(this.indexName...).Do(context.Background())
	return
}

func newFlush() OptimizeInterface {
	return &Flush{}
}
