package adsb_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/adsb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		dataByte, err := hex.DecodeString("8D4840D6202CC371C32CE0576098")
		require.NoError(t, err)

		msg := &adsb.Message{}

		assert.NoError(t, msg.Unmarshal(dataByte))

		assert.Equal(t, adsb.DownlinkFormat(17), msg.DownlinkFormat)
		assert.Equal(t, byte(5), msg.TransponderCapability)
		assert.Equal(t, uint32(4735190), msg.AircraftAddress)
		assert.Equal(t, adsb.TypeCode(0x11), msg.TypeCode)
	})
}
