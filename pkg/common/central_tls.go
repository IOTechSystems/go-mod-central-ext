// Copyright (C) 2022-2024 IOTech Ltd

package common

// TLS file settings
const (
	// General settings
	BaseOutputDir       = "/tmp/edgex/secrets"
	CaKeyFileName       = "ca.key"
	CaCertFileName      = "ca.crt"
	OpensslConfFileName = "openssl.conf"
	RsaKySize           = "4096"

	// MQTT specific settings
	MqttClientKeyFileName  = "mqtt.key"
	MqttClientCertFileName = "mqtt.crt"
	EnvMessageBusMqttTls   = "EDGECENTRAL_MESSAGEBUS_MQTT_TLS"
)
