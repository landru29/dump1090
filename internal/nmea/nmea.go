package nmea

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"unsafe"

	"github.com/landru29/dump1090/internal/errors"
)

const (
	bitSize  = 6 // nmea byte size
	byteSize = 8 // real byte size

	fieldTalkerID            = 0
	fieldFragmentCount       = 1
	fieldFragmentNumber      = 2
	fieldSequencialMessageID = 3
	fieldRadioChannelCode    = 4
	fieldPayload             = 5
	fieldChecksum            = 6

	ErrDataTooLong     errors.Error = "length is over data capacity"
	ErrUnsupportedType errors.Error = "unsupported data type"
)

type Fields [][]byte

func (f Fields) CheckSum() string {
	output := byte(0)

	data := bytes.Join(f, []byte(","))
	currentField := 0

	for idx, byteElement := range data {
		if idx == 0 && byteElement == '!' {
			continue
		}

		if byteElement == ',' {
			currentField++
		}

		if currentField == fieldChecksum && byteElement == '*' {
			break
		}

		output ^= byteElement
	}

	return strings.ToUpper(hex.EncodeToString([]byte{output}))
}

func addBytes(dest []uint8, data []uint8, bitPosition uint32, length uint32) {
	for idx := uint32(0); idx < length; idx++ {
		readBit := (data[idx/byteSize] << (idx % byteSize)) & 0x80

		dest[(bitPosition+idx)%bitSize] |= (readBit >> (uint8((bitPosition+idx)/bitSize) + 2))
	}
}

func AddData(dest []uint8, data any, bitPosition uint32, length uint32) (uint32, error) {
	if length > uint32(unsafe.Sizeof(data)) {
		return 0, ErrDataTooLong
	}

	switch value := data.(type) {
	case uint8:
		addBytes(dest, []uint8{value}, bitPosition, length)
		return length, nil
	case uint16:
		encoded := make([]uint8, 2)

		binary.BigEndian.AppendUint16(encoded, value)

		addBytes(dest, encoded, bitPosition, length)

		return length, nil
	case uint32:
		encoded := make([]uint8, 4)

		binary.BigEndian.AppendUint32(encoded, value)

		addBytes(dest, encoded, bitPosition, length)

		return length, nil
	case bool:
		if value {
			addBytes(dest, []byte{1}, bitPosition, 1)
		}

		return 1, nil
	case int8:
	case int16:
	case int32:
	}

	return 0, ErrUnsupportedType
}
