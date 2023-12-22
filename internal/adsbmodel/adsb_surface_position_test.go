package adsbmodel_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/adsbmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurfacePosition(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8C4841753AAB238733C8CD4020B1")
		require.NoError(t, err)

		require.NoError(t, adsbmodel.ModeS(dataByte).CheckSum())

		msg, err := adsbmodel.ModeS(dataByte).Message()
		require.NoError(t, err)

		assert.Equal(t, "surface position", msg.Name())

		_, ok := msg.(adsbmodel.SurfacePosition)
		assert.True(t, ok)
	})
}
