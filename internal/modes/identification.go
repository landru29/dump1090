package modes

const asciiTable = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"

// Identification id the aircraft identification.
type Identification struct {
	Message  ExtendedSquitter
	Category byte
	String   string
}

// Identification id the aircraft identification.
func (e ExtendedSquitter) Identification() (*Identification, error) {
	if e.Type() != MessageTypeAircraftIdentification {
		return nil, ErrWrongMessageType
	}

	//     1     //      2    //      3    //      4    //      5    //    6
	// 876543 21 // 8765 4321 // 87 654321 // 876543 21 // 8765 4321 // 87 654321
	// \----/ \--------/ \--------/ \----/    \----/ \--------/ \--------/ \----/
	//    0        1          2        3         4        5         6         7

	letters := make([]byte, 8)

	letters[0] = asciiTable[e.Message[1]>>2]
	letters[1] = asciiTable[(e.Message[1]&0x3)<<4+e.Message[2]>>4]
	letters[2] = asciiTable[(e.Message[2]&0xf)<<2+e.Message[3]>>6]
	letters[3] = asciiTable[e.Message[3]&0x3f]
	letters[4] = asciiTable[e.Message[4]>>2]
	letters[5] = asciiTable[(e.Message[4]&0x3)<<4+e.Message[5]>>4]
	letters[6] = asciiTable[(e.Message[5]&0xf)<<2+e.Message[6]>>6]
	letters[7] = asciiTable[e.Message[6]&0x3f]

	return &Identification{
		Message:  e,
		Category: e.Message[0] & 0x7,
		String:   string(letters),
	}, nil
}
