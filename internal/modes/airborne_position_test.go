package modes_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/landru29/dump1090/internal/modes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// +-------+-----+--------------+---+---+-------------------+-------------------+
// | TC    |     | ALT          | T | F | CPR-LAT           | CPR-LON           |
// +-------+-----+--------------+---+---+-------------------+-------------------+
// | 01011 | 000 | 110000111000 | 0 | 0 | 10110101101001000 | 01100100010101100 |
// | 01011 | 000 | 110000111000 | 0 | 1 | 10010000110101110 | 01100010000010010 |
// +-------+-----+--------------+---+---+-------------------+-------------------+

func TestAirbornPosition(t *testing.T) {
	t.Parallel()

	for idx, fixtureElt := range []struct {
		input    string
		expected modes.AirbornePosition
	}{
		{
			input: "8D40621D58C382D690C8AC2863A7",
			expected: modes.AirbornePosition{
				SurveillanceStatus: modes.SurveillanceStatusNoCondition,
				SingleAntennaFlag:  false,
				EncodedAltitude:    0xc38,
				TimeUTC:            false,
				OddFrame:           false,
				EncodedLatitude:    0x16b48,
				EncodedLongitude:   0xc8ac,
				Baro:               true,
			},
		},
		{
			input: "8D40621D58C386435CC412692AD6",
			expected: modes.AirbornePosition{
				SurveillanceStatus: modes.SurveillanceStatusNoCondition,
				SingleAntennaFlag:  false,
				EncodedAltitude:    0xc38,
				TimeUTC:            false,
				OddFrame:           true,
				EncodedLatitude:    0x121ae,
				EncodedLongitude:   0xc412,
				Baro:               true,
			},
		},
	} {
		fixture := fixtureElt

		t.Run(fmt.Sprintf("%d: %s", idx, fixture.input), func(t *testing.T) {
			t.Parallel()

			dataByte, err := hex.DecodeString(fixture.input)
			require.NoError(t, err)

			msg := modes.ExtendedSquitter{}

			require.NoError(t, msg.Unmarshal(dataByte))

			airbornePosition, err := msg.AirbornePosition()
			require.NoError(t, err)

			// reset time for comparison.
			airbornePosition.Time = time.Time{}

			assert.Equal(t, fixture.expected, *airbornePosition)
		})
	}
}
