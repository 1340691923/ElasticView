package escache

import (
	"crypto/tls"
	elasticV6 "github.com/olivere/elastic"
	"net/http"
)

func NewEsClientV6(esConnectConfig *EsConnect) (esClient *elasticV6.Client, err error) {

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	optList := []elasticV6.ClientOptionFunc{
		elasticV6.SetSniff(false),
		elasticV6.SetHealthcheck(false),
		elasticV6.SetURL(esConnectConfig.Ip),
		elasticV6.SetHttpClient(httpClient),
	}

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV6.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV6.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}
