package binary_test

import (
	"testing"

	"github.com/landru29/dump1090/internal/binary"
	"github.com/stretchr/testify/assert"
)

func TestReadBits(t *testing.T) {
	t.Parallel()

	t.Run("24 bits in the middle of a byte", func(t *testing.T) {
		t.Parallel()

		// 10111010 11111011 10011101 10110001 01110011 11010110 01101011 11111101 01101101 01011011
		//              1011 10011101 10110001 0111
		data := []byte{0xba, 0xfb, 0x9d, 0xb1, 0x73, 0xd6, 0x6b, 0xfd, 0x6d, 0x5b}

		bits := binary.ReadBits(data, 12, 24)

		assert.Equal(t, uint64(0xb9db17), bits)
	})

	t.Run("24 bits at the start of a byte", func(t *testing.T) {
		t.Parallel()

		// 10111010 11111011 10011101 10110001 01110011 11010110 01101011 11111101 01101101 01011011
		//          11111011 10011101 10110001
		data := []byte{0xba, 0xfb, 0x9d, 0xb1, 0x73, 0xd6, 0x6b, 0xfd, 0x6d, 0x5b}

		bits := binary.ReadBits(data, 8, 24)

		assert.Equal(t, uint64(0xfb9db1), bits)
	})

	t.Run("any size", func(t *testing.T) {
		t.Parallel()

		// 10111010 11111011 10011101 10110001 01110011 11010110 01101011 11111101 01101101 01011011
		//            111011 10011101 10110
		data := []byte{0xba, 0xfb, 0x9d, 0xb1, 0x73, 0xd6, 0x6b, 0xfd, 0x6d, 0x5b}

		bits := binary.ReadBits(data, 10, 19)

		assert.Equal(t, uint64(0x773b6), bits)
	})
}

func TestWriteBits(t *testing.T) {
	t.Parallel()

	t.Run("24 bits in the middle of a byte", func(t *testing.T) {
		t.Parallel()

		// 00000000 00001011 10011101 10110001 01110000 00000000 00000000 00000000 00000000 00000000
		//              1011 10011101 10110001 0111
		data := make([]byte, 10)

		binary.WriteBits(data, uint64(0xb9db17), 12, 24)

		assert.Equal(t, []byte{0x00, 0x0b, 0x9d, 0xb1, 0x70, 0x00, 0x00, 0x00, 0x00, 0x00}, data)
	})

	t.Run("24 bits at the start of a byte", func(t *testing.T) {
		t.Parallel()

		// 00000000 11111011 10011101 10110001 00000000 00000000 00000000 00000000 00000000 00000000
		//          11111011 10011101 10110001
		data := make([]byte, 10)

		binary.WriteBits(data, uint64(0xfb9db1), 8, 24)

		assert.Equal(t, []byte{0x00, 0xfb, 0x9d, 0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, data)
	})

	t.Run("any size", func(t *testing.T) {
		t.Parallel()

		// 00000000 00111011 10011101 10110000 00000000 00000000 00000000 00000000 00000000 00000000
		//            111011 10011101 10110
		data := make([]byte, 10)

		binary.WriteBits(data, uint64(0x773b6), 10, 19)

		assert.Equal(t, []byte{0x00, 0x3b, 0x9d, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, data)
	})

	t.Run("any size", func(t *testing.T) {
		t.Parallel()

		// 11111111 11111011 10011101 10110111 11111111 11111111 11111111 11111111 11111111 11111111
		//            111011 10011101 10110
		data := make([]byte, 10)

		for idx := range data {
			data[idx] = 0xff
		}

		binary.WriteBits(data, uint64(0x773b6), 10, 19)

		assert.Equal(t, []byte{0xff, 0xfb, 0x9d, 0xb7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, data)
	})

	t.Run("zero in one", func(t *testing.T) {
		t.Parallel()

		// 11111111 11000000 00000000 00000111 11111111 11111111 11111111 11111111 11111111 11111111
		//            000000 00000000 00000
		data := make([]byte, 10)

		for idx := range data {
			data[idx] = 0xff
		}

		binary.WriteBits(data, uint64(0), 10, 19)

		assert.Equal(t, []byte{0xff, 0xc0, 0x00, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, data)
	})

	t.Run("one in zero", func(t *testing.T) {
		t.Parallel()

		// 00000000 00111111 11111111 11111000 00000000 00000000 00000000 00000000 00000000 00000000
		//            111111 11111111 11111
		data := make([]byte, 10)

		binary.WriteBits(data, uint64(0x7ffff), 10, 19)

		assert.Equal(t, []byte{0x00, 0x3f, 0xff, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, data)
	})
}
