package model

import "github.com/landru29/dump1090/internal/binary"

// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
// ┃ DF  |                      Extended squitter                      | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             83                              |   24   ┃
// ┗━━━━━╈━━━━┯━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╈━━━━━━━━┛
//       ┃ CA | ICAO |                    Message                      ┃
//       ┠┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//       ┃ 3  |  24  |                       56                        ┃
//       ┗━━━━┷━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

// TypeCode is the type of the exetended squitter.
type TypeCode byte

// ExtendedSquitter is an extended squitter message.
type ExtendedSquitter struct {
	ModeS
	TransponderCapability byte
	AircraftAddress       ICAOAddr
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

// UnmarshalModeS parses the message.
func (e *ExtendedSquitter) UnmarshalModeS(data []byte) error { //nolint: cyclop
	if err := (&e.ModeS).UnmarshalModeS(data); err != nil {
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

	e.AircraftAddress = ICAOAddr((uint32(messageData[0]) << 16) + //nolint: gomnd
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

// ChecksumSquitter computes and check the checksum.
func (e ExtendedSquitter) ChecksumSquitter() error {
	remainder := binary.ChecksumSquitter(e.ModeS.Raw[:len(e.ModeS.Raw)-3])

	if remainder != e.ParityInterrogator {
		return ErrWrongCRC
	}

	return nil
}
