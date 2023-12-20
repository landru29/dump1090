package model

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
func (i ICAOAddr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + i.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *ICAOAddr) UnmarshalJSON(data []byte) error {
	if len(data) < minICAOsizeJSON || data[0] != '"' || data[len(data)-1] != '"' {
		return ErrWrongICAO
	}

	value, err := strconv.ParseUint(string(data[1:len(data)-1]), 16, 32)

	*i = ICAOAddr(value)

	return err
}

// String implements the Stringer interface.
func (i ICAOAddr) String() string {
	return strings.ToUpper(strconv.FormatUint(uint64(i), 16))
}
