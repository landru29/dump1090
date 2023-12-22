package adsbmodel_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/adsbmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentification(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D4840D6202CC371C32CE0576098")
		require.NoError(t, err)

		require.NoError(t, adsbmodel.ModeS(dataByte).CheckSum())

		msg, err := adsbmodel.ModeS(dataByte).Message()
		require.NoError(t, err)

		assert.Equal(t, "identification", msg.Name())

		identification, ok := msg.(adsbmodel.Identification)
		assert.True(t, ok)

		assert.Equal(t, "KLM1023 ", identification.String())
		assert.Equal(t, adsbmodel.Category(0), identification.Category())
	})
}
