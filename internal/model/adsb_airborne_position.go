package model

import "time"

//       ┏━━━━━━━┓
//       ┃ 8-18  ┃
//       ┃ 20-23 ┃
//       ┣━━━━━━━╇━━━━┯━━━━━┯━━━━━┯━━━┯━━━┯━━━━━━━━━┯━━━━━━━━━┓
//       ┃  TC   | SS | SAF | ALT | T | F | LAT-CPR | LON-CPR ┃
//       ┠┈┈┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┼┈┈┈┼┈┈┈┈┈┈┈┈┈┼┈┈┈┈┈┈┈┈┈┨
//       ┃   5   |  2 |  1  |  12 | 1 | 1 |    17   |   17    ┃
//       ┗━━━━━━━┷━━━━┷━━━━━┷━━━━━┷━━━┷━━━┷━━━━━━━━━┷━━━━━━━━━┛

// SurveillanceStatus is the surveillance status.
type SurveillanceStatus byte

const (
	// SurveillanceStatusNoCondition is no condition.
	SurveillanceStatusNoCondition SurveillanceStatus = iota

	// SurveillanceStatusPermanentAlert is permanent alert.
	SurveillanceStatusPermanentAlert

	// SurveillanceStatusTemporaryAlert is temporary alert.
	SurveillanceStatusTemporaryAlert

	// SurveillanceStatusSpiCondition is SPI condition.
	SurveillanceStatusSpiCondition
)

// AirbornePosition is the airborne position (GNSS and Baro).
type AirbornePosition struct {
	SurveillanceStatus SurveillanceStatus
	SingleAntennaFlag  bool
	EncodedAltitude    uint16
	TimeUTC            bool
	OddFrame           bool
	EncodedLatitude    uint32
	EncodedLongitude   uint32
	Time               time.Time
	Baro               bool
}

// AirbornePosition is the airborne position of the aircraft.
func (e ExtendedSquitter) AirbornePosition() (*AirbornePosition, error) {
	if e.Type != MessageTypeAirbornePositionBaroAltitude && e.Type != MessageTypeAirbornePositionGnssHeight {
		return nil, ErrWrongMessageType
	}

	//     0        //    1     //      2         //      3    //      4    //    5     //    6
	// 87654 32  1  // 87654321 // 8765 4  3   21 // 87654321 // 8765432 1 // 87654321 // 87654321
	// \---/ \/  |     \--------------/ |  |   \-----------------------/ \-----------------------/
	//   TC  SS SAG         Enc Alt     T Odd        Enc Latitude             Enc Longitude

	return &AirbornePosition{
		SurveillanceStatus: SurveillanceStatus((e.Message[0] & 0x6) >> 1), //nolint: gomnd
		SingleAntennaFlag:  map[byte]bool{1: true, 0: false}[e.Message[0]&0x1],
		EncodedAltitude:    (uint16(e.Message[1]) << 4) | (uint16(e.Message[2]) >> 4),                                        //nolint: gomnd,lll
		TimeUTC:            map[byte]bool{1: true, 0: false}[(e.Message[2]>>3)&0x1],                                          //nolint: gomnd,lll
		OddFrame:           map[byte]bool{1: true, 0: false}[(e.Message[2]>>2)&0x1],                                          //nolint: gomnd,lll
		EncodedLatitude:    ((uint32(e.Message[2]) & 0x3) << 15) | (uint32(e.Message[3]) << 7) | (uint32(e.Message[4]) >> 1), //nolint: gomnd,lll
		EncodedLongitude:   ((uint32(e.Message[4]) & 0x1) << 16) | (uint32(e.Message[5]) << 8) | uint32(e.Message[6]),        //nolint: gomnd,lll
		Time:               time.Now(),
		Baro:               e.Type == MessageTypeAirbornePositionBaroAltitude,
	}, nil
}
