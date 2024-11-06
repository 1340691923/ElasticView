package model

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/goccy/go-json"
	"net/http"
	"strings"
)

type EsConnect struct {
	Ip      string `json:"ip" db:"ip"`
	User    string `json:"user" db:"user"`
	Pwd     string `json:"pwd" db:"pwd"`
	Version string `json:"version" db:"version"`
	RootPEM string `json:"rootpem" db:"rootpem"`
	CertPEM string `json:"certpem" db:"certpem"`
	KeyPEM  string `json:"keypem" db:"keypem"`
	Header  string `json:"header" db:"header"`
}

func (this *EsConnect) ToEsSdkCfg(cfg *config.Config, connectId int) *proto.Config {

	header := []vo.HeaderKv{}

	json.Unmarshal([]byte(this.Header), &header)

	httpHeader := http.Header{}

	for _, v := range header {
		httpHeader.Set(v.Key, v.Value)
	}

	return &proto.Config{
		ConnectId:               connectId,
		Cfg:                     cfg,
		Version:                 this.Version,
		Addresses:               strings.Split(this.Ip, ","),
		Username:                this.User,
		Password:                this.Pwd,
		CloudID:                 "",
		APIKey:                  "",
		ServiceToken:            "",
		CertificateFingerprint:  "",
		Header:                  httpHeader,
		CACert:                  nil,
		RetryOnStatus:           nil,
		DisableRetry:            false,
		EnableRetryOnTimeout:    false,
		MaxRetries:              0,
		CompressRequestBody:     false,
		DiscoverNodesOnStart:    false,
		DiscoverNodesInterval:   0,
		EnableMetrics:           false,
		EnableDebugLogger:       false,
		EnableCompatibilityMode: false,
		DisableMetaHeader:       false,
		UseResponseCheckOnly:    false,
		Transport:               nil,
		RootPEM:                 this.RootPEM,
		CertPEM:                 this.CertPEM,
		KeyPEM:                  this.KeyPEM,
	}
}