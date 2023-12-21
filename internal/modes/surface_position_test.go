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

// +-------+---------+---+---------+---+---+-------------------+-------------------+
// | TC    | MOV     | S | TRK     | T | F | CPR-LAT           | CPR-LON           |
// +-------+---------+-------------+---+---+-------------------+-------------------+
// | 00111 | 0101010 | 1 | 0110010 | 0 | 0 | 11100001110011001 | 11100100011001101 |
// | 00111 | 0101000 | 1 | 0100011 | 0 | 1 | 01001100100011111 | 11010111010111101 |
// +-------+---------+---+---------+---+---+-------------------+-------------------+

func TestSurfacePosition(t *testing.T) {
	t.Parallel()

	for idx, fixtureElt := range []struct {
		input    string
		expected modes.SurfacePosition
	}{
		{
			input: "8C4841753AAB238733C8CD4020B1",
			expected: modes.SurfacePosition{
				GroundTrack:      140.625,
				TimeUTC:          false,
				OddFrame:         false,
				EncodedLatitude:  0x1c399,
				EncodedLongitude: 0x1c8cd,
				Speed:            18,
			},
		},
		{
			input: "8C4841753A8A35323FAEBDAC702D",
			expected: modes.SurfacePosition{
				GroundTrack:      98.4375,
				TimeUTC:          false,
				OddFrame:         true,
				EncodedLatitude:  0x991f,
				EncodedLongitude: 0x1aebd,
				Speed:            16,
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

			surfacePosition, err := msg.SurfacePosition()
			require.NoError(t, err)

			// reset time for comparison.
			surfacePosition.Time = time.Time{}

			assert.Equal(t, fixture.expected, *surfacePosition)
		})
	}
}
