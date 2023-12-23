package model

//       ┏━━━━┓
//       ┃ 31 ┃
//       ┣━━━━╇━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┯━━━━┓
//       ┃ TC | CA | C1 | C2 | C3 | C4 | C5 | C6 | C7 | C8 ┃
//       ┠┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┼┈┈┈┈┨
//       ┃ 5  |  3 |  6 |  6 |  6 |  6 |  6 |  6 |  6 |  6 ┃
//       ┗━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┷━━━━┛

const (
	asciiTable         = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"
	identificationName = "identification"
)

// Category is the aircraft category.
type Category byte

// Identification id the aircraft identification.
type Identification struct {
	ExtendedSquitter
}

// String implement the Stringer interface.
func (i Identification) String() string {
	message := i.Message()
	letters := make([]byte, 8) //nolint: gomnd

	letters[0] = asciiTable[message[1]>>2]
	letters[1] = asciiTable[(message[1]&0x3)<<4+message[2]>>4] //nolint: gomnd
	letters[2] = asciiTable[(message[2]&0xf)<<2+message[3]>>6] //nolint: gomnd
	letters[3] = asciiTable[message[3]&0x3f]
	letters[4] = asciiTable[message[4]>>2]
	letters[5] = asciiTable[(message[4]&0x3)<<4+message[5]>>4] //nolint: gomnd
	letters[6] = asciiTable[(message[5]&0xf)<<2+message[6]>>6] //nolint: gomnd
	letters[7] = asciiTable[message[6]&0x3f]

	return string(letters)
}

// Name implements the Message interface.
func (i Identification) Name() string {
	return identificationName
}

// Category is the aircraft category.
func (i Identification) Category() Category {
	return Category(i.Message()[0] & 0x7) //nolint: gomnd
}
