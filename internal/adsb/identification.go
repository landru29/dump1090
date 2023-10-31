package adsb

const asciiTable = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"

// Identification id the aircraft identification.
type Identification struct {
	Message  Message
	Category byte
	String   string
}

// Identification id the aircraft identification.
func (m Message) Identification() (*Identification, error) {
	if m.Type() != MessageTypeAircraftIdentification {
		return nil, ErrWrongMessageType
	}

	//     1     //      2    //      3    //      4    //      5    //    6
	// 876543 21 // 8765 4321 // 87 654321 // 876543 21 // 8765 4321 // 87 654321
	// \----/ \--------/ \--------/ \----/    \----/ \--------/ \--------/ \----/
	//    0        1          2        3         4        5         6         7

	letters := make([]byte, 8)

	letters[0] = asciiTable[m.Message[1]>>2]
	letters[1] = asciiTable[(m.Message[1]&0x3)<<4+m.Message[2]>>4]
	letters[2] = asciiTable[(m.Message[2]&0xf)<<2+m.Message[3]>>6]
	letters[3] = asciiTable[m.Message[3]&0x3f]
	letters[4] = asciiTable[m.Message[4]>>2]
	letters[5] = asciiTable[(m.Message[4]&0x3)<<4+m.Message[5]>>4]
	letters[6] = asciiTable[(m.Message[5]&0xf)<<2+m.Message[6]>>6]
	letters[7] = asciiTable[m.Message[6]&0x3f]

	return &Identification{
		Message:  m,
		Category: m.Message[0] & 0x7,
		String:   string(letters),
	}, nil
}
