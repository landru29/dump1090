package adsbmodel_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/adsbmodel"
	"github.com/stretchr/testify/require"
)

func TestChecksum(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D58C382D690C8AC2863A7")
		require.NoError(t, err)

		require.NoError(t, adsbmodel.ModeS(dataByte).CheckSum())
	})

	t.Run("ko", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D59C382D690C8AC2863A7")
		require.NoError(t, err)

		require.Error(t, adsbmodel.ModeS(dataByte).CheckSum())
	})
}
