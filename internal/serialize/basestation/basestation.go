// Package basestation is the BaseStation serializer.
package basestation

import (
	"bytes"
	"fmt"

	"github.com/landru29/dump1090/internal/dump"
)

// Serializer is the BaseStation serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	if ac == nil {
		return nil, nil
	}

	switch aircraft := ac.(type) {
	case dump.Aircraft:
		return s.Serialize([]*dump.Aircraft{&aircraft})

	case *dump.Aircraft:
		if aircraft == nil {
			return nil, nil
		}

		return []byte(message(*aircraft)), nil

	case []dump.Aircraft:
		output := [][]byte{}
		for _, ac := range aircraft {
			data, err := s.Serialize(ac)
			if err != nil {
				return nil, err
			}

			if len(data) != 0 {
				output = append(output, data)
			}
		}

		return bytes.Join(output, []byte("\n")), nil

	case []*dump.Aircraft:
		acs := []dump.Aircraft{}
		for _, ac := range aircraft {
			if ac != nil {
				acs = append(acs, *ac)
			}
		}

		return s.Serialize(acs)
	}

	return nil, nil
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/basestation"
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "base-station"
}

func message(aircraft dump.Aircraft) string { //nolint: cyclop
	alert := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Message.Alert()]

	emergency := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Message.Emergency()]

	ground := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Message.Ground()]

	spi := map[bool]int{
		false: 0,
		true:  -1,
	}[aircraft.Message.Indent()]

	switch {
	case aircraft.Message.DownlinkFormat == 0:
		return fmt.Sprintf("MSG,5,,,%s,,,,,,,%d,,,,,,,,,,", aircraft.Message.HexIdent(), aircraft.Message.Altitude)
	case aircraft.Message.DownlinkFormat == 4:
		return fmt.Sprintf("MSG,5,,,%s,,,,,,,%d,,,,,,,%d,%d,%d,%d", aircraft.Message.HexIdent(), aircraft.Message.Altitude, alert, emergency, spi, ground) //nolint: lll
	case aircraft.Message.DownlinkFormat == 5:
		return fmt.Sprintf("MSG,6,,,%s,,,,,,,,,,,,,%d,%d,%d,%d,%d", aircraft.Message.HexIdent(), aircraft.Message.Identity, alert, emergency, spi, ground) //nolint: lll
	case aircraft.Message.DownlinkFormat == 11:
		return fmt.Sprintf("MSG,8,,,%s,,,,,,,,,,,,,,,,,", aircraft.Message.HexIdent())
	case aircraft.Message.DownlinkFormat == 17 && aircraft.Message.Type == 4:
		return fmt.Sprintf("MSG,1,,,%s,,,,,,%s,,,,,,,,0,0,0,0", aircraft.Message.HexIdent(), aircraft.Message.Flight)
	case aircraft.Message.DownlinkFormat == 17 && aircraft.Message.Type >= 8 && aircraft.Message.Type <= 18:
		if aircraft.Lat == 0 && aircraft.Lon == 0 {
			return fmt.Sprintf("MSG,3,,,%s,,,,,,,%d,,,,,,,0,0,0,0", aircraft.Message.HexIdent(), aircraft.Message.Altitude)
		}

		return fmt.Sprintf("MSG,3,,,%s,,,,,,,%d,,,%1.5f,%1.5f,,,0,0,0,0", aircraft.Message.HexIdent(), aircraft.Message.Altitude, aircraft.Lat, aircraft.Lon) //nolint: lll
	case aircraft.Message.DownlinkFormat == 17 && aircraft.Message.Type == 19 && aircraft.Message.SubType == 1:
		verticalRate := int64(aircraft.Message.VertRate-1) * 64 * map[bool]int64{true: -1, false: 1}[aircraft.Message.VertRateNegative] //nolint: lll

		return fmt.Sprintf("MSG,4,,,%s,,,,,,,,%d,%d,,,%d,,0,0,0,0", aircraft.Message.HexIdent(), aircraft.Speed, aircraft.Track, verticalRate) //nolint: lll
	case aircraft.Message.DownlinkFormat == 21:
		return fmt.Sprintf("MSG,6,,,%s,,,,,,,,,,,,,%d,%d,%d,%d,%d", aircraft.Message.HexIdent(), aircraft.Message.Identity, alert, emergency, spi, ground) //nolint: lll
	}

	return ""
}
