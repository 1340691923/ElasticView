package es

import (
	elasticV7 "github.com/olivere/elastic/v7"
)

type EsClientV7 struct {
	Client          *elasticV7.Client
	esConnectConfig EsConnect
}

func NewEsClientV7(esConnectConfig *EsConnect) (esClient *elasticV7.Client, err error) {

	optList := []elasticV7.ClientOptionFunc{
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheck(false),
	}

	optList = append(optList, elasticV7.SetURL(esConnectConfig.Ip))

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV7.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV7.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}
