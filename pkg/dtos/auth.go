// Copyright (C) 2025 IOTech Ltd

package dtos

// AuthGraphQL defines the content for authorizing a GraphQL API request.
type AuthGraphQL struct {
	// Path is an HTTP-like path represented as a regex pattern for ACL control.
	// It must follow the format: /service-endpoint/field-name
	//
	// e.g. /alarms-service/graphql/Alarms
	// "/alarms-service/graphql" is the service endpoint, "Alarms" is the QUERY field name.
	//
	// e.g. /alarms-service/graphql/DisableAlarm
	// "/alarms-service/graphql" is the service endpoint, "DisableAlarm" is the MUTATION field name.
	Path string `json:"path" validate:"required,edgex-dto-none-empty-string"`
	// Method is the GraphQL operation type, which is named this way to be consistent with the one in AccessPolicy.
	Method string `json:"operation" validate:"oneof=QUERY MUTATION SUBSCRIPTION,required"`
}
