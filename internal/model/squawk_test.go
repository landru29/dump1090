package model_test

import (
	"encoding/json"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalSquawk(t *testing.T) {
	ident := model.Squawk(7123)

	out, err := json.Marshal(ident)
	require.NoError(t, err)

	assert.Equal(t, "7123", string(out))
}

func TestUnmarshalSquawk(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		var icao model.Squawk

		err := json.Unmarshal([]byte("7123"), &icao)
		require.NoError(t, err)

		assert.Equal(t, model.Squawk(7123), icao)
	})

	t.Run("too high", func(t *testing.T) {
		var icao model.Squawk

		err := json.Unmarshal([]byte("17123"), &icao)
		require.Error(t, err)
	})

	t.Run("digit", func(t *testing.T) {
		var icao model.Squawk

		err := json.Unmarshal([]byte(`8900`), &icao)
		require.Error(t, err)
	})
}
