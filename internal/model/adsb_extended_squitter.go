package model

// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
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
//
//
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

const extendedSquitterName = "extended squitter"

// ExtendedSquitter is an extended squitter message.
type ExtendedSquitter struct {
	ModeS
}

// TransponderCapability is the CA.
func (e ExtendedSquitter) TransponderCapability() TransponderCapability {
	return TransponderCapability(e.ModeS[0] & 0x07) //nolint: gomnd
}

// AircraftAddress implements the Squitter interface.
func (e ExtendedSquitter) AircraftAddress() ICAOAddr {
	return ICAOAddr((uint32(e.ModeS[1]) << 16) + //nolint: gomnd
		(uint32(e.ModeS[2]) << 8) + //nolint: gomnd
		uint32(e.ModeS[3]))
}

// Message is the extended squitter message.
func (e ExtendedSquitter) Message() []byte {
	return e.ModeS[4:]
}

// TypeCode is the TC.
func (e ExtendedSquitter) TypeCode() TypeCode {
	return TypeCode(e.ModeS[4] >> 3).Code() //nolint: gomnd
}

// Decode decodes the mode-s frame.
func (e ExtendedSquitter) Decode() (Message, error) { //nolint: ireturn
	if e.DownlinkFormat() == DownlinkFormatExtendedSquitter && len(e.ModeS) == 112/8 {
		switch e.TypeCode() { //nolint: exhaustive
		case TypeCodeAircraftIdentification:
			return Identification{ExtendedSquitter: e}, nil

		case TypeCodeAircraftOperationStatus:
			return OperationStatus{ExtendedSquitter: e}, nil

		case TypeCodeSurfacePosition:
			return SurfacePosition{ExtendedSquitter: e}, nil

		case TypeCodeAirbornePositionBaroAltitude, TypeCodeAirbornePositionGNSSHeight:
			return AirbornePosition{ExtendedSquitter: e}, nil

		case TypeCodeAirborneVelocities:
			return AirborneVelocity{ExtendedSquitter: e}, nil

		default:
			return nil, ErrUnsupportedFormat
		}
	}

	return nil, ErrUnsupportedFormat
}

// Name implements the Squitter interface.
func (e ExtendedSquitter) Name() string {
	return extendedSquitterName
}
