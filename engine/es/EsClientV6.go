package es

import (
	elasticV6 "github.com/olivere/elastic"
)

func NewEsClientV6(esConnectConfig *EsConnect) (esClient *elasticV6.Client, err error) {

	optList := []elasticV6.ClientOptionFunc{
		elasticV6.SetSniff(false),
		elasticV6.SetHealthcheck(false),
	}

	optList = append(optList, elasticV6.SetURL(esConnectConfig.Ip))

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV6.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV6.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}
