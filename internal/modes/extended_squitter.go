// Package modes is the Mode S.
package modes

import (
	"github.com/landru29/dump1090/internal/errors"
)

const (
	ErrWrongMessageSize errors.Error = "wrong message size"

	ErrWrongMessageType errors.Error = "wrong message type"

	ErrWrongCRC errors.Error = "wrong CRC"

	ErrUnsupportedDownlinkFormat errors.Error = "unsupported downlink format"
)

type TypeCode byte

// ExtendedSquitter is an extended squitter message.
type ExtendedSquitter struct {
	ModeS
	TransponderCapability byte
	AircraftAddress       uint32
	TypeCode              TypeCode
	Message               []byte
}

type MessageType int

const (
	// MessageTypeUnsupported is the unsupported message type.
	MessageTypeUnsupported MessageType = iota
	// MessageTypeAircraftIdentification is the aircraft identification.
	MessageTypeAircraftIdentification
	// MessageTypeSurfacePosition is the surface position.
	MessageTypeSurfacePosition
	// MessageTypeAirbornePositionBaroAltitude is the airborne position (w/Baro Altitude).
	MessageTypeAirbornePositionBaroAltitude
	// MessageTypeAirborneVelocities is the airborne velocities.
	MessageTypeAirborneVelocities
	// MessageTypeAirbornePositionGnssHeight is the airborne position (w/GNSS Height).
	MessageTypeAirbornePositionGnssHeight
	// MessageTypeAircraftStatus is the aircraft status.
	MessageTypeAircraftStatus
	// MessageTypeTargetStateAndStatusInformation is the target state and status information.
	MessageTypeTargetStateAndStatusInformation
	// MessageTypeAircraftOperationStatus is the aircraft operation status.
	MessageTypeAircraftOperationStatus
)

// Unmarshal parses the message.
func (e *ExtendedSquitter) Unmarshal(data []byte) error {
	if err := (&e.ModeS).Unmarshal(data); err != nil {
		return err
	}

	messageData := data[1 : len(data)-3]

	length := len(messageData)
	if length < 10 {
		return ErrWrongMessageSize
	}

	if e.ModeS.DownlinkFormat != DownlinkFormatExtendedSquitter && e.ModeS.DownlinkFormat != DownlinkFormatExtendedSquitterNonTransponder {
		return ErrUnsupportedDownlinkFormat
	}

	e.TransponderCapability = data[0] & 0x07

	e.AircraftAddress = (uint32(messageData[0]) << 16) + (uint32(messageData[1]) << 8) + uint32(messageData[2])

	e.TypeCode = TypeCode((messageData[3] & 0xf8) >> 3)

	e.Message = messageData[3:]

	return e.checksum()
}

// Type is the type of message.
func (e ExtendedSquitter) Type() MessageType {
	switch {
	case e.TypeCode > 0 && e.TypeCode < 5:
		// Aircraft Identification.
		return MessageTypeAircraftIdentification

	case e.TypeCode > 4 && e.TypeCode < 9:
		// Surface position.
		return MessageTypeSurfacePosition

	case e.TypeCode > 8 && e.TypeCode < 19:
		// Airborne position (w/Baro Altitude).
		return MessageTypeAirbornePositionBaroAltitude

	case e.TypeCode == 19:
		// Airborne velocities.
		return MessageTypeAirborneVelocities

	case e.TypeCode > 19 && e.TypeCode < 23:
		// Airborne position (w/GNSS Height).
		return MessageTypeAirbornePositionGnssHeight

	case e.TypeCode == 28:
		// Aircraft status.
		return MessageTypeAircraftStatus

	case e.TypeCode == 29:
		// Target state and status information.
		return MessageTypeTargetStateAndStatusInformation

	case e.TypeCode == 31:
		// Aircraft operation status.
		return MessageTypeAircraftOperationStatus

	default:
		return MessageTypeUnsupported
	}
}
