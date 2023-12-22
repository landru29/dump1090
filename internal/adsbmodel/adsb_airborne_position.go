package adsbmodel

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
func (s AirbornePosition) Name() string {
	return airbornePositionName
}
