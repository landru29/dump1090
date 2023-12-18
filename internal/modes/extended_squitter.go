// Package modes is the Mode S.
package modes

import (
	"github.com/landru29/dump1090/internal/errors"
)

const (
	// ErrWrongMessageSize is when the message size is not coherent.
	ErrWrongMessageSize errors.Error = "wrong message size"

	// ErrWrongMessageType is when the message type is not recognized.
	ErrWrongMessageType errors.Error = "wrong message type"

	// ErrWrongCRC is when a wrong CRC was encountered.
	ErrWrongCRC errors.Error = "wrong CRC"

	// ErrUnsupportedDownlinkFormat is when the downloink format is not supported.
	ErrUnsupportedDownlinkFormat errors.Error = "unsupported downlink format"
)

// AircraftAddress is the ICAO address of an aircraft.
type AircraftAddress uint32

// TypeCode is the type of the exetended squitter.
type TypeCode byte

// ExtendedSquitter is an extended squitter message.
type ExtendedSquitter struct {
	ModeS
	TransponderCapability byte
	AircraftAddress       AircraftAddress
	TypeCode              TypeCode
	Message               []byte
	Type                  MessageType
}

// MessageType is the type of the message in the extended squitter.
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
func (e *ExtendedSquitter) Unmarshal(data []byte) error { //nolint: cyclop
	if err := (&e.ModeS).Unmarshal(data); err != nil {
		return err
	}

	messageData := data[1 : len(data)-3]

	length := len(messageData)
	if length < 10 { //nolint: gomnd
		return ErrWrongMessageSize
	}

	if e.ModeS.DownlinkFormat != DownlinkFormatExtendedSquitter &&
		e.ModeS.DownlinkFormat != DownlinkFormatExtendedSquitterNonTransponder {
		return ErrUnsupportedDownlinkFormat
	}

	e.TransponderCapability = data[0] & 0x07 //nolint: gomnd

	e.AircraftAddress = AircraftAddress((uint32(messageData[0]) << 16) + //nolint: gomnd
		(uint32(messageData[1]) << 8) + //nolint: gomnd
		uint32(messageData[2]))

	e.TypeCode = TypeCode((messageData[3] & 0xf8) >> 3) //nolint: gomnd

	e.Message = messageData[3:]

	switch {
	case e.TypeCode > 0 && e.TypeCode < 5:
		// Aircraft Identification.
		e.Type = MessageTypeAircraftIdentification

	case e.TypeCode > 4 && e.TypeCode < 9:
		// Surface position.
		e.Type = MessageTypeSurfacePosition

	case e.TypeCode > 8 && e.TypeCode < 19:
		// Airborne position (w/Baro Altitude).
		e.Type = MessageTypeAirbornePositionBaroAltitude

	case e.TypeCode == 19: //nolint: gomnd
		// Airborne velocities.
		e.Type = MessageTypeAirborneVelocities

	case e.TypeCode > 19 && e.TypeCode < 23:
		// Airborne position (w/GNSS Height).
		e.Type = MessageTypeAirbornePositionGnssHeight

	case e.TypeCode == 28: //nolint: gomnd
		// Aircraft status.
		e.Type = MessageTypeAircraftStatus

	case e.TypeCode == 29: //nolint: gomnd
		// Target state and status information.
		e.Type = MessageTypeTargetStateAndStatusInformation

	case e.TypeCode == 31: //nolint: gomnd
		// Aircraft operation status.
		e.Type = MessageTypeAircraftOperationStatus

	default:
		e.Type = MessageTypeUnsupported
	}

	return e.ChecksumSquitter()
}
