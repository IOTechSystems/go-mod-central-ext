// Copyright (C) 2023 IOTech Ltd

package dbc

const (
	ServiceName = "ServiceName"

	Canbus          = "CANbus"
	J1939           = "J1939"
	Network         = "Network"
	Standard        = "Standard"
	ID              = "ID"
	DataSize        = "DataSize"
	Sender          = "Sender"
	PGN             = "PGN"
	CommType        = "CommType"
	CommTypeTCP     = "TCP"
	Port            = "Port"
	NetType         = "NetType"
	NetTypeEthernet = "Ethernet"

	BitStart      = "bitStart"
	BitLen        = "bitLen"
	LittleEndian  = "littleEndian"
	ReceiverNames = "receiverNames"
	MuxSignal     = "muxSignal"
	MuxNum        = "muxNum"
	IsSigned      = "isSigned"

	messageIDExtendedFlag = 0x80000000
	j1939PGNOffset        = 8
	j1939PGNMask          = 0x3FFFF
)
