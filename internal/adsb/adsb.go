// Package adsb is the Automatic Dependent Surveillance-Broadcast.
package adsb

import (
	"github.com/landru29/dump1090/internal/errors"
)

const (
	ErrWrongMessageSize errors.Error = "wrong message size"

	ErrWrongMessageType errors.Error = "wrong message type"

	ErrWrongCRC errors.Error = "wrong CRC"

	ErrUnsupportedDownlinkFormat errors.Error = "unsupported downlink format"
)

// DownlinkFormat is the 5 first bits of an ADSB message.
type DownlinkFormat byte

const (
	// DownlinkFormatShortAirAirSurveillance is Short air-air surveillance (ACAS) => message size: 56 bits.
	DownlinkFormatShortAirAirSurveillance DownlinkFormat = 0
	// DownlinkFormatAltitudeReply is Altitude reply => message size: 56 bits.
	DownlinkFormatAltitudeReply DownlinkFormat = 4
	// DownlinkFormatIdentityReply is Identity reply => message size: 56 bits.
	DownlinkFormatIdentityReply DownlinkFormat = 5
	// DownlinkFormatAllCallReply is All-call reply => message size: 56 bits.
	DownlinkFormatAllCallReply DownlinkFormat = 11
	// DownlinkFormatLongAirAirSurveillance is Long air-air surveillance (ACAS) => message size: 112 bits.
	DownlinkFormatLongAirAirSurveillance DownlinkFormat = 16
	// DownlinkFormatExtendedSquitter is Extended squitter => message size: 112 bits.
	DownlinkFormatExtendedSquitter DownlinkFormat = 17
	// DownlinkFormatExtendedSquitterNonTransponder is Extended squitter, non transponder => message size: 112 bits.
	DownlinkFormatExtendedSquitterNonTransponder DownlinkFormat = 18
	// DownlinkFormatMilitaryExtendedSquitter is Military extended squitter => message size: 112 bits.
	DownlinkFormatMilitaryExtendedSquitter DownlinkFormat = 19
	// DownlinkFormatCommBWithAltitudeReply is Comm-B, with altitude reply => message size: 112 bits.
	DownlinkFormatCommBWithAltitudeReply DownlinkFormat = 20
	// DownlinkFormatCommBWithIdentityReply is Comm-B, with identity reply => message size: 112 bits.
	DownlinkFormatCommBWithIdentityReply DownlinkFormat = 21
	// DownlinkFormatCommDExtendedLengthMessage is Comm-D, extended length message => message size: 112 bits.
	DownlinkFormatCommDExtendedLengthMessage DownlinkFormat = 24
)

type TypeCode byte

// Message is an ADSB message.
type Message struct {
	DownlinkFormat        DownlinkFormat
	TransponderCapability byte
	AircraftAddress       uint32
	TypeCode              TypeCode
	Message               []byte
	ParityInterrogator    uint32
	Raw                   []byte
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
func (m *Message) Unmarshal(data []byte) error {
	length := len(data)
	if length < 7 {
		return ErrWrongMessageSize
	}

	m.DownlinkFormat = DownlinkFormat((data[0] & 0xf8) >> 3)

	if m.DownlinkFormat != 17 && m.DownlinkFormat != 18 {
		return ErrUnsupportedDownlinkFormat
	}

	m.TransponderCapability = data[0] & 0x07

	m.AircraftAddress = (uint32(data[1]) << 16) + (uint32(data[2]) << 8) + uint32(data[3])

	m.TypeCode = TypeCode((data[4] & 0xf8) >> 3)

	m.ParityInterrogator = (uint32(data[length-3]) << 16) + (uint32(data[length-2]) << 8) + uint32(data[length-1])

	m.Message = data[4 : length-3]

	m.Raw = data[:length-3]

	return m.checksum()
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
