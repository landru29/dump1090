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

// Message is an ADSB message.
type Message struct {
	TransponderCapability byte
	AircraftAddress       uint32
	TypeCode              TypeCode
	Message               []byte
	ModeS                 Frame
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
func (m *Message) Unmarshal(data Frame) error {
	m.ModeS = data

	length := len(data.Data)
	if length < 10 {
		return ErrWrongMessageSize
	}

	if m.ModeS.DownlinkFormat != DownlinkFormatExtendedSquitter && m.ModeS.DownlinkFormat != DownlinkFormatExtendedSquitterNonTransponder {
		return ErrUnsupportedDownlinkFormat
	}

	m.TransponderCapability = data.Raw[0] & 0x07

	m.AircraftAddress = (uint32(data.Data[0]) << 16) + (uint32(data.Data[1]) << 8) + uint32(data.Data[2])

	m.TypeCode = TypeCode((data.Data[3] & 0xf8) >> 3)

	m.Message = data.Data[3:]

	return m.checksumSquitter()
}

// Type is the type of message.
func (m Message) Type() MessageType {
	switch {
	case m.TypeCode > 0 && m.TypeCode < 5:
		// Aircraft Identification.
		return MessageTypeAircraftIdentification

	case m.TypeCode > 4 && m.TypeCode < 9:
		// Surface position.
		return MessageTypeSurfacePosition

	case m.TypeCode > 8 && m.TypeCode < 19:
		// Airborne position (w/Baro Altitude).
		return MessageTypeAirbornePositionBaroAltitude

	case m.TypeCode == 19:
		// Airborne velocities.
		return MessageTypeAirborneVelocities

	case m.TypeCode > 19 && m.TypeCode < 23:
		// Airborne position (w/GNSS Height).
		return MessageTypeAirbornePositionGnssHeight

	case m.TypeCode == 28:
		// Aircraft status.
		return MessageTypeAircraftStatus

	case m.TypeCode == 29:
		// Target state and status information.
		return MessageTypeTargetStateAndStatusInformation

	case m.TypeCode == 31:
		// Aircraft operation status.
		return MessageTypeAircraftOperationStatus

	default:
		return MessageTypeUnsupported
	}
}
