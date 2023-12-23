package model_test

import (
	"encoding/csv"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("8D4840D6202CC371C32CE0576098")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())

		squitter, err := model.ModeS(dataByte).Squitter()
		require.NoError(t, err)

		require.Equal(t, "extended squitter", squitter.Name())

		extendedSquitter, ok := squitter.(model.ExtendedSquitter)
		assert.True(t, ok)

		assert.Equal(t, model.DownlinkFormat(17), extendedSquitter.DownlinkFormat())
		assert.Equal(t, model.TransponderCapabilityAirborne, extendedSquitter.TransponderCapability())
		assert.Equal(t, model.ICAOAddr(4735190), extendedSquitter.AircraftAddress())
	})

	t.Run("any ADS-B message", func(t *testing.T) {
		t.Parallel()

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

			addrStr := record[2]

			t.Run(msgStr, func(t *testing.T) {
				dataByte, err := hex.DecodeString(msgStr)
				require.NoError(t, err)

				icaoAddr, err := strconv.ParseUint(addrStr, 16, 64)
				require.NoError(t, err)

				require.NoError(t, model.ModeS(dataByte).CheckSum())

				squitter, err := model.ModeS(dataByte).Squitter()
				require.NoError(t, err)

				require.Equal(t, "extended squitter", squitter.Name())

				extendedSquitter, ok := squitter.(model.ExtendedSquitter)
				assert.True(t, ok)

				assert.Equal(t, model.ICAOAddr(icaoAddr), extendedSquitter.AircraftAddress())
			})
		}
	})
}
