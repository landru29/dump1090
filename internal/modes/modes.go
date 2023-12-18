package modes

// +---------------------------------------------------------------------------+
// |                                  Mode S                                   |
// +---------------------------------------------------------------------------+
// |                                    112                                    |
// +-----+------------------------------------------------------------+--------+
// | DF  |                     Extended squitter                      | Parity |
// +-----+------------------------------------------------------------+--------+
// |  5  |                             83                             |   24   |
// +-----+----+------+------------------------------------------------+--------+
//       | CA | ICAO |                    Message                     |
//       +----+------+------------------------------------------------+
//       | 3  |  24  |                       56                       |
//       +----+------+----+-------------------------------------------+
//                   | TC |                  Payload                  |
//                   +----+-------------------------------------------+
//                   | 5  |                    51                     |
//                   +----+-------------------------------------------+

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

// ModeS is a ModeS frame.
type ModeS struct {
	DownlinkFormat     DownlinkFormat
	ParityInterrogator uint32
	Raw                []byte
}

// Unmarshal is the mode-S unmarshaler.
func (m *ModeS) Unmarshal(data []byte) error {
	if len(data) < 4 { //nolint: gomnd
		return ErrWrongMessageSize
	}

	length := len(data)

	m.DownlinkFormat = DownlinkFormat((data[0] & 0xf8) >> 3) //nolint: gomnd

	m.Raw = data

	m.ParityInterrogator = (uint32(data[length-3]) << 16) + (uint32(data[length-2]) << 8) + //nolint: gomnd
		uint32(data[length-1])

	return nil
}
