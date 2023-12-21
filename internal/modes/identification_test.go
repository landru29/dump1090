package modes_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/dump1090/internal/modes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentification(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("8D4840D6202CC371C32CE0576098")
		require.NoError(t, err)

		msg := modes.ExtendedSquitter{}

		require.NoError(t, msg.Unmarshal(dataByte))

		ident, err := msg.Identification()
		require.NoError(t, err)

		assert.Equal(t, "KLM1023 ", ident.String)
		assert.Equal(t, byte(0), ident.Category)
	})
}
