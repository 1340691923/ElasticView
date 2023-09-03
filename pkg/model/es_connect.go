package model

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"strings"
)

type EsConnect struct {
	Ip      string `json:"ip" db:"ip"`
	User    string `json:"user" db:"user"`
	Pwd     string `json:"pwd" db:"pwd"`
	Version int    `json:"version" db:"version"`
	RootPEM string `json:"rootpem" db:"rootpem"`
	CertPEM string `json:"certpem" db:"certpem"`
	KeyPEM  string `json:"keypem" db:"keypem"`
}

func (this *EsConnect) ToEsSdkCfg() proto.Config {
	return proto.Config{
		Version:                 this.Version,
		Addresses:               strings.Split(this.Ip, ","),
		Username:                this.User,
		Password:                this.Pwd,
		CloudID:                 "",
		APIKey:                  "",
		ServiceToken:            "",
		CertificateFingerprint:  "",
		Header:                  nil,
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
		RetryBackoff:            nil,
		Transport:               nil,
		RootPEM:                 this.RootPEM,
		CertPEM:                 this.CertPEM,
		KeyPEM:                  this.KeyPEM,
	}
}
