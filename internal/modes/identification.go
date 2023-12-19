package modes

//       ┏━━━━┓
//       ┃ 31 ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//       ┃ TC | CA | C1 | C2 | C3 | C4 | C5 | C6 | C7 | C8 ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//       ┃ 5  |  3 |  6 |  6 |  6 |  6 |  6 |  6 |  6 |  6 ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛

const asciiTable = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"

// Identification id the aircraft identification.
type Identification struct {
	Category byte
	String   string
}

// Identification is the aircraft identification.
func (e ExtendedSquitter) Identification() (*Identification, error) {
	if e.Type != MessageTypeAircraftIdentification {
		return nil, ErrWrongMessageType
	}

	//     1     //      2    //      3    //      4    //      5    //    6
	// 876543 21 // 8765 4321 // 87 654321 // 876543 21 // 8765 4321 // 87 654321
	// \----/ \--------/ \--------/ \----/    \----/ \--------/ \--------/ \----/
	//    0        1          2        3         4        5         6         7

	letters := make([]byte, 8) //nolint: gomnd

	letters[0] = asciiTable[e.Message[1]>>2]
	letters[1] = asciiTable[(e.Message[1]&0x3)<<4+e.Message[2]>>4] //nolint: gomnd
	letters[2] = asciiTable[(e.Message[2]&0xf)<<2+e.Message[3]>>6] //nolint: gomnd
	letters[3] = asciiTable[e.Message[3]&0x3f]
	letters[4] = asciiTable[e.Message[4]>>2]
	letters[5] = asciiTable[(e.Message[4]&0x3)<<4+e.Message[5]>>4] //nolint: gomnd
	letters[6] = asciiTable[(e.Message[5]&0xf)<<2+e.Message[6]>>6] //nolint: gomnd
	letters[7] = asciiTable[e.Message[6]&0x3f]

	return &Identification{
		Category: e.Message[0] & 0x7, //nolint: gomnd
		String:   string(letters),
	}, nil
}
