package nmea

import (
	"bytes"

	"github.com/landru29/dump1090/internal/dump"
)

const (
	speedOverGroundScale = 10
)

// VesselType is a type of vessel.
type VesselType int

const (
	// VesselTypeAircraft is an aircraft.
	VesselTypeAircraft = iota
	// VesselTypeHelicopter is a helicopter.
	VesselTypeHelicopter
)

// Serializer is the nmea serializer.
type Serializer struct {
	mmsiVessel VesselType
	mid        uint16
}

// New is a new NMEA serializer.
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
				fields, err := payload{
					MMSI:             s.MMSI(ac.Addr),
					Longitude:        ac.Lon,
					Latitude:         ac.Lat,
					SpeedOverGround:  float64(ac.Speed) / speedOverGroundScale,
					PositionAccuracy: true,
					CourseOverGround: float64(ac.Track),
					TrueHeading:      uint16(ac.Track),
					NavigationStatus: navigationStatusAground,
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

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/nmea"
}

// MMSI ...
func (s Serializer) MMSI(addr uint32) uint32 {
	out := uint32(s.mid%1000)*10000 + 10000000 //nolint: gomnd

	switch s.mmsiVessel {
	case VesselTypeAircraft:
		out += 1000
	case VesselTypeHelicopter:
		out += 5000
	}

	return out + addr%1000
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "nmea"
}
