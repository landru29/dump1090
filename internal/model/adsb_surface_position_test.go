package model_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurfacePosition(t *testing.T) {
	t.Parallel()

	for idx, fixtureElt := range []struct {
		input            string
		groundTrack      float64
		timeUTC          bool
		oddFrame         bool
		encodedLatitude  uint32
		encodedLongitude uint32
		speed            float64
	}{
		{
			input:            "8C4841753AAB238733C8CD4020B1",
			groundTrack:      140.625,
			timeUTC:          false,
			oddFrame:         false,
			encodedLatitude:  0x1c399,
			encodedLongitude: 0x1c8cd,
			speed:            18,
		},
		{
			input:            "8C4841753A8A35323FAEBDAC702D",
			groundTrack:      98.4375,
			timeUTC:          false,
			oddFrame:         true,
			encodedLatitude:  0x991f,
			encodedLongitude: 0x1aebd,
			speed:            16,
		},
	} {
		fixture := fixtureElt

		t.Run(fmt.Sprintf("%d: %s", idx, fixture.input), func(t *testing.T) {
			t.Parallel()
			dataByte, err := hex.DecodeString(fixture.input)
			require.NoError(t, err)

			require.NoError(t, model.ModeS(dataByte).CheckSum())

			squitter, err := model.ModeS(dataByte).Squitter()
			require.NoError(t, err)

			require.Equal(t, "extended squitter", squitter.Name())

			extendedSquitter, ok := squitter.(model.ExtendedSquitter)
			assert.True(t, ok)

			msg, err := extendedSquitter.Decode()
			require.NoError(t, err)

			assert.Equal(t, "surface position", msg.Name())

			position, ok := msg.(model.SurfacePosition)
			assert.True(t, ok)

			assert.Equal(t, fixture.groundTrack, position.GroundTrack()) //nolint: testifylint
			assert.Equal(t, fixture.timeUTC, position.TimeUTC())
			assert.Equal(t, fixture.oddFrame, position.OddFrame())
			assert.Equal(t, fixture.encodedLatitude, position.EncodedLatitude())
			assert.Equal(t, fixture.encodedLongitude, position.EncodedLongitude())
			assert.Equal(t, fixture.speed, position.Speed()) //nolint: testifylint
		})
	}
}
