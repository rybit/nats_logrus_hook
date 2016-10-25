package nhook

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/nats-io/nats"
)

// NatsConfig represents the minimum entries that are needed to connect to Nats over TLS
type NatsConfig struct {
	CAFiles  []string `json:"ca_files"`
	KeyFile  string   `json:"key_file"`
	CertFile string   `json:"cert_file"`
	Servers  []string `json:"servers"`
}

// ServerString will build the proper string for nats connect
func (config *NatsConfig) ServerString() string {
	return strings.Join(config.Servers, ",")
}

// LogFields will return all the fields relevant to this config
func (config *NatsConfig) LogFields() logrus.Fields {
	return logrus.Fields{
		"servers":   config.Servers,
		"ca_files":  config.CAFiles,
		"key_file":  config.KeyFile,
		"cert_file": config.CertFile,
	}
}

// TLSConfig will load the TLS certificate
func (config *NatsConfig) TLSConfig() (*tls.Config, error) {
	pool := x509.NewCertPool()
	for _, caFile := range config.CAFiles {
		caData, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		if !pool.AppendCertsFromPEM(caData) {
			return nil, fmt.Errorf("Failed to add CA cert at %s", caFile)
		}
	}

	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// ConnectToNats will do a TLS connection to the nats servers specified
func ConnectToNats(config *NatsConfig) (*nats.Conn, error) {
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, err
	}

	return nats.Connect(config.ServerString(), nats.Secure(tlsConfig))
}

// ConnectToNatsWithError will do a TLS connection to the nats servers specified
func ConnectToNatsWithError(config *NatsConfig, eHandler nats.ErrHandler) (*nats.Conn, error) {
	tlsConfig, err := config.TLSConfig()
	if err != nil {
		return nil, err
	}

	if eHandler != nil {
		return nats.Connect(config.ServerString(), nats.Secure(tlsConfig), nats.ErrorHandler(eHandler))
	} else {
		return nats.Connect(config.ServerString(), nats.Secure(tlsConfig))
	}
}
