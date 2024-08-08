package proto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	estransportV8 "github.com/elastic/elastic-transport-go/v8/elastictransport"
	elasticV6 "github.com/elastic/go-elasticsearch/v6"
	estransportV6 "github.com/elastic/go-elasticsearch/v6/estransport"
	elasticV7 "github.com/elastic/go-elasticsearch/v7"
	estransportV7 "github.com/elastic/go-elasticsearch/v7/estransport"
	elasticV8 "github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"time"
)

type Config struct {
	Version   int
	Addresses []string // A list of Elasticsearch nodes to use.
	Username  string   // Username for HTTP Basic Authentication.
	Password  string   // Password for HTTP Basic Authentication.

	CloudID                string // Endpoint for the Elastic Service (https://elastic.co/cloud).
	APIKey                 string // Base64-encoded token for authorization; if set, overrides username/password and service token.
	ServiceToken           string // Service token for authorization; if set, overrides username/password.
	CertificateFingerprint string // SHA256 hex fingerprint given by Elasticsearch on first launch.

	Header http.Header // Global HTTP request header.

	// PEM-encoded certificate authorities.
	// When set, an empty certificate pool will be created, and the certificates will be appended to it.
	// The option is only valid when the transport is not specified, or when it's http.Transport.
	CACert []byte

	RetryOnStatus        []int // List of status codes for retry. Default: 502, 503, 504.
	DisableRetry         bool  // Default: false.
	EnableRetryOnTimeout bool  // Default: false.
	MaxRetries           int   // Default: 3.

	CompressRequestBody  bool // Default: false.
	DiscoverNodesOnStart bool // Discover nodes when initializing the client. Default: false.

	DiscoverNodesInterval time.Duration // Discover nodes periodically. Default: disabled.

	EnableMetrics           bool // Enable the metrics collection.
	EnableDebugLogger       bool // Enable the debug logging.
	EnableCompatibilityMode bool // Enable sends compatibility header

	DisableMetaHeader    bool // Disable the additional "X-Elastic-Client-Meta" HTTP header.
	UseResponseCheckOnly bool

	RetryBackoff func(attempt int) time.Duration // Optional backoff duration. Default: nil.

	Transport http.RoundTripper // The HTTP transport object.

	RootPEM string
	CertPEM string
	KeyPEM  string
}

func (this Config) ConvertV6(
	logger estransportV6.Logger,
	selector estransportV6.Selector,
	connectionPoolFunc func([]*estransportV6.Connection, estransportV6.Selector) estransportV6.ConnectionPool,
) (elasticV6.Config, error) {

	cfg := elasticV6.Config{
		Addresses:             this.Addresses,
		Username:              this.Username,
		Password:              this.Password,
		CloudID:               this.CloudID,
		APIKey:                this.APIKey,
		Header:                this.Header,
		CACert:                this.CACert,
		RetryOnStatus:         this.RetryOnStatus,
		DisableRetry:          this.DisableRetry,
		EnableRetryOnTimeout:  this.EnableRetryOnTimeout,
		MaxRetries:            this.MaxRetries,
		DiscoverNodesOnStart:  this.DiscoverNodesOnStart,
		DiscoverNodesInterval: this.DiscoverNodesInterval,
		EnableMetrics:         this.EnableMetrics,
		EnableDebugLogger:     this.EnableDebugLogger,
		RetryBackoff:          this.RetryBackoff,
		Transport:             this.Transport,
	}

	if logger != nil {
		cfg.Logger = logger
	}
	if selector != nil {
		cfg.Selector = selector
	}
	if connectionPoolFunc != nil {
		cfg.ConnectionPoolFunc = connectionPoolFunc
	}

	transport, err := this.getTransport()
	if err != nil {
		return cfg, err
	}

	cfg.Transport = transport

	return cfg, nil
}

