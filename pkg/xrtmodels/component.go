// Copyright (C) 2023 IOTech Ltd

package xrtmodels

// constants of device connector settings
const (
	Name                                     = "Name"
	EventTopic                               = "EventTopic"
	ReplyTopic                               = "ReplyTopic"
	RequestTopic                             = "RequestTopic"
	TelemetryTopic                           = "TelemetryTopic"
	OPCUAServerRequestTimeout                = "OPCUAServerRequestTimeout"
	OPCUAServerUseTelemetryValues            = "OPCUAServerUseTelemetryValues"
	OPCUAServerStaleTelemetryValueTime       = "OPCUAServerStaleTelemetryValueTime"
	OPCUAServerTopicMiddlewarePrefix         = "OPCUAServerTopicMiddlewarePrefix"
	OPCUAServerUseMiddlewarePrefixRequest    = "OPCUAServerUseMiddlewarePrefixRequest"
	OPCUAServerUseMiddlewarePrefixReply      = "OPCUAServerUseMiddlewarePrefixReply"
	OPCUAServerUseMiddlewarePrefixTelemetry  = "OPCUAServerUseMiddlewarePrefixTelemetry"
	OPCUAServerUseMiddlewarePrefixEvent      = "OPCUAServerUseMiddlewarePrefixEvent"
	OPCUAServerUseMiddlewarePrefixEdgeXEvent = "OPCUAServerUseMiddlewarePrefixEdgeXEvent"
	OPCUAServerEdgeXEventTopicBase           = "OPCUAServerEdgeXEventTopicBase"

	DeviceServiceCategory      = "XRT::DeviceService"
	DeviceServiceRunningStatus = "Running"
)

// constants of env to override the device connector settings
const (
	EnvXRTOPCUAServerRequestTimeout                = "XRT_OPCUA_SERVER_REQUEST_TIMEOUT"
	EnvXRTOPCUAServerUseTelemetryValues            = "XRT_OPCUA_SERVER_USE_TELEMETRY_VALUES"
	EnvXRTOPCUAServerStaleTelemetryValueTime       = "XRT_OPCUA_SERVER_STALE_TELEMETRY_VALUE_TIME"
	EnvXRTOPCUAServerTopicMiddlewarePrefix         = "XRT_OPCUA_SERVER_TOPIC_MIDDLEWARE_PREFIX"
	EnvXRTOPCUAServerUseMiddlewarePrefixRequest    = "XRT_OPCUA_SERVER_USE_MIDDLEWARE_PREFIX_REQUEST"
	EnvXRTOPCUAServerUseMiddlewarePrefixReply      = "XRT_OPCUA_SERVER_USE_MIDDLEWARE_PREFIX_REPLY"
	EnvXRTOPCUAServerUseMiddlewarePrefixTelemetry  = "XRT_OPCUA_SERVER_USE_MIDDLEWARE_PREFIX_TELEMETRY"
	EnvXRTOPCUAServerUseMiddlewarePrefixEvent      = "XRT_OPCUA_SERVER_USE_MIDDLEWARE_PREFIX_EVENT"
	EnvXRTOPCUAServerUseMiddlewarePrefixEdgeXEvent = "XRT_OPCUA_SERVER_USE_MIDDLEWARE_PREFIX_EDGEX_EVENT"
	EnvXRTOPCUAServerEdgeXEventTopicBase           = "XRT_OPCUA_SERVER_EDGEX_EVENT_TOPIC_BASE"
)

type Component struct {
	Category string         `json:"category"`
	Config   map[string]any `json:"config"`
	Name     string         `json:"name"`
	State    string         `json:"state"`
	Type     string         `json:"type"`
}
