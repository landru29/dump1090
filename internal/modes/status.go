package modes

//       ┏━━━━━┓
//       ┃ 1-4 ┃
//       ┣━━━━━╇━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
//       ┃ TC  | ST |                                                            ┃
//       ┠┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┨
//       ┃ 5   |  3 |                         48                                 ┃
//       ┗━━━━━╈━━━━╇━━━━━┯━━━━┯━━━━━┯━━━━━━┯━━━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━┯━━━━━┫
//             ┃ =0 | CC  | OM | Ver | NICs | NACp | BAQ | SIL | BAI | HRD | Res ┃
//             ┗━━━━╅┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┨
//                  ┃ 16  | 16 |  3  |  1   |  4   |  2  |  2  |  1  |  1  |  2  ┃
//             ┏━━━━╇━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┫
//             ┃ =1 | CC  | OM | Ver | NICs | NACp | Res | SIL | BAI | HRD | Res ┃
//             ┗━━━━╅┈┈┈┈┈┼┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┼┈┈┈┈┈┨
//                  ┃ 16  | 16 |  3  |  1   |  4   |  2  |  2  |  1  |  1  |  2  ┃
//                  ┗━━━━━┷━━━━┷━━━━━┷━━━━━━┷━━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┷━━━━━┛

// OperationStatus is the operation status.
type OperationStatus struct{}

// OperationStatus is the operation status of the aircraft.
func (e ExtendedSquitter) OperationStatus() (*OperationStatus, error) {
	if e.Type != MessageTypeAircraftOperationStatus {
		return nil, ErrWrongMessageType
	}

	return &OperationStatus{}, nil
}
