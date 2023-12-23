package model

//       ┏━━━━┓
//       ┃ 19 ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━━━┯━━━━━┯━━━━┯━━━━━┯━━━━━━┯━━━━━━┓
//       ┃ TC | ST | IC | IFR | NUC |     | VrSrc | Svr | VR | Res | SDif | dAlt ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  1 |  1  |  3  |  22 |   1   |  1  | 9  |  2  |  1   |  7   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━━━┷━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┛

const airborneVelocityName = "airborne velocity"

// AirborneVelocity is the surface position.
type AirborneVelocity struct {
	ExtendedSquitter
}

// Name implements the Message interface.
func (s AirborneVelocity) Name() string {
	return airborneVelocityName
}
