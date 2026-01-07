// Copyright (C) 2026 IOTech Ltd

package protobuf

import (
	"encoding/json"
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"

	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
)

// ConvertEventToProtobuf converts dtos.Event to protobuf Event
func ConvertEventToProtobuf(event dtos.Event) (*Event, error) {
	pbEvent := &Event{
		ApiVersion:  event.ApiVersion,
		Id:          event.Id,
		DeviceName:  event.DeviceName,
		ProfileName: event.ProfileName,
		SourceName:  event.SourceName,
		Origin:      event.Origin,
	}

	pbEvent.Readings = make([]*Reading, len(event.Readings))
	for i, reading := range event.Readings {
		pbReading, err := convertReadingToProtobuf(reading)
		if err != nil {
			return nil, fmt.Errorf("failed to convert reading %d: %w", i, err)
		}
		pbEvent.Readings[i] = pbReading
	}

	if len(event.Tags) > 0 {
		jsonBytes, err := json.Marshal(event.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		pbEvent.Tags = jsonBytes
	}

	return pbEvent, nil
}

// convertReadingToProtobuf converts dtos.BaseReading to protobuf Reading
func convertReadingToProtobuf(reading dtos.BaseReading) (*Reading, error) {
	pbReading := &Reading{
		Id:           reading.Id,
		Origin:       reading.Origin,
		DeviceName:   reading.DeviceName,
		ResourceName: reading.ResourceName,
		ProfileName:  reading.ProfileName,
		ValueType:    reading.ValueType,
		Units:        reading.Units,
		IsNull:       reading.IsNull(),
	}

	switch reading.ValueType {
	case common.ValueTypeBinary:
		pbReading.BinaryValue = reading.BinaryValue
		pbReading.MediaType = reading.MediaType
	case common.ValueTypeObject, common.ValueTypeObjectArray:
		if reading.ObjectValue != nil {
			// Encode object value as MessagePack bytes to preserve type information
			msgpackBytes, err := msgpack.Marshal(reading.ObjectValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal object value: %w", err)
			}
			pbReading.ObjectValue = msgpackBytes
		}
	default:
		// SimpleReading or NumericReading
		if reading.NumericValue != nil {
			msgpackBytes, err := msgpack.Marshal(reading.NumericValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal numeric value: %w", err)
			}
			pbReading.NumericValue = msgpackBytes
		} else {
			pbReading.Value = reading.Value
		}
	}

	if len(reading.Tags) > 0 {
		jsonBytes, err := json.Marshal(reading.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal reading tags: %w", err)
		}
		pbReading.Tags = jsonBytes
	}

	return pbReading, nil
}

// DecodeProtobufToEvent decodes protobuf bytes back to EdgeX Event
func DecodeProtobufToEvent(data []byte) (*dtos.Event, error) {
	var pbEvent Event
	if err := proto.Unmarshal(data, &pbEvent); err != nil {
		return nil, err
	}

	event := &dtos.Event{
		Id:          pbEvent.GetId(),
		DeviceName:  pbEvent.GetDeviceName(),
		ProfileName: pbEvent.GetProfileName(),
		SourceName:  pbEvent.GetSourceName(),
		Origin:      pbEvent.GetOrigin(),
	}
	event.ApiVersion = pbEvent.GetApiVersion()

	event.Readings = make([]dtos.BaseReading, len(pbEvent.GetReadings()))
	for i, pbReading := range pbEvent.GetReadings() {
		reading, err := convertProtobufToReading(pbReading)
		if err != nil {
			return nil, fmt.Errorf("failed to convert reading %d: %w", i, err)
		}
		event.Readings[i] = reading
	}

	if len(pbEvent.GetTags()) > 0 {
		if err := json.Unmarshal(pbEvent.GetTags(), &event.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}

	return event, nil
}

// convertProtobufToReading converts protobuf Reading to dtos.BaseReading
func convertProtobufToReading(pbReading *Reading) (dtos.BaseReading, error) {
	reading := dtos.BaseReading{
		Id:           pbReading.GetId(),
		Origin:       pbReading.GetOrigin(),
		DeviceName:   pbReading.GetDeviceName(),
		ResourceName: pbReading.GetResourceName(),
		ProfileName:  pbReading.GetProfileName(),
		ValueType:    pbReading.GetValueType(),
		Units:        pbReading.GetUnits(),
	}

	switch pbReading.GetValueType() {
	case common.ValueTypeBinary:
		reading.BinaryReading = dtos.BinaryReading{
			BinaryValue: pbReading.GetBinaryValue(),
			MediaType:   pbReading.GetMediaType(),
		}
	case common.ValueTypeObject, common.ValueTypeObjectArray:
		if len(pbReading.GetObjectValue()) > 0 {
			var objectValue any
			if err := msgpack.Unmarshal(pbReading.GetObjectValue(), &objectValue); err != nil {
				return reading, fmt.Errorf("failed to unmarshal object value: %w", err)
			}
			reading.ObjectReading = dtos.ObjectReading{
				ObjectValue: objectValue,
			}
		}
	default:
		// SimpleReading or NumericReading
		if len(pbReading.GetNumericValue()) > 0 {
			var numericValue any
			if err := msgpack.Unmarshal(pbReading.GetNumericValue(), &numericValue); err != nil {
				return reading, fmt.Errorf("failed to unmarshal numeric value: %w", err)
			}
			// MessagePack preserves the original value type, so we can use it directly
			reading.NumericReading = dtos.NumericReading{
				NumericValue: numericValue,
			}
		} else if pbReading.GetValue() != "" {
			reading.SimpleReading = dtos.SimpleReading{
				Value: pbReading.GetValue(),
			}
		}
	}

	if len(pbReading.GetTags()) > 0 {
		if err := json.Unmarshal(pbReading.GetTags(), &reading.Tags); err != nil {
			return reading, fmt.Errorf("failed to unmarshal reading tags: %w", err)
		}
	}

	return reading, nil
}
