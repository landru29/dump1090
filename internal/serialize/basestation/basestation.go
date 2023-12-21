// Package basestation is the BaseStation serializer.
package basestation

import (
	"bytes"
	"fmt"

	"github.com/landru29/dump1090/internal/model"
)

// Serializer is the BaseStation serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(planes ...any) ([]byte, error) {
	output := [][]byte{}

	for _, ac := range planes {
		switch aircraft := ac.(type) {
		case model.Aircraft:
			data, err := s.Serialize(&aircraft)
			if err != nil {
				return nil, err
			}

			output = append(output, data)

		case *model.Aircraft:
			if aircraft != nil {
				output = append(output, []byte(message(*aircraft)))
			}
		case []model.Aircraft:
			data, err := s.Serialize(model.UntypeArray(aircraft)...)
			if err != nil {
				return nil, err
			}

			output = append(output, data)
		case []*model.Aircraft:
			data, err := s.Serialize(model.UntypeArray(aircraft)...)
			if err != nil {
				return nil, err
			}

			output = append(output, data)
		}
	}

	return bytes.Join(output, []byte("\n")), nil
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/basestation"
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "base-station"
}

func message(aircraft model.Aircraft) string {
	alert := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Alert()]

	emergency := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Emergency()]

	ground := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Ground()]

	spi := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Indent()]

	switch {
	case aircraft.LastDownlinkFormat == 0:
		return fmt.Sprintf("MSG,5,,,%s,,,,,,,%d,,,,,,,,,,", aircraft.Identity, aircraft.Altitude)

	case aircraft.LastDownlinkFormat == 4:
		return fmt.Sprintf("MSG,5,,,%s,,,,,,,%d,,,,,,,%d,%d,%d,%d", aircraft.Identity, aircraft.Altitude, alert, emergency, spi, ground) //nolint: lll

	case aircraft.LastDownlinkFormat == 5:
		return fmt.Sprintf("MSG,6,,,%s,,,,,,,,,,,,,%d,%d,%d,%d,%d", aircraft.Identity, aircraft.Identity, alert, emergency, spi, ground) //nolint: lll

	case aircraft.LastDownlinkFormat == 11:
		return fmt.Sprintf("MSG,8,,,%s,,,,,,,,,,,,,,,,,", aircraft.Identity)

	case aircraft.LastDownlinkFormat == 17 && aircraft.LastType == 4:
		return fmt.Sprintf("MSG,1,,,%s,,,,,,%s,,,,,,,,0,0,0,0", aircraft.Identity, aircraft.Flight)

	case aircraft.LastDownlinkFormat == 17 && aircraft.LastType >= 8 && aircraft.LastType <= 18:
		if aircraft.Position == nil {
			return fmt.Sprintf("MSG,3,,,%s,,,,,,,%d,,,,,,,0,0,0,0", aircraft.Identity, aircraft.Altitude)
		}

		return fmt.Sprintf("MSG,3,,,%s,,,,,,,%d,,,%1.5f,%1.5f,,,0,0,0,0", aircraft.Identity, aircraft.Altitude, aircraft.Position.Latitude, aircraft.Position.Longitude) //nolint: lll

	case aircraft.LastDownlinkFormat == 17 && aircraft.LastType == 19 && aircraft.LastSubType == 1:
		return fmt.Sprintf("MSG,4,,,%s,,,,,,,,%d,%d,,,%d,,0,0,0,0", aircraft.Identity, aircraft.Speed, aircraft.Track, aircraft.VerticalRate) //nolint: lll

	case aircraft.LastDownlinkFormat == 21:
		return fmt.Sprintf("MSG,6,,,%s,,,,,,,,,,,,,%d,%d,%d,%d,%d", aircraft.Identity, aircraft.Identity, alert, emergency, spi, ground) //nolint: lll
	}

	return ""
}
