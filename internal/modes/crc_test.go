package modes_test

import (
	"testing"

	"github.com/landru29/dump1090/internal/modes"
	"github.com/stretchr/testify/assert"
)

func TestCRC(t *testing.T) {
	remainder := modes.Checksumv2([]byte{0x8D, 0x40, 0x6B, 0x90, 0x20, 0x15, 0xA6, 0x78, 0xD4, 0xD2, 0x20})

	// crc32q := crc32.MakeTable(0x01fff409)
	// crc := crc32.Checksum([]byte{0x8D, 0x40, 0x6B, 0x90, 0x20, 0x15, 0xA6, 0x78, 0xD4, 0xD2, 0x20}, crc32q)
	// assert.Equal(t, 0, crc)

	assert.Equal(t, 0, remainder)
}
