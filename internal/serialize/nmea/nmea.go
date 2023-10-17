// Package nmea is the nmea serializer.
package nmea

import (
	"bytes"

	"github.com/landru29/dump1090/internal/dump"

	nmeaencoder "github.com/landru29/dump1090/internal/nmea"
)

type VesselType int

const (
	VesselTypeAircraft = iota
	VesselTypeHelicopter
)

// Serializer is the nmea serializer.
type Serializer struct {
	mmsiVessel VesselType
	mid        uint16
}

func New(mmsiVessel VesselType, mid uint16) *Serializer {
	return &Serializer{
		mmsiVessel: mmsiVessel,
		mid:        mid,
	}
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
		output := [][]byte{}
		for _, ac := range aircraft {
			if ac.Lon != 0 || ac.Lat != 0 {
				fields, err := nmeaencoder.Payload{
					MMSI:             s.MMSI(ac.Addr),
					Longitude:        ac.Lon,
					Latitude:         ac.Lat,
					SpeedOverGround:  float64(ac.Speed) / 10,
					PositionAccuracy: true,
					CourseOverGround: float64(ac.Track),
					TrueHeading:      uint16(ac.Track),
					NavigationStatus: nmeaencoder.NavigationStatusAground,
				}.Fields()
				if err != nil {
					return nil, err
				}
				output = append(output, []byte(fields.String()))
			}
		}

		if len(output) == 0 {
			return nil, nil
		}

		return bytes.Join(output, []byte("\n")), nil
	}

	return nil, nil
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/nmea"
}

func (s Serializer) MMSI(addr uint32) uint32 {
	out := uint32(s.mid%1000)*10000 + 10000000

	switch s.mmsiVessel {
	case VesselTypeAircraft:
		out += 1000
	case VesselTypeHelicopter:
		out += 5000
	}

	return out + addr%1000
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "nmea"
}
