package model

import (
	"fmt"
	"strings"
	"time"
)

// verticalRate := int64(aircraft.Message.VertRate-1)
//               * 64
//               * map[bool]int64{true: -1, false: 1}[aircraft.Message.VertRateNegative]

// Aircraft is an aircraft description.
type Aircraft struct {
	IcaoAddress        ICAOAddr  `json:"icao_address"`
	Altitude           int64     `json:"altitude"`
	Position           *Position `json:"position,omitempty"`
	Flight             string    `json:"flight"`             /* Flight number */
	Addr               uint32    `json:"icao"`               /* ICAO address */
	Speed              int       `json:"speed"`              /* Velocity computed from EW and NS components. */
	Track              int       `json:"track"`              /* Angle of flight. */
	Identity           Squawk    `json:"identity"`           /* 13 bits identity (from transponder). */
	LastUpdate         time.Time `json:"last_update"`        /* Time at which the last packet was received. */
	LastFlightStatus   int       `json:"last_flight_status"` /* Flight status for DF4,5,20,21 */
	LastDownlinkFormat int       `json:"downlink_format"`    /* Downlink format # */
	LastType           int
	LastSubType        int
	VerticalRate       int64 `json:"vertical_rate"`
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