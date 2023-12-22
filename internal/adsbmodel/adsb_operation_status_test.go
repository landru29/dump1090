package adsbmodel_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/adsbmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationStatus(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D4CA92BF8230006004BB8FB39CA")
		require.NoError(t, err)

		require.NoError(t, adsbmodel.ModeS(dataByte).CheckSum())

		msg, err := adsbmodel.ModeS(dataByte).Message()
		require.NoError(t, err)

		assert.Equal(t, "operation status", msg.Name())

		_, ok := msg.(adsbmodel.OperationStatus)
		assert.True(t, ok)
	})
}
