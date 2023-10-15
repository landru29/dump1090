package nmea

import "math"

type NavigationStatus uint8

const (
	NavigationStatusUnderWayUsingEngine NavigationStatus = iota
	NavigationStatusAtAnchor
	NavigationStatusNotUnderCommand
	NavigationStatusRestrictedManoeuverability
	NavigationStatusConstrainedByHerDraught
	NavigationStatusMoored
	NavigationStatusAground
	NavigationStatusEngagedInFishing
	NavigationStatusUnderWaySailing
	NavigationStatusReservedForFutureAmendmentOfNavigationalStatusForHSC
	NavigationStatusReservedForFutureAmendmentµOfNavigationalStatusForWIG
	NavigationStatusPowerDrivenVesselTowingAstern
	NavigationStatusPowerDrivenVesselPushingAheadOrTowingAlongside
	NavigationStatusReservedForFutureUse
	NavigationStatusAisSartIsActive
	NavigationStatusUndefined
)

type ManeuverIndicator uint8

const (
	Notavailable ManeuverIndicator = iota
	NoSpecialManeuver
	SpecialManeuver
)

type RadioChannel string

const (
	RadioChannelA RadioChannel = "A"
	RadioChannelB RadioChannel = "B"
)

// Payload is the VDM / VDO payload
type Payload struct {
	MMSI              uint32            // 8-37 (30)
	NavigationStatus  NavigationStatus  // 38-41 (4)
	RateOfTurn        float64           // 42-49 (8)
	SpeedOverGround   float64           // 50-59 (10)
	PositionAccuracy  bool              // 60-60 (1)
	Longitude         float64           // 61-88 (28)
	Latitude          float64           // 89-115 (27)
	CourseOverGround  float64           // 116-127 (12)
	TrueHeading       uint16            // 128-136 (9)
	TimeStampSecond   uint8             // 137-142 (6)
	ManeuverIndicator ManeuverIndicator // 143-144 (2)
	RaimFlag          bool              // 148-148 (1)
	RadioStatus       uint32            // 149-167 (19)
	RadioChannel      RadioChannel
}

const (
	RateOfTurnRightMoreFiveDegPerMin = 710.0
	RateOfTurnLeftMoreFiveDegPerMin  = -710.0
	RateNoTurnInfo                   = -1.7e+308
)

func (p Payload) Binary() (string, error) {
	encoded := make([]uint8, 28)
	if _, err := PayloadAddData(encoded, uint8(1), 0, 6); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, p.MMSI, 8, 30); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, uint8(p.NavigationStatus), 38, 4); err != nil {
		return "", err
	}

	rot := 128.0
	if p.RateOfTurn != RateNoTurnInfo {
		rot = 4.733 * math.Sqrt(math.Abs(p.RateOfTurn))
		if p.RateOfTurn < 0 {
			rot = -rot
		}
		if rot > 126 {
			rot = 127
		}

		if rot < -126 {
			rot = -127
		}
	}

	if _, err := PayloadAddData(encoded, int8(rot), 42, 8); err != nil {
		return "", err
	}

	sog := uint16(math.Abs(p.SpeedOverGround * 10))
	if _, err := PayloadAddData(encoded, sog, 50, 10); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, p.PositionAccuracy, 60, 1); err != nil {
		return "", err
	}

	lng := (int64(p.Longitude*600000.0+324000000) % 216000000) - 108000000
	if _, err := PayloadAddData(encoded, lng, 61, 28); err != nil {
		return "", err
	}

	lat := int64(p.Latitude * 600000.0)
	if _, err := PayloadAddData(encoded, lat, 89, 27); err != nil {
		return "", err
	}

	cog := uint16(math.Abs(p.CourseOverGround * 10))
	if _, err := PayloadAddData(encoded, cog, 116, 12); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, p.TrueHeading, 128, 9); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, p.TimeStampSecond, 137, 6); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, uint8(p.ManeuverIndicator), 143, 2); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, p.RaimFlag, 148, 1); err != nil {
		return "", err
	}

	if _, err := PayloadAddData(encoded, p.RadioStatus, 149, 19); err != nil {
		return "", err
	}

	return EncodeBinaryPayload(encoded), nil
}

func (p Payload) Fields() (Fields, error) {
	binaryPayload, err := p.Binary()
	if err != nil {
		return nil, err
	}

	output := Fields{
		[]byte("!AIVDM"),
		[]byte("1"), // fragment count
		[]byte("1"), // fragment number
		[]byte{},    // sequential message ID
		[]byte(p.RadioChannel),
		[]byte(binaryPayload),
		[]byte("0*"),
	}

	checkSum := output.CheckSum()

	output[len(output)-1] = append(output[len(output)-1], []byte(checkSum)...)

	return output, nil
}
