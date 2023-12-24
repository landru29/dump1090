package model

import (
	"github.com/landru29/dump1090/internal/compactposition"
	"github.com/landru29/dump1090/internal/errors"
)

const (
	errFrameOddEven errors.Error = "frames must be odd and even"
)

// Position is a GPS position.
type Position struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

// Positionner is a frame containing position informations.
type Positionner interface {
	// EncodedLatitude is the encoded latitude.
	EncodedLatitude() uint32

	// EncodedLongitude is the encoded longitude.
	EncodedLongitude() uint32

	// OddFrame defines if the frame is odd or even.
	OddFrame() bool
}

// DecodePosition decodes the current position with another frame.
func DecodePosition(current Positionner, other Positionner) (*Position, error) {
	if current.OddFrame() == other.OddFrame() {
		return nil, errFrameOddEven
	}

	if current.OddFrame() {
		latOdd, latEven := compactposition.DecodeLatitude(current.EncodedLatitude(), other.EncodedLatitude())
		lngOdd, _ := compactposition.DecodeLongitude(current.EncodedLongitude(), other.EncodedLongitude(), latOdd, latEven)

		return &Position{
			Latitude:  latOdd,
			Longitude: lngOdd,
		}, nil
	}

	latOdd, latEven := compactposition.DecodeLatitude(other.EncodedLatitude(), current.EncodedLatitude())
	_, lngEven := compactposition.DecodeLongitude(other.EncodedLongitude(), current.EncodedLongitude(), latOdd, latEven)

	return &Position{
		Latitude:  latEven,
		Longitude: lngEven,
	}, nil
}
