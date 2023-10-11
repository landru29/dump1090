package nmea_test

import (
	"testing"

	"github.com/landru29/dump1090/internal/nmea"
	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	assert.Equal(t,
		"!AIVDM,1,1,,A,15RTgt0PAso;90TKcjM8h6g208CQ,0*4A",
		nmea.Payload{
			MMSI:              371798000,
			NavigationStatus:  nmea.NavigationStatusUnderWaySailing,
			RateOfTurn:        nmea.RateOfTurnLeftMoreFiveDegPerMin,
			SpeedOverGround:   12.3,
			PositionAccuracy:  true,
			Longitude:         -123.39538333333333,
			Latitude:          48.38163333333333,
			CourseOverGround:  224,
			TrueHeading:       215,
			TimeStampSecond:   33,
			ManeuverIndicator: 0,
			RaimFlag:          false,
			RadioStatus:       34017,
			RadioChannel:      nmea.RadioChannelA,
		}.Fields().String(),
	)
}
