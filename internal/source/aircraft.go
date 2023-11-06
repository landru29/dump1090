package source

import "time"

// Aircraft is an aircraft.
type Aircraft struct {
	Addr     uint32    `json:"icao"`     /* ICAO address */
	HexAddr  string    `json:"icao_hex"` /* Printable ICAO address */
	Flight   string    `json:"flight"`   /* Flight number */
	Altitude int       `json:"altitude"` /* Altitude */
	Speed    int       `json:"speed"`    /* Velocity computed from EW and NS components. */
	Track    int       `json:"track"`    /* Angle of flight. */
	Seen     time.Time `json:"-"`        /* Time at which the last packet was received. */
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"` /* Coordinated obtained from CPR encoded data. */
}
