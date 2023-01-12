package escache

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"strings"

	elasticV7 "github.com/olivere/elastic/v7"
)

type EsClientV7 struct {
	Client          *elasticV7.Client
	esConnectConfig EsConnect
}

func NewEsClientV7(esConnectConfig *EsConnect) (esClient *elasticV7.Client, err error) {

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

	optList := []elasticV7.ClientOptionFunc{
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheck(false),
		elasticV7.SetURL(strings.Split(esConnectConfig.Ip,",")...),
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

	optList := []elasticV7.ClientOptionFunc{
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheck(false),
		elasticV7.SetURL(strings.Split(esConnectConfig.Ip,",")...),
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
