// Package nmea is the nmea serializer.
package nmea

import (
	"bytes"
	"encoding/hex"
)

/*
	Payload (https://gpsd.gitlab.io/gpsd/AIVDM.html#_types_1_2_and_3_position_report_class_a)

	| Bits    | Size  |   Description               |
	|---------+-------+-----------------------------|
	| 0-5     | 6     |   Message Type              |
	| 6-7     | 2     |   Repeat Indicator          |
	| 8-37    | 30    |   MMSI                      |
	| 38-41   | 4     |   Navigation Status         |
	| 42-49   | 8     |   Rate of Turn (ROT)        |
	| 50-59   | 10    |   Speed Over Ground (SOG)   |
	| 60-60   | 1     |   Position Accuracy         |
	| 61-88   | 28    |   Longitude                 |
	| 89-115  | 27    |   Latitude                  |
	| 116-127 | 12    |   Course Over Ground (COG)  |
	| 128-136 | 9     |   True Heading (HDG)        |
	| 137-142 | 6     |   Time Stamp                |
	| 143-144 | 2     |   Maneuver Indicator        |
	| 145-147 | 3     |   Spare                     |
	| 148-148 | 1     |   RAIM flag                 |
	| 149-167 | 19    |   Radio status              |

*/

const (
	preambule = "!"

	talkerID = "AIVDM"

	fieldSize = 7

	fieldTalkerID            = 0
	fieldFragmentCount       = 1
	fieldFragmentNumber      = 2
	fieldSequencialMessageID = 3
	fieldRadioChannelCode    = 4
	fieldPayload             = 5
	fieldChecksum            = 6
)

type fields [][]byte

type payload struct {
	MessageType       uint8  // 6 bits
	RepeatIndicator   uint8  // 2 bits
	MMSI              uint32 // 30 bits
	NavigationStatus  uint8  // 4 bits
	RateOfTurn        int    // 8 bits
	SpeedOverGround   uint16 // 10 bits
	PositionAccuracy  uint8  // 1 bits
	Longitude         int32  // 28 bits
	Latitude          int32  // 27 bits
	CourseOverGround  uint16 // 12 bits
	TrueHeading       uint16 // 9 bits
	TimeStamp         uint8  // 6 bits
	ManeuverIndicator uint32 // 2 bits
	Spare             uint8  // 3 bits
	RaimFlag          uint8  // 1 bits
	RadioStatus       uint8  // 19 bits
}

// Serializer is the nmea serializer.
type Serializer struct {
	converter map[byte]byte
}

func New() *Serializer {
	output := Serializer{
		converter: map[byte]byte{},
	}

	for idx := byte(0); idx < 0x40; idx++ {
		if idx < 0x28 {
			output.converter[idx] = idx + 0x30
		} else {
			output.converter[idx] = idx + 0x38
		}
	}

	return &output
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	message := make(fields, fieldSize)
	message[fieldTalkerID] = []byte(talkerID)
	message[fieldFragmentCount] = []byte("1")
	message[fieldFragmentNumber] = []byte("1")
	message[fieldSequencialMessageID] = nil
	message[fieldRadioChannelCode] = []byte("B")

	message[fieldChecksum] = []byte("0*" + message.checkSum())

	return bytes.Join([][]byte{[]byte(preambule), bytes.Join(message, []byte(","))}, nil), nil
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "text/plain"
}

func (f fields) checkSum() string {
	output := byte(0)

	data := bytes.Join(f, []byte(","))

	for _, byteElement := range data {
		output ^= byteElement
	}

	return hex.EncodeToString([]byte{output})
}
