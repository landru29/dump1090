package binary_test

import (
	"testing"

	"github.com/landru29/dump1090/internal/binary"
	"github.com/stretchr/testify/assert"
)

// https://mode-s.org/decode/content/ads-b/8-error-control.html#ads-b-parity
func TestCRC(t *testing.T) {
	t.Parallel()

	remainder := binary.ChecksumSquitter([]byte{0x8D, 0x40, 0x6B, 0x90, 0x20, 0x15, 0xA6, 0x78, 0xD4, 0xD2, 0x20})

	assert.Equal(t, uint32(0xaa4bda), remainder)
}
