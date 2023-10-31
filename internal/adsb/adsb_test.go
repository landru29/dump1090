package adsb_test

import (
	"encoding/csv"
	"encoding/hex"
	"io"
	"os"
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
		assert.Equal(t, adsb.TypeCode(0x04), msg.TypeCode)
		assert.Equal(t, uint32(0x576098), msg.ParityInterrogator)
		assert.Len(t, msg.Message, 56/8)
	})

	t.Run("any ADS-B message", func(t *testing.T) {
		file, err := os.Open("testdata/sample_data_adsb.csv")
		require.NoError(t, err)

		defer func(closer io.Closer) {
			_ = closer.Close()
		}(file)

		reader := csv.NewReader(file)

		records, err := reader.ReadAll()
		require.NoError(t, err)

		for _, record := range records {
			msgStr := record[1]
			t.Run(msgStr, func(t *testing.T) {
				dataByte, err := hex.DecodeString(msgStr)
				require.NoError(t, err)

				msg := &adsb.Message{}

				assert.NoError(t, msg.Unmarshal(dataByte))
			})
		}
	})
}
