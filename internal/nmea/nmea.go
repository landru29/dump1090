package nmea

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"strings"

	"github.com/landru29/dump1090/internal/errors"
)

const (
	bitSize  = 6 // nmea byte size
	byteSize = 8 // real byte size

	// fieldTalkerID            = 0
	// fieldFragmentCount       = 1
	// fieldFragmentNumber      = 2
	// fieldSequencialMessageID = 3
	// fieldRadioChannelCode    = 4
	// fieldPayload             = 5
	// fieldChecksum            = 6

	ErrDataTooLong     errors.Error = "length is over data capacity"
	ErrUnsupportedType errors.Error = "unsupported data type"
)

type Fields [][]byte

func (f Fields) CheckSum() string {
	output := byte(0)

	data := bytes.Join(f, []byte(","))
	currentField := 0

	for idx, byteElement := range data {
		if idx == 0 && byteElement == '!' || byteElement == '$' {
			continue
		}

		if byteElement == ',' {
			currentField++
		}

		if currentField == len(f)-1 && byteElement == '*' {
			break
		}

		output ^= byteElement
	}

	return strings.ToUpper(hex.EncodeToString([]byte{output}))
}

func (f Fields) String() string {
	return string(bytes.Join(f, []byte(",")))
}

func payloadAddBytes(dest []uint8, data []uint8, bitPosition uint8, length uint8) {
	startInputBit := uint8(len(data))*8 - length
	for idx := uint8(0); idx < length; idx++ {
		readBit := (data[(startInputBit+idx)/byteSize] << ((startInputBit + idx) % byteSize)) & 0x80

		writeBit := readBit >> (uint8((bitPosition+idx)%bitSize) + 2)
		dest[(bitPosition+idx)/bitSize] |= writeBit
	}
}

func PayloadAddData(dest []uint8, data any, bitPosition uint8, length uint8) (uint8, error) {
	var encoded []uint8

	switch value := data.(type) {
	case uint64:
		encoded = make([]uint8, 8)
		binary.BigEndian.PutUint64(encoded, value)

	case bool:
		encoded = []uint8{0}
		if value {
			encoded = []uint8{1}
		}

	case int64:
		encoded = make([]uint8, 8)
		binary.BigEndian.PutUint64(encoded, uint64(value))

	case uint8:
		return PayloadAddData(dest, uint64(value), bitPosition, length)
	case uint16:
		return PayloadAddData(dest, uint64(value), bitPosition, length)
	case uint32:
		return PayloadAddData(dest, uint64(value), bitPosition, length)

	case int8:
		return PayloadAddData(dest, int64(value), bitPosition, length)
	case int16:
		return PayloadAddData(dest, int64(value), bitPosition, length)
	case int32:
		return PayloadAddData(dest, int64(value), bitPosition, length)

	case float64:
		return PayloadAddData(dest, math.Float64bits(value), bitPosition, length)
	case float32:
		return PayloadAddData(dest, math.Float64bits(float64(value)), bitPosition, length)

	default:
		return 0, ErrUnsupportedType

	}

	payloadAddBytes(dest, encoded, bitPosition, length)

	return length, nil
}

func EncodeBinaryPayload(input []uint8) string {
	str := ""
	for idx, elt := range input {
		if (elt & 0x3f) > 0x27 {
			str = str[:idx] + string(byte(elt)+56)
		} else {
			str = str[:idx] + string(byte(elt)+48)
		}
	}
	return str
}
