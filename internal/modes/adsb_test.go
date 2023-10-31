package modes_test

import (
	"bufio"
	"encoding/csv"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/landru29/dump1090/internal/modes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		dataByte, err := hex.DecodeString("8D4840D6202CC371C32CE0576098")
		require.NoError(t, err)

		modeS := &modes.Frame{}

		assert.NoError(t, modeS.Unmarshal(dataByte))

		msg := modes.Message{}

		assert.NoError(t, msg.Unmarshal(*modeS))

		assert.Equal(t, modes.DownlinkFormat(17), msg.ModeS.DownlinkFormat)
		assert.Equal(t, byte(5), msg.TransponderCapability)
		assert.Equal(t, uint32(4735190), msg.AircraftAddress)
		assert.Equal(t, modes.TypeCode(0x04), msg.TypeCode)
		assert.Equal(t, uint32(0x576098), msg.ModeS.ParityInterrogator)
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

				modeS := &modes.Frame{}

				assert.NoError(t, modeS.Unmarshal(dataByte))

				msg := modes.Message{}

				assert.NoError(t, msg.Unmarshal(*modeS))
			})
		}
	})

	t.Run("from Dump1090", func(t *testing.T) {
		file, err := os.Open("testdata/dump1090.txt")
		require.NoError(t, err)

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			line := scanner.Text()

			t.Run(line, func(t *testing.T) {
				dataByte, err := hex.DecodeString(line[1 : len(line)-1])
				require.NoError(t, err)

				modeS := &modes.Frame{}

				assert.NoError(t, modeS.Unmarshal(dataByte))
			})
		}
	})
}
