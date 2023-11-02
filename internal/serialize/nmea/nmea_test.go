package nmea //nolint: testpackage

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func displayBytes(t *testing.T, data []uint8) {
	t.Helper()

	for _, elt := range data {
		fmt.Printf("%06b ", elt) //nolint: forbidigo
	}

	fmt.Println() //nolint: forbidigo
}

func TestCheckSum(t *testing.T) {
	sentences := []string{
		"!AIVDM,1,1,,B,177KQJ5000G?tO`K>RA1wUbN0TKH,0*5C",
		"!AIVDM,1,1,,A,15RTgt0PAso;90TKcjM8h6g208CQ,0*4A",
		"!AIVDM,1,1,,A,16SteH0P00Jt63hHaa6SagvJ087r,0*42",
	}

	for _, sentence := range sentences {
		value := sentence

		t.Run(value, func(t *testing.T) {
			fields := fields(bytes.Split([]byte(value), []byte{','}))

			assert.Equal(t, value[len(value)-2:], fields.checkSum())
		})
	}
}

func TestAddData(t *testing.T) {
	for _, elt := range []struct {
		expected    []uint8
		input       any
		bitPosition uint8
		length      uint8
	}{
		{
			expected:    []uint8{0b000111, 0b111000, 0, 0, 0, 0, 0, 0, 0, 0},
			input:       uint8(63),
			bitPosition: 3,
			length:      6,
		},
	} {
		fixture := elt

		t.Run("", func(t *testing.T) {
			encoded := make([]uint8, 10)

			_, _ = payloadAddData(encoded, fixture.input, fixture.bitPosition, fixture.length)

			assert.Equal(t, fixture.expected, encoded)

			displayBytes(t, encoded)
		})
	}
}

func TestEncode(t *testing.T) {
	assert.Equal(t,
		"123456789:;<=>?@ABCDEFGHIJKL",
		encodeBinaryPayload([]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}),
	)
}
