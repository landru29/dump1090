package model

//       ┏━━━━━┓
//       ┃ 4-9 ┃
//       ┣━━━━━╇━━━━━┯━━━┯━━━━━┯━━━┯━━━┯━━━━━━━━━┯━━━━━━━━━┓
//       ┃ TC  | MOV | S | TRK | T | F | LAT-CPR | LON-CPR ┃
//       ┠┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┼┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┨
//       ┃ 5   |  7  | 1 |  7  | 1 | 1 |    17   |   17    ┃
//       ┗━━━━━┷━━━━━┷━━━┷━━━━━┷━━━┷━━━┷━━━━━━━━━┷━━━━━━━━━┛

const surfacePositionName = "surface position"

// SurfacePosition is the surface position.
type SurfacePosition struct {
	ExtendedSquitter
}

// Name implements the Message interface.
func (p SurfacePosition) Name() string {
	return surfacePositionName
}

// GroundTrack is the ground track.
func (p SurfacePosition) GroundTrack() float64 {
	message := p.Message()

	groundTrackByte := ((message[1] & 0x7) << 4) | ((message[2]) >> 4) //nolint: gomnd
	groundTrack := float64(-1)

	if (message[1]>>3)&0x1 == 1 { //nolint: gomnd
		groundTrack = float64(360) * float64(groundTrackByte) / float64(128) //nolint: gomnd
	}

	return groundTrack
}

// TimeUTC define whether the time is UTC or not.
func (p SurfacePosition) TimeUTC() bool {
	return map[byte]bool{1: true, 0: false}[(p.Message()[2]>>3)&0x1] //nolint: gomnd
}

// OddFrame defines if the frame is odd or even.
func (p SurfacePosition) OddFrame() bool {
	return map[byte]bool{1: true, 0: false}[(p.Message()[2]>>2)&0x1] //nolint: gomnd
}

// EncodedLatitude is the encoded latitude.
func (p SurfacePosition) EncodedLatitude() uint32 {
	message := p.Message()

	return ((uint32(message[2]) & 0x3) << 15) | (uint32(message[3]) << 7) | (uint32(message[4]) >> 1) //nolint: gomnd
}

// EncodedLongitude is the encoded longitude.
func (p SurfacePosition) EncodedLongitude() uint32 {
	message := p.Message()

	return ((uint32(message[4]) & 0x1) << 16) | (uint32(message[5]) << 8) | uint32(message[6]) //nolint: gomnd
}

// Speed is the ground speed.
func (p SurfacePosition) Speed() float64 {
	var speed float64

	message := p.Message()

	encodedMovement := ((message[0] & 0x7) << 4) | ((message[1]) >> 4) //nolint: gomnd

	switch {
	case encodedMovement == 0:
		speed = -1
	case encodedMovement == 1:
		speed = 0
	case encodedMovement < 8: //nolint: gomnd
		speed = float64(0.125) * float64(encodedMovement-1) //nolint: gomnd
	case encodedMovement < 12: //nolint: gomnd
		speed = float64(1) + float64(0.25)*float64(encodedMovement-9) //nolint: gomnd
	case encodedMovement < 38: //nolint: gomnd
		speed = float64(2) + float64(0.5)*float64(encodedMovement-13) //nolint: gomnd
	case encodedMovement < 93: //nolint: gomnd
		speed = float64(15) + float64(encodedMovement-39) //nolint: gomnd
	case encodedMovement < 108: //nolint: gomnd
		speed = float64(70) + float64(2)*float64(encodedMovement-94) //nolint: gomnd
	case encodedMovement < 123: //nolint: gomnd
		speed = float64(100) + float64(5)*float64(encodedMovement-109) //nolint: gomnd
	default:
		speed = 175.1
	}

	return speed
}
