package adsbmodel

import (
	"github.com/landru29/dump1090/internal/binary"
	"github.com/landru29/dump1090/internal/errors"
)

// ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
// ┃                                  Mode S                                    ┃
// ┠┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
// ┃                                    112                                     ┃
// ┣━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┫
// ┃ DF  |                      Extended squitter                      | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             83                              |   24   ┃
// ┗━━━━━╈━━━━┯━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╈━━━━━━━━┛
//       ┃ CA | ICAO |                    Message                      ┃
//       ┠┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//       ┃ 3  |  24  |                       56                        ┃
//       ┗━━━━┷━━━━━━╈━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
//                   ┃ TC |                  Payload                   ┃
//                   ┠┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//                   ┃ 5  |                    51                      ┃
//                   ┗━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
//               Aircraft Identification TC=1-4
//                   ┏━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//                   ┃ TC | CA | C1 | C2 | C3 | C4 | C5 | C6 | C7 | C8 ┃
//                   ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//                   ┃ 5  |  3 |  6 |  6 |  6 |  6 |  6 |  6 |  6 |  6 ┃
//                   ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛
//                Surface Position TC=4-9
//                   ┏━━━━┯━━━━━┯━━━━┯━━━━━┯━━━┯━━━┯━━━━━━━━━┯━━━━━━━━━┓
//                   ┃ TC | MOV | S  | TRK | T | F | LAT-CPR | LON-CPR ┃
//                   ┠┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┼┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┨
//                   ┃ 5  |  7  | 1  |  7  | 1 | 1 |    17   |   17    ┃
//                   ┗━━━━┷━━━━━┷━━━━┷━━━━━┷━━━┷━━━┷━━━━━━━━━┷━━━━━━━━━┛
//                Airborn position TC=8-18, 20-23
//                   ┏━━━━┯━━━━┯━━━━━┯━━━━━┯━━━┯━━━┯━━━━━━━━━┯━━━━━━━━━┓
//                   ┃ TC | SS | SAF | ALT | T | F | LAT-CPR | LON-CPR ┃
//                   ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┼┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┨
//                   ┃ 5  |  2 |  1  |  12 | 1 | 1 |    17   |   17    ┃
//                   ┗━━━━┷━━━━┷━━━━━┷━━━━━┷━━━┷━━━┷━━━━━━━━━┷━━━━━━━━━┛
//                Airborn velocity TC=19
//                   ┏━━━━┯━━━━┯━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━━━┯━━━━━┯━━━━┯━━━━━┯━━━━━━┯━━━━━━┓
//                   ┃ TC | ST | IC | IFR | NUC |     | VrSrc | Svr | VR | Res | SDif | dAlt ┃
//                   ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┨
//                   ┃ 5  |  3 |  1 |  1  |  3  |  22 |   1   |  1  | 9  |  2  |  1   |  7   ┃
//                   ┗━━━━┷━━━━┷━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━━━┷━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┛
//                Aircraft status TC=28
//                Target State And Status Information TC=29
//                Aircraft Operation Status TC=31
//                   ┏━━━━┯━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
//                   ┃ TC | ST |                                       ┃
//                   ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//                   ┃ 5  |  2 |                 48                    ┃
//                   ┗━━━━┷━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
//
// ┏━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━┓
// ┃ name | Description                 | bits ┃
// ┣━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━┫
// ┃ DF   | Downlink Format             |   5  ┃
// ┃ CA   | Capability                  |   3  ┃
// ┃ TC   | Type Code                   |   5  ┃
// ┃ ICAO | Aircraft unique identifier  |  24  ┃
// ┗━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━┛
//
//

// DownlinkFormat is the 5 first bits of an ADSB message.
type DownlinkFormat uint8

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

	// ErrUnsupportedFormat is when the mode-s format is not supported.
	ErrUnsupportedFormat errors.Error = "unsupported  format"

	// ErrWrongCRC is when a wrong CRC was encountered.
	ErrWrongCRC errors.Error = "wrong CRC"
)

// ModeS is a ModeS frame.
type ModeS []byte

// DownlinkFormat is the DF.
func (m ModeS) DownlinkFormat() DownlinkFormat {
	return DownlinkFormat((m[0] & 0xf8) >> 3) //nolint: gomnd
}

// ParityInterrogator is the Parity.
func (m ModeS) ParityInterrogator() uint32 {
	length := len(m)

	return (uint32(m[length-3]) << 16) + (uint32(m[length-2]) << 8) + //nolint: gomnd
		uint32(m[length-1])
}

// Message decode the mode-s frame.
func (m ModeS) Message() (Message, error) { //nolint: ireturn
	if m.DownlinkFormat() == DownlinkFormatExtendedSquitter && len(m) == 112/8 {
		squitter := ExtendedSquitter{ModeS: m}

		switch squitter.TypeCode() { //nolint: exhaustive
		case TypeCodeAircraftIdentification:
			return Identification{ExtendedSquitter: squitter}, nil

		case TypeCodeAircraftOperationStatus:
			return OperationStatus{ExtendedSquitter: squitter}, nil

		case TypeCodeSurfacePosition:
			return SurfacePosition{ExtendedSquitter: squitter}, nil

		case TypeCodeAirbornePositionBaroAltitude, TypeCodeAirbornePositionGNSSHeight:
			return AirbornePosition{ExtendedSquitter: squitter}, nil

		case TypeCodeAirborneVelocities:
			return AirborneVelocity{ExtendedSquitter: squitter}, nil

		default:
			return nil, ErrUnsupportedFormat
		}
	}

	return nil, ErrUnsupportedFormat
}

// CheckSum checks the integrity of the message.
func (m ModeS) CheckSum() error {
	remainder := binary.ChecksumSquitter(m[:len(m)-3])

	if remainder != m.ParityInterrogator() {
		return ErrWrongCRC
	}

	return nil
}

// ChecksumSquitter computes and check the checksum.
// func (e ExtendedSquitter) ChecksumSquitter() error {
// 	remainder := binary.ChecksumSquitter(e.ModeS.Raw[:len(e.ModeS.Raw)-3])

// 	if remainder != e.ParityInterrogator {
// 		return ErrWrongCRC
// 	}

// 	return nil
// }

// Message is a received message.
type Message interface {
	Name() string
}
