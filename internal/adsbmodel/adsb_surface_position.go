package adsbmodel

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
func (s SurfacePosition) Name() string {
	return surfacePositionName
}
