package compactposition_test

import (
	"testing"

	"github.com/landru29/dump1090/internal/compactposition"
	"github.com/stretchr/testify/assert"
)

func TestDecodeLatitude(t *testing.T) {
	t.Parallel()

	latOdd, latEven := compactposition.DecodeLatitude(74158, 93000)
	assert.InDelta(t, 52.26578017412606, latOdd, 0.0000000000001)
	assert.InDelta(t, 52.25720214843750, latEven, 0.0000000000001)
}

func TestDecodeLongitude(t *testing.T) {
	t.Parallel()

	latOdd, latEven := compactposition.DecodeLatitude(74158, 93000)
	assert.InDelta(t, 52.26578017412606, latOdd, 0.0000000000001)
	assert.InDelta(t, 52.25720214843750, latEven, 0.0000000000001)

	lngOdd, lngEven := compactposition.DecodeLongitude(50194, 51372, latOdd, latEven)
	assert.InDelta(t, 3.72599833720439, lngOdd, 0.0000000000001)
	assert.InDelta(t, 3.91937255859375, lngEven, 0.0000000000001)
}
