// Package nmea is the nmea serializer.
package nmea

import (
	"github.com/landru29/dump1090/internal/dump"

	nmeaencoder "github.com/landru29/dump1090/internal/nmea"
)

// Serializer is the nmea serializer.
type Serializer struct{}

func New() *Serializer {
	return &Serializer{}
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	if ac == nil {
		return nil, nil
	}

	switch aircraft := ac.(type) {
	case dump.Aircraft:
		return s.Serialize([]*dump.Aircraft{&aircraft})
	case *dump.Aircraft:
		return s.Serialize([]*dump.Aircraft{aircraft})
	case []dump.Aircraft:
		out := make([]*dump.Aircraft, len(aircraft))
		for idx := range aircraft {
			out[idx] = &aircraft[idx]
		}
		return s.Serialize(out)
	case []*dump.Aircraft:
		return []byte(nmeaencoder.Payload{}.Fields().String()), nil
	}

	return nil, nil
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "text/plain"
}
