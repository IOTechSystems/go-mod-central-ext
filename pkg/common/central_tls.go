// Copyright (C) 2022-2024 IOTech Ltd

package common

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

// TLS file settings
const (
	// General settings
	BaseOutputDir       = "/tmp/edgex/secrets"
	CaKeyFileName       = "ca.key"
	CaCertFileName      = "ca.crt"
	OpensslConfFileName = "openssl.conf"
	RsaKySize           = "4096"

	// Redis specific settings
	RedisTlsCertOutputDir = BaseOutputDir + "/redis"
	RedisTlsSecretName    = "redis-tls"
	RedisKeyFileName      = "redis.key"
	RedisCsrFileName      = "redis.csr"
	RedisCertFileName     = "redis.crt"

	// MQTT specific settings
	MqttClientKeyFileName  = "mqtt.key"
	MqttClientCertFileName = "mqtt.crt"
	EnvMessageBusMqttTls   = "EDGECENTRAL_MESSAGEBUS_MQTT_TLS"
)

// CreateRedisTlsConfigFromPEM loads TLS certificates from PEM encoded data and creates Redis TLS config
func CreateRedisTlsConfigFromPEM(certPEMBlock, keyPEMBlock, caPEMBlock []byte) (*tls.Config, errors.EdgeX) {
	var tlsConfig *tls.Config
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return tlsConfig, errors.NewCommonEdgeX(errors.KindServerError, "fail to parse the Redis TLS key pair", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caPEMBlock); !ok {
		return tlsConfig, errors.NewCommonEdgeX(errors.KindServerError, "failed to parse ca to pool", nil)
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
		// skip server side SSL verification, primarily for self-signed certs
		InsecureSkipVerify: true, // nolint:gosec
	}
	return tlsConfig, nil
}
