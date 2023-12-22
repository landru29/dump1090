package adsbmodel_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/adsbmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAirbornePosition(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D58C382D690C8AC2863A7")
		require.NoError(t, err)

		require.NoError(t, adsbmodel.ModeS(dataByte).CheckSum())

		msg, err := adsbmodel.ModeS(dataByte).Message()
		require.NoError(t, err)

		assert.Equal(t, "airborne position", msg.Name())

		_, ok := msg.(adsbmodel.AirbornePosition)
		assert.True(t, ok)
	})
}
