package model_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAirbornePosition(t *testing.T) {
	t.Parallel()

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

			assert.Equal(t, "airborne position", msg.Name())

			position, ok := msg.(model.AirbornePosition)
			assert.True(t, ok)

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
}