func (this Config) ConvertV7(
	logger estransportV7.Logger,
	selector estransportV7.Selector,
	connectionPoolFunc func([]*estransportV7.Connection, estransportV7.Selector) estransportV7.ConnectionPool,
) (elasticV7.Config, error) {

	cfg := elasticV7.Config{
		Addresses:               this.Addresses,
		Username:                this.Username,
		Password:                this.Password,
		CloudID:                 this.CloudID,
		APIKey:                  this.APIKey,
		ServiceToken:            this.ServiceToken,
		Header:                  this.Header,
		CACert:                  this.CACert,
		RetryOnStatus:           this.RetryOnStatus,
		DisableRetry:            this.DisableRetry,
		EnableRetryOnTimeout:    this.EnableRetryOnTimeout,
		MaxRetries:              this.MaxRetries,
		DiscoverNodesOnStart:    this.DiscoverNodesOnStart,
		DiscoverNodesInterval:   this.DiscoverNodesInterval,
		EnableMetrics:           this.EnableMetrics,
		EnableDebugLogger:       this.EnableDebugLogger,
		DisableMetaHeader:       this.DisableMetaHeader,
		RetryBackoff:            this.RetryBackoff,
		Transport:               this.Transport,
	}

	if logger != nil {
		cfg.Logger = logger
	}
	if selector != nil {
		cfg.Selector = selector
	}
	if connectionPoolFunc != nil {
		cfg.ConnectionPoolFunc = connectionPoolFunc
	}

	transport, err := this.getTransport()
	if err != nil {
		return cfg, err
	}

	cfg.Transport = transport

	return cfg, nil
}

func (this Config) ConvertV8(
	logger estransportV8.Logger,
	selector estransportV8.Selector,
	connectionPoolFunc func([]*estransportV8.Connection, estransportV8.Selector) estransportV8.ConnectionPool,
) (elasticV8.Config, error) {

	cfg := elasticV8.Config{
		Addresses:               this.Addresses,
		Username:                this.Username,
		Password:                this.Password,
		CloudID:                 this.CloudID,
		APIKey:                  this.APIKey,
		ServiceToken:            this.ServiceToken,
		CertificateFingerprint:  this.CertificateFingerprint,
		Header:                  this.Header,
		CACert:                  this.CACert,
		RetryOnStatus:           this.RetryOnStatus,
		DisableRetry:            this.DisableRetry,
		MaxRetries:              this.MaxRetries,
		CompressRequestBody:     this.CompressRequestBody,
		DiscoverNodesOnStart:    this.DiscoverNodesOnStart,
		DiscoverNodesInterval:   this.DiscoverNodesInterval,
		EnableMetrics:           this.EnableMetrics,
		EnableDebugLogger:       this.EnableDebugLogger,
		EnableCompatibilityMode: this.EnableCompatibilityMode,
		DisableMetaHeader:       this.DisableMetaHeader,
		RetryBackoff:            this.RetryBackoff,
		Transport:               this.Transport,
	}

	if logger != nil {
		cfg.Logger = logger
	}
	if selector != nil {
		cfg.Selector = selector
	}
	if connectionPoolFunc != nil {
		cfg.ConnectionPoolFunc = connectionPoolFunc
	}

	transport, err := this.getTransport()
	if err != nil {
		return cfg, err
	}

	cfg.Transport = transport

	return cfg, nil
}

func (this Config) getTransport() (http.RoundTripper, error) {
	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	if this.RootPEM != "" {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(this.RootPEM))
		if !ok {
			return nil, fmt.Errorf("failed to parse root certificate")
		}
		tlsClientConfig.RootCAs = roots
	}

	if this.CertPEM != "" && this.KeyPEM != "" {
		cert, err := tls.X509KeyPair([]byte(this.CertPEM), []byte(this.KeyPEM))
		if err != nil {
			return nil, err
		}

		tlsClientConfig.Certificates = []tls.Certificate{cert}
	}

	transport := &http.Transport{TLSClientConfig: tlsClientConfig}

	return transport, nil
}
