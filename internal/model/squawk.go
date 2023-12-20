package model

import (
	"math"
	"strconv"
	"strings"

	"github.com/landru29/dump1090/internal/errors"
)

const (
	// ErrWrongSquawk is when the Squawk is not valid.
	ErrWrongSquawk errors.Error = "squawk digits are lower or equal to 7"

	maxSquawkDigit = 7

	maxSquawk Squawk = 7777

	// SquawkHijacker is the squawk to set when being under hijacker.
	SquawkHijacker Squawk = 7500

	// SquawkRadioFailure is the squawk to set when the onboard radio is broken.
	SquawkRadioFailure Squawk = 7600

	// SquawkMayday is the squawk to set when panic onboard.
	SquawkMayday Squawk = 7700
)

// Squawk is the transponder code.
type Squawk uint16

// MarshalJSON implements the json.Marshaler interface.
func (s Squawk) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Squawk) UnmarshalJSON(data []byte) error {
	out, err := strconv.ParseUint(string(data), 10, 16)
	if err != nil {
		return err
	}

	*s = Squawk(out)

	if *s > maxSquawk {
		return ErrWrongSquawk
	}

	for idx := 0; idx < 4; idx++ {
		if s.DigitAt(idx) > maxSquawkDigit {
			return ErrWrongSquawk
		}
	}

	return nil
}

// String implements the Stringer interface.
func (s Squawk) String() string {
	return strings.ToUpper(strconv.FormatUint(uint64(s), 10))
}

// DigitAt gets the digit at a specific position.
func (s Squawk) DigitAt(position int) uint8 {
	return uint8((uint16(s) / uint16(math.Pow10(position))) % 10) //nolint: gomnd
}
