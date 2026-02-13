// Copyright (C) 2023-2026 IOTech Ltd

package dbc

import (
	"fmt"
	"math"
	"strconv"

	edgexCommon "github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	edgexDtos "github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"go.einride.tech/can/pkg/descriptor"
)

func valueType(s *descriptor.Signal) string {
	_, offsetFrac := math.Modf(s.Offset)
	_, scaleFrac := math.Modf(s.Scale)
	if offsetFrac != 0 || scaleFrac != 0 {
		return edgexCommon.ValueTypeFloat64
	}
	if s.IsSigned {
		return edgexCommon.ValueTypeInt64
	} else {
		return edgexCommon.ValueTypeUint64
	}
}

func ConvertDBCtoProfile(data []byte) (Converter[[]edgexDtos.DeviceProfile], errors.EdgeX) {
	converter, err := newDeviceProfileDBC(data)
	if err != nil {
		return nil, err
	}

	err = converter.ConvertToDTO()
	if err != nil {
		return nil, err
	}

	return converter, nil
}

func ConvertDBCtoDevice(data []byte, args map[string]string) (Converter[[]edgexDtos.Device], errors.EdgeX) {
	converter, err := newDeviceDBC(data, args)
	if err != nil {
		return nil, err
	}

	err = converter.ConvertToDTO()
	if err != nil {
		return nil, err
	}

	return converter, nil
}

func getOriginalCanId(canID uint32) string {
	id := canID | messageIDExtendedFlag
	return strconv.FormatUint(uint64(id), 10)
}

func getPGN(canID uint32) string {
	// J1939 PGN bit start from 9, length is 18
	pgn := (canID >> j1939PGNOffset) & j1939PGNMask
	return fmt.Sprintf("%X", pgn)
}
