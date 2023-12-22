package model

//       ┏━━━━┓
//       ┃ 19 ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━━━┯━━━━━┯━━━━┯━━━━━┯━━━━━━┯━━━━━━┓
//       ┃ TC | ST | IC | IFR | NUC |     | VrSrc | Svr | VR | Res | SDif | dAlt ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┨
//       ┃ 5  |  3 |  1 |  1  |  3  |  22 |   1   |  1  | 9  |  2  |  1   |  7   ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━━━┷━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┛

// AirborneVelocity is the airborne velocity.
type AirborneVelocity struct{}

// AirborneVelocity is the airborne velocity of the aircraft.
func (e ExtendedSquitter) AirborneVelocity() (*AirborneVelocity, error) {
	if e.Type != MessageTypeAirborneVelocities {
		return nil, ErrWrongMessageType
	}

	return &AirborneVelocity{}, nil
}
