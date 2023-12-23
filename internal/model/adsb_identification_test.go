package model_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentification(t *testing.T) {
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

		msg, err := extendedSquitter.Decode()
		require.NoError(t, err)

		assert.Equal(t, "identification", msg.Name())

		identification, ok := msg.(model.Identification)
		assert.True(t, ok)

		assert.Equal(t, "KLM1023 ", identification.String())
		assert.Equal(t, model.Category(0), identification.Category())
	})
}
