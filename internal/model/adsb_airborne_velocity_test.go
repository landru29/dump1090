package model_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAirborneVelocity(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D34620499083E383008054D8CB4")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())

		squitter, err := model.ModeS(dataByte).Squitter()
		require.NoError(t, err)

		require.Equal(t, "extended squitter", squitter.Name())

		extendedSquitter, ok := squitter.(model.ExtendedSquitter)
		assert.True(t, ok)

		msg, err := extendedSquitter.Decode()
		require.NoError(t, err)

		assert.Equal(t, "airborne velocity", msg.Name())

		_, ok = msg.(model.AirborneVelocity)
		assert.True(t, ok)
	})
}
