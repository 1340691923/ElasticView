package escache

import (
	"crypto/tls"
	elasticV7 "github.com/olivere/elastic/v7"
	"net/http"
)

type EsClientV7 struct {
	Client          *elasticV7.Client
	esConnectConfig EsConnect
}

func NewEsClientV7(esConnectConfig *EsConnect) (esClient *elasticV7.Client, err error) {

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	optList := []elasticV7.ClientOptionFunc{
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheck(false),
		elasticV7.SetURL(esConnectConfig.Ip),
		elasticV7.SetHttpClient(httpClient),
	}

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV7.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV7.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}

func NewEsClientV8(esConnectConfig *EsConnect) (esClient *elasticV7.Client, err error) {

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	optList := []elasticV7.ClientOptionFunc{
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheck(false),
		elasticV7.SetURL(esConnectConfig.Ip),
		elasticV7.SetHttpClient(httpClient),
	}

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV7.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV7.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}
