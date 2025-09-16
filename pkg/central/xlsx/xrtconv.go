// Copyright (C) 2025 IOTech Ltd

package xlsx

import (
	"fmt"
	"strconv"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/xrtmodels"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"

	"github.com/spf13/cast"
)

// toXrtProperties converts the protocol properties to specified data type when importing devices by excel file
func toXrtProperties(protocol string, protocolProperties map[string]interface{}) errors.EdgeX {
	intProperties, floatProperties, boolProperties := xrtmodels.PropertyConversionList(protocol)

	for _, p := range intProperties {
		propertyValue, ok := protocolProperties[p]
		if ok {
			// convert property value from interface{} to string, then to int
			val, err := strconv.Atoi(cast.ToString(propertyValue))
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("fail to convert %v to int", p), err)
			}
			protocolProperties[p] = val
		}
	}

	for _, p := range floatProperties {
		propertyValue, ok := protocolProperties[p]
		if ok {
			// convert property value from interface{} to string, then to float
			val, err := strconv.ParseFloat(cast.ToString(propertyValue), 64)
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("fail to convert %v to float", p), err)
			}
			protocolProperties[p] = val
		}
	}

	for _, p := range boolProperties {
		propertyValue, ok := protocolProperties[p]
		if ok {
			// convert property value from interface{} to string, then to bool
			val, err := strconv.ParseBool(cast.ToString(propertyValue))
			if err != nil {
				return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("fail to convert %v to bool", p), err)
			}
			protocolProperties[p] = val
		}
	}
	return nil
}
