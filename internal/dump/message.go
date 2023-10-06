package dump

/*
#cgo LDFLAGS: -lrtlsdr -lm
#include "dump1090.h"
*/
import "C"

// l = snprintf(p,buflen,
// 	"{\"hex\":\"%s\", \"flight\":\"%s\", \"lat\":%f, "
// 	"\"lon\":%f, \"altitude\":%d, \"track\":%d, "
// 	"\"speed\":%d},\n",
// 	a->hexaddr, a->flight, a->lat, a->lon, a->altitude, a->track,
// 	a->speed);

type Message struct {
	/* Generic fields */
	MsgBits        int    /* Number of bits in message */
	DownlinkFormat int    /* Downlink format # */
	CRCok          int    /* True if CRC was valid */
	CRC            uint32 /* Message CRC */
	Errorbit       int    /* Bit corrected. -1 if no bit corrected. */
	Aa1            int    /* ICAO Address bytes 1 */
	Aa2            int    /* ICAO Address bytes 2 */
	Aa3            int    /* ICAO Address bytes 3 */
	PhaseCorrected int    /* True if phase correction was applied. */

	/* DF 11 */
	Capabilities int /* Responder capabilities. */

	/* DF 17 */
	Type           int /* Extended squitter message type. */
	SubType        int /* Extended squitter message subtype. */
	HeadingIsValid int
	Heading        int
	AircraftType   int
	Fflag          int    /* 1 = Odd, 0 = Even CPR message. */
	Tflag          int    /* UTC synchronized? */
	RawLatitude    int    /* Non decoded latitude */
	RawLongitude   int    /* Non decoded longitude */
	Flight         string /* 8 chars flight number. */
	EwDir          int    /* 0 = East, 1 = West. */
	EwVelocity     int    /* E/W velocity. */
	NsDir          int    /* 0 = North, 1 = South. */
	NsVelocity     int    /* N/S velocity. */
	VertRateSource int    /* Vertical rate source. */
	VertRateSign   int    /* Vertical rate sign. */
	VertRate       int    /* Vertical rate. */
	Velocity       int    /* Computed from EW and NS velocity. */

	/* DF4, DF5, DF20, DF21 */
	FlightStatus int /* Flight status for DF4,5,20,21 */
	Dr           int /* Request extraction of downlink request. */
	Um           int /* Request extraction of downlink request. */
	Identity     int /* 13 bits identity (Squawk). */

	/* Fields used by multiple message types. */
	Altitude int
	Unit     int
}

func newMessage(message *C.modesMessage) Message {
	return Message{
		MsgBits:        int(message.msgbits),
		DownlinkFormat: int(message.msgtype),
		CRCok:          int(message.crcok),
		CRC:            uint32(message.crc),
		Errorbit:       int(message.errorbit),
		Aa1:            int(message.aa1),
		Aa2:            int(message.aa2),
		Aa3:            int(message.aa3),
		PhaseCorrected: int(message.phase_corrected),
		Capabilities:   int(message.ca),
		Type:           int(message.metype),
		SubType:        int(message.mesub),
		HeadingIsValid: int(message.heading_is_valid),
		Heading:        int(message.heading),
		AircraftType:   int(message.aircraft_type),
		Fflag:          int(message.fflag),
		Tflag:          int(message.tflag),
		RawLatitude:    int(message.raw_latitude),
		RawLongitude:   int(message.raw_longitude),
		Flight:         C.GoString(&message.flight[0]),
		EwDir:          int(message.ew_dir),
		EwVelocity:     int(message.ew_velocity),
		NsDir:          int(message.ns_dir),
		NsVelocity:     int(message.ns_velocity),
		VertRateSource: int(message.vert_rate_source),
		VertRateSign:   int(message.vert_rate_sign),
		VertRate:       int(message.vert_rate),
		Velocity:       int(message.velocity),
		FlightStatus:   int(message.fs),
		Dr:             int(message.dr),
		Um:             int(message.um),
		Identity:       int(message.identity),
		Altitude:       int(message.altitude),
		Unit:           int(message.unit),
	}
}
