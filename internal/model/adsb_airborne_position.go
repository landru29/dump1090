package model

//       ┏━━━━━━━┓
//       ┃ 8-18  ┃
//       ┃ 20-23 ┃
//       ┣━━━━━━━╇━━━━┯━━━━━┯━━━━━┯━━━┯━━━┯━━━━━━━━━┯━━━━━━━━━┓
//       ┃  TC   | SS | SAF | ALT | T | F | LAT-CPR | LON-CPR ┃
//       ┠┈┈┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┼┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┨
//       ┃   5   |  2 |  1  |  12 | 1 | 1 |    17   |   17    ┃
//       ┗━━━━━━━┷━━━━┷━━━━━┷━━━━━┷━━━┷━━━┷━━━━━━━━━┷━━━━━━━━━┛

const airbornePositionName = "airborne position"

// AirbornePosition is the surface position.
type AirbornePosition struct {
	ExtendedSquitter
}

// Name implements the Message interface.
func (p AirbornePosition) Name() string {
	return airbornePositionName
}

// SurveillanceStatus is the surveillance status.
func (p AirbornePosition) SurveillanceStatus() SurveillanceStatus {
	return SurveillanceStatus((p.Message()[0] & 0x6) >> 1) //nolint: gomnd
}

// SingleAntennaFlag defines if the antenna is single or dual.
func (p AirbornePosition) SingleAntennaFlag() bool {
	return map[byte]bool{1: true, 0: false}[p.Message()[0]&0x1]
}

// EncodedAltitude is the encoded altitude.
func (p AirbornePosition) EncodedAltitude() uint16 {
	message := p.Message()

	return (uint16(message[1]) << 4) | (uint16(message[2]) >> 4) //nolint: gomnd
}

// TimeUTC define whether the time is UTC or not.
func (p AirbornePosition) TimeUTC() bool {
	return map[byte]bool{1: true, 0: false}[(p.Message()[2]>>3)&0x1] //nolint: gomnd
}

// OddFrame defines if the frame is odd or even.
func (p AirbornePosition) OddFrame() bool {
	return map[byte]bool{1: true, 0: false}[(p.Message()[2]>>2)&0x1] //nolint: gomnd
}

// EncodedLatitude is the encoded latitude.
func (p AirbornePosition) EncodedLatitude() uint32 {
	message := p.Message()

	return ((uint32(message[2]) & 0x3) << 15) | //nolint: gomnd
		(uint32(message[3]) << 7) | //nolint: gomnd
		(uint32(message[4]) >> 1)
}

// EncodedLongitude is the encoded longitude.
func (p AirbornePosition) EncodedLongitude() uint32 {
	message := p.Message()

	return ((uint32(message[4]) & 0x1) << 16) | //nolint: gomnd
		(uint32(message[5]) << 8) | //nolint: gomnd
		uint32(message[6])
}

// Baro defines whether the message is based on a baro altitude or GNSS altitude.
func (p AirbornePosition) Baro() bool {
	return p.TypeCode() == TypeCodeAirbornePositionBaroAltitude
}
