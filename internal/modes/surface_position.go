package modes

import "time"

// SurfacePosition is the surface position of the aircraft.
type SurfacePosition struct {
	GroundTrack      float64
	TimeUTC          bool
	OddFrame         bool
	EncodedLatitude  uint32
	EncodedLongitude uint32
	Time             time.Time
	Speed            float64
}

// SurfacePosition is the Surface position of the aircraft.
func (e ExtendedSquitter) SurfacePosition() (*SurfacePosition, error) {
	if e.Type != MessageTypeSurfacePosition {
		return nil, ErrWrongMessageType
	}

	//     0        //    1         //      2         //      3    //      4    //    5     //    6
	// 87654 321  // 8765  4    321 // 8765   4  3   21 // 87654321 // 8765432 1 // 87654321 // 87654321
	// \---/ \----------/  |    \---------/   |  |   \-----------------------/ \-----------------------/
	//   TC    Movement   stat    GndTrack    T Odd         Enc Latitude             Enc Longitude

	encodedMovement := ((e.Message[0] & 0x7) << 4) | ((e.Message[1]) >> 4) //nolint: gomnd
	var speed float64

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

	groundTrackByte := ((e.Message[1] & 0x7) << 4) | ((e.Message[2]) >> 4) //nolint: gomnd
	groundTrack := float64(-1)

	if (e.Message[1]>>3)&0x1 == 1 { //nolint: gomnd
		groundTrack = float64(360) * float64(groundTrackByte) / float64(128) //nolint: gomnd
	}

	return &SurfacePosition{
		GroundTrack:      groundTrack,
		TimeUTC:          map[byte]bool{1: true, 0: false}[(e.Message[2]>>3)&0x1],                                          //nolint: gomnd,lll
		OddFrame:         map[byte]bool{1: true, 0: false}[(e.Message[2]>>2)&0x1],                                          //nolint: gomnd,lll
		EncodedLatitude:  ((uint32(e.Message[2]) & 0x3) << 15) | (uint32(e.Message[3]) << 7) | (uint32(e.Message[4]) >> 1), //nolint: gomnd,lll
		EncodedLongitude: ((uint32(e.Message[4]) & 0x1) << 16) | (uint32(e.Message[5]) << 8) | uint32(e.Message[6]),        //nolint: gomnd,lll
		Time:             time.Now(),
		Speed:            speed,
	}, nil
}
