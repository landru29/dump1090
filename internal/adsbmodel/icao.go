package adsbmodel

import (
	"strconv"
	"strings"

	"github.com/landru29/dump1090/internal/errors"
)

const (
	// ErrWrongICAO is when trying to unmarshal the wrong value.
	ErrWrongICAO errors.Error = "wrong ICAO address"

	minICAOsizeJSON = 2
)

// ICAOAddr is the ICAO aircraft address.
type ICAOAddr uint32

// MarshalJSON implements the json.Marshaler interface.
func (a ICAOAddr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *ICAOAddr) UnmarshalJSON(data []byte) error {
	if len(data) < minICAOsizeJSON || data[0] != '"' || data[len(data)-1] != '"' {
		return ErrWrongICAO
	}

	value, err := strconv.ParseUint(string(data[1:len(data)-1]), 16, 32)

	*a = ICAOAddr(value)

	return err
}

// String implements the Stringer interface.
func (a ICAOAddr) String() string {
	return strings.ToUpper(strconv.FormatUint(uint64(a), 16))
}
