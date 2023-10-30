// Package adsb is the Automatic Dependent Surveillance-Broadcast.
package adsb

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
}

// Unmarshal parses the message.
func (m *Message) Unmarshal(data []byte) error {
	m.DownlinkFormat = DownlinkFormat((data[0] & 0xf8) >> 3)

	m.TransponderCapability = data[0] & 0x07

	m.AircraftAddress = (uint32(data[1]) << 16) + (uint32(data[2]) << 8) + uint32(data[3])

	m.TypeCode = TypeCode((data[0] & 0xf8) >> 3)

	switch {
	case m.TypeCode > 0 && m.TypeCode < 5:
		// Aircraft Identification.
	case m.TypeCode > 4 && m.TypeCode < 9:
		// Surface position.
	case m.TypeCode > 8 && m.TypeCode < 19:
		// Airborne position (w/Baro Altitude).
	case m.TypeCode == 19:
		// Airborn velocities.
	case m.TypeCode > 19 && m.TypeCode < 23:
		// Airborne position (w/GNSS Height).
	case m.TypeCode == 28:
		// Aircraft status.
	case m.TypeCode == 29:
		// Target state and status information.
	case m.TypeCode == 31:
		// Aircraft operation status.
	}

	return nil
}
