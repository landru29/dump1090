package model_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationStatus(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D4CA92BF8230006004BB8FB39CA")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())

		squitter, err := model.ModeS(dataByte).Squitter()
		require.NoError(t, err)

		require.Equal(t, "extended squitter", squitter.Name())

		extendedSquitter, ok := squitter.(model.ExtendedSquitter)
		assert.True(t, ok)

		msg, err := extendedSquitter.Decode()
		require.NoError(t, err)

		assert.Equal(t, "operation status", msg.Name())

		_, ok = msg.(model.OperationStatus)
		assert.True(t, ok)
	})
}
