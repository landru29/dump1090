package adsbmodel

// ┏━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┓
// ┃ DF  |                      Extended squitter                      | Parity ┃
// ┠┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┨
// ┃  5  |                             83                              |   24   ┃
// ┗━━━━━╈━━━━┯━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╈━━━━━━━━┛
//       ┃ CA | ICAO |                    Message                      ┃
//       ┠┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//       ┃ 3  |  24  |                       56                        ┃
//       ┗━━━━┷━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛

// ExtendedSquitter is an extended squitter message.
type ExtendedSquitter struct {
	ModeS
}

// TransponderCapability is the CA.
func (e ExtendedSquitter) TransponderCapability() TransponderCapability {
	return TransponderCapability(e.ModeS[0] & 0x07) //nolint: gomnd
}

// AircraftAddress is the ICAO.
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
