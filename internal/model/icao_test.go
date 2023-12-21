package model_test

import (
	"encoding/json"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalICAO(t *testing.T) {
	t.Parallel()

	icao := model.ICAOAddr(0xabcdef)

	out, err := json.Marshal(icao)
	require.NoError(t, err)

	assert.Equal(t, `"ABCDEF"`, string(out))
}

func TestUnmarshalICAO(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		var icao model.ICAOAddr

		err := json.Unmarshal([]byte(`"123456"`), &icao)
		require.NoError(t, err)

		assert.Equal(t, model.ICAOAddr(0x123456), icao)
	})

	t.Run("empty value", func(t *testing.T) {
		t.Parallel()

		var icao model.ICAOAddr

		err := json.Unmarshal(nil, &icao)
		require.Error(t, err)
	})

	t.Run("missing quotes", func(t *testing.T) {
		t.Parallel()

		var icao model.ICAOAddr

		err := json.Unmarshal([]byte(`123456`), &icao)
		require.Error(t, err)
	})
}
