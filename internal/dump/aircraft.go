package dump

/*
#cgo LDFLAGS: -lrtlsdr -lm
#include "dump1090.h"
*/
import "C"

import (
	"fmt"
	"strings"
	"time"
)

// Aircraft is an aircraft.
type Aircraft struct {
	Addr     uint32    `json:"icao"`     /* ICAO address */
	HexAddr  string    `json:"icao_hex"` /* Printable ICAO address */
	Flight   string    `json:"flight"`   /* Flight number */
	Altitude int       `json:"altitude"` /* Altitude */
	Speed    int       `json:"speed"`    /* Velocity computed from EW and NS components. */
	Track    int       `json:"track"`    /* Angle of flight. */
	Seen     time.Time `json:"-"`        /* Time at which the last packet was received. */
	Messages int64     `json:"-"`        /* Number of Mode S messages received. */
	/* Encoded latitude and longitude as extracted by odd and even
	 * CPR encoded messages. */
	OddCPRlat  int     `json:"-"`
	OddCPRlon  int     `json:"-"`
	EvenCPRlat int     `json:"-"`
	EvenCPRlon int     `json:"-"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"` /* Coordinated obtained from CPR encoded data. */

	OddCPRtime  time.Time `json:"-"`
	EvenCPRtime time.Time `json:"-"`
	CountryCode string    `json:"country_code"`

	Message Message `json:"-"`
}

func (a Aircraft) String() string {
	return fmt.Sprintf("hex:      %s\nflight:   %s\naltitude: %d\nspeed:    %d\ntrack:    %d\nlat:      %f\nlon:      %f\nseen:     %s\n", //nolint: lll
		a.HexAddr, a.Flight, a.Altitude, a.Speed, a.Track, a.Lat, a.Lon, a.Seen)
}

func newAircraft(aircraft *C.aircraft, msg *C.modesMessage) Aircraft {
	return Aircraft{
		Addr:        uint32(aircraft.addr),
		HexAddr:     C.GoString(&aircraft.hexaddr[0]),
		Flight:      strings.Trim(C.GoString(&aircraft.flight[0]), " "),
		Altitude:    int(aircraft.altitude),
		Speed:       int(aircraft.speed),
		Track:       int(aircraft.track),
		Seen:        time.Unix(int64(aircraft.seen), 0),
		Messages:    int64(aircraft.messages),
		OddCPRlat:   int(aircraft.odd_cprlat),
		OddCPRlon:   int(aircraft.odd_cprlon),
		EvenCPRlat:  int(aircraft.even_cprlat),
		EvenCPRlon:  int(aircraft.even_cprlon),
		Lat:         float64(aircraft.lat),
		Lon:         float64(aircraft.lon),
		OddCPRtime:  time.Unix(int64(aircraft.odd_cprtime)/1000, (int64(aircraft.odd_cprtime)%1000)*1000000),   //nolint: gomnd,lll
		EvenCPRtime: time.Unix(int64(aircraft.even_cprtime)/1000, (int64(aircraft.even_cprtime)%1000)*1000000), //nolint: gomnd,lll
		Message:     newMessage(msg),
	}
}
