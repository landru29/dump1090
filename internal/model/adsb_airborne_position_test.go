package model_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func airbornPosition(t *testing.T, input string) model.AirbornePosition {
	t.Helper()

	dataByte, err := hex.DecodeString(input)
	require.NoError(t, err)

	require.NoError(t, model.ModeS(dataByte).CheckSum())

	squitter, err := model.ModeS(dataByte).Squitter()
	require.NoError(t, err)

	require.Equal(t, "extended squitter", squitter.Name())

	extendedSquitter, ok := squitter.(model.ExtendedSquitter)
	assert.True(t, ok)

	msg, err := extendedSquitter.Decode()
	require.NoError(t, err)

	assert.Equal(t, "airborne position", msg.Name())

	position, ok := msg.(model.AirbornePosition)
	assert.True(t, ok)

	return position
}

func TestAirbornePosition(t *testing.T) {
	t.Parallel()

	t.Run("position ok", func(t *testing.T) {
		t.Parallel()

		current := airbornPosition(t, "8D40621D58C386435CC412692AD6")
		other := airbornPosition(t, "8D40621D58C382D690C8AC2863A7")

		pos, err := current.DecodePosition(other)
		require.NoError(t, err)
		require.NotNil(t, pos)

		assert.InDelta(t, 52.26578017412606, pos.Latitude, 0.0000000000001)
		assert.InDelta(t, 3.72599833720439, pos.Longitude, 0.0000000000001)
	})

	t.Run("position wrong frames", func(t *testing.T) {
		t.Parallel()

		current := airbornPosition(t, "8D40621D58C386435CC412692AD6")

		_, err := current.DecodePosition(current)
		require.Error(t, err)
	})

	t.Run("raw data", func(t *testing.T) {
		for idx, fixtureElt := range []struct {
			input              string
			surveillanceStatus model.SurveillanceStatus
			singleAntennaFlag  bool
			encodedAltitude    uint16
			timeUTC            bool
			oddFrame           bool
			encodedLatitude    uint32
			encodedLongitude   uint32
			baro               bool
		}{
			{
				input:              "8D40621D58C386435CC412692AD6",
				surveillanceStatus: model.SurveillanceStatusNoCondition,
				singleAntennaFlag:  false,
				encodedAltitude:    0xc38,
				timeUTC:            false,
				oddFrame:           true,
				encodedLatitude:    0x121ae,
				encodedLongitude:   0xc412,
				baro:               true,
			},
			{
				input:              "8D40621D58C382D690C8AC2863A7",
				surveillanceStatus: model.SurveillanceStatusNoCondition,
				singleAntennaFlag:  false,
				encodedAltitude:    0xc38,
				timeUTC:            false,
				oddFrame:           false,
				encodedLatitude:    0x16b48,
				encodedLongitude:   0xc8ac,
				baro:               true,
			},
		} {
			fixture := fixtureElt

			t.Run(fmt.Sprintf("%d: %s", idx, fixture.input), func(t *testing.T) {
				t.Parallel()

				position := airbornPosition(t, fixture.input)

				assert.Equal(t, fixture.surveillanceStatus, position.SurveillanceStatus())
				assert.Equal(t, fixture.singleAntennaFlag, position.SingleAntennaFlag())
				assert.Equal(t, fixture.encodedAltitude, position.EncodedAltitude())
				assert.Equal(t, fixture.timeUTC, position.TimeUTC())
				assert.Equal(t, fixture.oddFrame, position.OddFrame())
				assert.Equal(t, fixture.encodedLatitude, position.EncodedLatitude())
				assert.Equal(t, fixture.encodedLongitude, position.EncodedLongitude())
				assert.Equal(t, fixture.baro, position.Baro())
			})
		}
	})
}
