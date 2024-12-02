// Copyright (C) 2024 IOTech Ltd

package models

// KeyData contains the signing or verification key for the JWT token
type KeyData struct {
	Issuer string
	Type   string
	Key    string
}
