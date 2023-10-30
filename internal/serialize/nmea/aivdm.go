package nmea

import "math"

type navigationStatus uint8

const (
	navigationStatusUnderWayUsingEngine navigationStatus = iota
	navigationStatusAtAnchor
	navigationStatusNotUnderCommand
	navigationStatusRestrictedManoeuverability
	navigationStatusConstrainedByHerDraught
	navigationStatusMoored
	navigationStatusAground
	navigationStatusEngagedInFishing
	navigationStatusUnderWaySailing
	navigationStatusReservedForFutureAmendmentOfNavigationalStatusForHSC
	navigationStatusReservedForFutureAmendmentOfNavigationalStatusForWIG
	navigationStatusPowerDrivenVesselTowingAstern
	navigationStatusPowerDrivenVesselPushingAheadOrTowingAlongside
	navigationStatusReservedForFutureUse
	navigationStatusAisSartIsActive
	navigationStatusUndefined
)

type maneuverIndicator uint8

const (
	notavailable      maneuverIndicator = iota //nolint: unused
	noSpecialManeuver                          //nolint: unused
	specialManeuver                            //nolint: unused
)

type radioChannel string

const (
	radioChannelA radioChannel = "A"
	radioChannelB radioChannel = "B"
)

// payload is the VDM / VDO payload
type payload struct {
	MMSI              uint32            // 8-37 (30)
	NavigationStatus  navigationStatus  // 38-41 (4)
	RateOfTurn        float64           // 42-49 (8)
	SpeedOverGround   float64           // 50-59 (10)
	PositionAccuracy  bool              // 60-60 (1)
	Longitude         float64           // 61-88 (28)
	Latitude          float64           // 89-115 (27)
	CourseOverGround  float64           // 116-127 (12)
	TrueHeading       uint16            // 128-136 (9)
	TimeStampSecond   uint8             // 137-142 (6)
	ManeuverIndicator maneuverIndicator // 143-144 (2)
	RaimFlag          bool              // 148-148 (1)
	RadioStatus       uint32            // 149-167 (19)
	RadioChannel      radioChannel
}

const (
	rateOfTurnRightMoreFiveDegPerMin = 710.0
	rateOfTurnLeftMoreFiveDegPerMin  = -710.0
	rateNoTurnInfo                   = -1.7e+308
)

func (p payload) Binary() (string, error) { //nolint: funlen,cyclop
	encoded := make([]uint8, 28)
	if _, err := payloadAddData(encoded, uint8(1), 0, 6); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, p.MMSI, 8, 30); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, uint8(p.NavigationStatus), 38, 4); err != nil {
		return "", err
	}

	rot := 128.0
	if p.RateOfTurn != rateNoTurnInfo {
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

	if _, err := payloadAddData(encoded, int8(rot), 42, 8); err != nil {
		return "", err
	}

	sog := uint16(math.Abs(p.SpeedOverGround * 10))
	if _, err := payloadAddData(encoded, sog, 50, 10); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, p.PositionAccuracy, 60, 1); err != nil {
		return "", err
	}

	lng := (int64(p.Longitude*600000.0+324000000) % 216000000) - 108000000
	if _, err := payloadAddData(encoded, lng, 61, 28); err != nil {
		return "", err
	}

	lat := int64(p.Latitude * 600000.0)
	if _, err := payloadAddData(encoded, lat, 89, 27); err != nil {
		return "", err
	}

	cog := uint16(math.Abs(p.CourseOverGround * 10))
	if _, err := payloadAddData(encoded, cog, 116, 12); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, p.TrueHeading, 128, 9); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, p.TimeStampSecond, 137, 6); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, uint8(p.ManeuverIndicator), 143, 2); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, p.RaimFlag, 148, 1); err != nil {
		return "", err
	}

	if _, err := payloadAddData(encoded, p.RadioStatus, 149, 19); err != nil {
		return "", err
	}

	return encodeBinaryPayload(encoded), nil
}

func (p payload) Fields() (fields, error) {
	binaryPayload, err := p.Binary()
	if err != nil {
		return nil, err
	}

	output := fields{
		[]byte("!AIVDM"),
		[]byte("1"), // fragment count
		[]byte("1"), // fragment number
		[]byte{},    // sequential message ID
		[]byte(p.RadioChannel),
		[]byte(binaryPayload),
		[]byte("0*"),
	}

	checkSum := output.checkSum()

	output[len(output)-1] = append(output[len(output)-1], []byte(checkSum)...)

	return output, nil
}
