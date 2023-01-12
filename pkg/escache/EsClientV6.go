package escache

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"strings"

	elasticV6 "github.com/olivere/elastic"
)

func NewEsClientV6(esConnectConfig *EsConnect) (esClient *elasticV6.Client, err error) {

	var tlsClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if esConnectConfig.RootPEM != "" {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(esConnectConfig.RootPEM))
		if !ok {
			err = errors.New("failed to parse root certificate")
			return
		}
		tlsClientConfig.RootCAs = roots
	}

	if esConnectConfig.CertPEM != "" && esConnectConfig.KeyPEM != "" {
		var cert tls.Certificate
		cert, err = tls.X509KeyPair([]byte(esConnectConfig.CertPEM), []byte(esConnectConfig.KeyPEM))
		if err != nil {
			return
		}

		tlsClientConfig.Certificates = []tls.Certificate{cert}
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsClientConfig,
		},
	}

	optList := []elasticV6.ClientOptionFunc{
		elasticV6.SetSniff(false),
		elasticV6.SetHealthcheck(false),
		elasticV6.SetURL(strings.Split(esConnectConfig.Ip,",")...),
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
