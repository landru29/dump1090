package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/landru29/dump1090/internal/source"
)

// verticalRate := int64(aircraft.Message.VertRate-1)
//               * 64
//               * map[bool]int64{true: -1, false: 1}[aircraft.Message.VertRateNegative]

// Aircraft is an aircraft description.
type Aircraft struct {
	Identification *Identification

	IcaoAddress        ICAOAddr  `json:"icaoAddress"`
	Altitude           int64     `json:"altitude"`
	Position           *Position `json:"position,omitempty"`
	Flight             string    `json:"flight"`           /* Flight number */
	Addr               uint32    `json:"icao"`             /* ICAO address */
	Speed              int       `json:"speed"`            /* Velocity computed from EW and NS components. */
	Track              int       `json:"track"`            /* Angle of flight. */
	Identity           Squawk    `json:"identity"`         /* 13 bits identity (from transponder). */
	LastUpdate         time.Time `json:"lastUpdate"`       /* Time at which the last packet was received. */
	LastFlightStatus   int       `json:"lastFlightStatus"` /* Flight status for DF4,5,20,21 */
	LastDownlinkFormat int       `json:"downlinkFormat"`   /* Downlink format # */
	LastType           int
	LastSubType        int
	VerticalRate       int64 `json:"verticalRate"`
}

// String implements the Stringer interface.
func (a Aircraft) String() string {
	return strings.Join(
		[]string{
			fmt.Sprintf("hex:      %X", a.Addr),
			fmt.Sprintf("flight:   %s", a.Flight),
			fmt.Sprintf("altitude: %d", a.Altitude),
			fmt.Sprintf("speed:    %d", a.Speed),
			fmt.Sprintf("track:    %d", a.Track),
			fmt.Sprintf("lat:      %f", a.Position.Latitude),
			fmt.Sprintf("lng:      %f", a.Position.Longitude),
			fmt.Sprintf("seen:     %s", a.LastUpdate),
		},
		"\n",
	)
}

// Emergency ...
func (a Aircraft) Emergency() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.Identity == SquawkHijacker || a.Identity == SquawkRadioFailure || a.Identity == SquawkMayday)
}

// Alert ...
func (a Aircraft) Alert() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.LastFlightStatus == 2 || a.LastFlightStatus == 3 || a.LastFlightStatus == 4)
}

// Ground ...
func (a Aircraft) Ground() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.LastFlightStatus == 1 || a.LastFlightStatus == 3)
}

// Indent ...
func (a Aircraft) Indent() bool {
	return (a.LastDownlinkFormat == 4 || a.LastDownlinkFormat == 5 || a.LastDownlinkFormat == 21) &&
		(a.LastFlightStatus == 4 || a.LastFlightStatus == 5)
}

// UnmarshalModeS is the mode-s unmarshaler.
// func (a *Aircraft) UnmarshalModeS(data []byte) error {
// 	extendedSquitter := &ExtendedSquitter{}

// 	err := extendedSquitter.UnmarshalModeS(data)
// 	if err != nil {
// 		return err
// 	}

// 	switch extendedSquitter.Type { //nolint: gocritic, exhaustive
// 	case MessageTypeAircraftIdentification:
// 		id, err := extendedSquitter.Identification()
// 		if err != nil {
// 			return err
// 		}

// 		a.Identification = id
// 	}

// 	return nil
// }

// ToSource ...
func (a Aircraft) ToSource() source.Aircraft {
	return source.Aircraft{}
}
