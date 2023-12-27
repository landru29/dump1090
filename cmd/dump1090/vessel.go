package main

import (
	"fmt"

	"github.com/landru29/dump1090/internal/serialize/nmea"
)

type vessel nmea.VesselType

// String implements the pflag.Value interface.
func (v vessel) String() string {
	return map[nmea.VesselType]string{
		nmea.VesselTypeAircraft:   "aircraft",
		nmea.VesselTypeHelicopter: "helicopter",
	}[nmea.VesselType(v)]
}

// Set implements the pflag.Value interface.
func (v *vessel) Set(str string) error {
	vesselType, ok := map[string]nmea.VesselType{
		"aircraft":   nmea.VesselTypeAircraft,
		"helicopter": nmea.VesselTypeHelicopter,
	}[str]
	if !ok {
		return fmt.Errorf("unknow vessel type %s", str)
	}

	*v = vessel(vesselType)

	return nil
}

// Type implements the pflag.Value interface.
func (v vessel) Type() string {
	return "vessel type"
}
