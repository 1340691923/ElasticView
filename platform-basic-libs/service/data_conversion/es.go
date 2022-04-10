package data_conversion

import (
	"context"
	"encoding/json"
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

type dealwithData interface {
	fn1(indexName string, typeName string, data util.Map) interface{}
	fn2(data []interface{}) error
}

type Es6 struct {
}

func (this *Es6) fn1(indexName string, typeName string, data util.Map) interface{} {
	return elasticV6.NewBulkIndexRequest().Index(indexName).Type(typeName).Doc(data)
}

func (this *Es6) fn2(data []interface{}, connect *es.EsConnect) error {

	esClinet, err := es.NewEsClientV6(connect)

	if err != nil {
		return err
	}
	bulkRequest := esClinet.Bulk()
	for _, buffer := range data {
		bulkRequest.Add(buffer.(elasticV6.BulkableRequest))
	}
	res, err := bulkRequest.Do(context.Background())

	switch {
	case res.Errors:
		resStr, _ := json.Marshal(res)
		logs.Logger.Sugar().Errorf("res", string(resStr))
	case err != nil:
		logs.Logger.Sugar().Errorf("出现错误", err)
	default:

	}

	return nil
}

type Es7 struct {
}

func (this *Es7) fn1(indexName string, typeName string, data util.Map) interface{} {
	return elasticV7.NewBulkIndexRequest().Index(indexName).Doc(data)
}

func (this *Es7) fn2(data []interface{}, connect *es.EsConnect) error {

	esClinet, err := es.NewEsClientV7(connect)

	if err != nil {
		return err
	}
	bulkRequest := esClinet.Bulk()
	for _, buffer := range data {
		bulkRequest.Add(buffer.(elasticV7.BulkableRequest))
	}
	res, err := bulkRequest.Do(context.Background())

	switch {
	case res.Errors:
		resStr, _ := json.Marshal(res)
		logs.Logger.Sugar().Errorf("res", string(resStr))
	case err != nil:
		logs.Logger.Sugar().Errorf("出现错误", err)
	default:

	}

	return nil
}
