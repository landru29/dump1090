// Package compactposition decodes CPR latitudes and longitudes.
package compactposition

import "math"

const numberOfLatitudeZones = 15

// longitudeZoneNumber yields the number of longitude zones between 1 and 59.
func longitudeZoneNumber(latitude float64) int {
	switch {
	case latitude == 0:
		return 59
	case latitude == 87:
		return 2
	case latitude == -87:
		return 2
	case latitude > 87:
		return 1
	case latitude < -87:
		return 1
	default:
		return int(math.Floor(
			2 * math.Pi /
				math.Acos(
					1-
						(1-math.Cos(math.Pi/(2*float64(numberOfLatitudeZones))))/
							math.Pow(math.Cos(math.Pi*latitude/180.0), 2),
				),
		))
	}
}

func latitudeZoneSize(odd bool) float64 {
	if odd {
		return 360.0 / float64(4*numberOfLatitudeZones-1)
	}

	return 360.0 / float64(4*numberOfLatitudeZones)
}

func longitudeZoneCount(odd bool, latitude float64) float64 {
	if odd {
		return math.Max(float64(longitudeZoneNumber(latitude-1)), 1)
	}

	return math.Max(float64(longitudeZoneNumber(latitude)), 1)
}

func longitudeZoneSize(odd bool, latitude float64) float64 {
	return 360.0 / longitudeZoneCount(odd, latitude)
}

// DecodeLatitude decodes both odd and even latitudes.
func DecodeLatitude(latCprOdd uint32, latCprEven uint32) (float64, float64) {
	latitudeZoneIndex := math.Floor(59.0*float64(latCprEven)/131072.0 - 60.0*float64(latCprOdd)/131072.0 + 0.5)

	latitudeEven := latitudeZoneSize(false) * (math.Mod(latitudeZoneIndex, 60.0) + float64(latCprEven)/131072.0)
	latitudeOdd := latitudeZoneSize(true) * (math.Mod(latitudeZoneIndex, 59.0) + float64(latCprOdd)/131072.0)

	if latitudeEven > 270 {
		latitudeEven -= 360.0
	}

	if latitudeOdd > 270 {
		latitudeOdd -= 360.0
	}

	return latitudeOdd, latitudeEven
}

// DecodeLongitude decodes both odd and even longitudes.
func DecodeLongitude(
	lngCprOdd uint32, lngCprEven uint32,
	latitudeOdd float64, latitudeEven float64,
) (float64, float64) {
	lngZoneNumOdd := float64(longitudeZoneNumber(latitudeOdd))
	lngZoneNumEven := float64(longitudeZoneNumber(latitudeEven))

	longitudeIndex := math.Floor(
		(float64(lngCprEven)/131072.0)*(lngZoneNumEven-1) -
			(float64(lngCprOdd)/131072.0)*lngZoneNumOdd +
			0.5,
	)

	longitudeEven := longitudeZoneSize(false, latitudeEven) *
		(math.Mod(longitudeIndex, longitudeZoneCount(false, latitudeEven)) + float64(lngCprEven)/131072.0)
	longitudeOdd := longitudeZoneSize(true, latitudeOdd) *
		(math.Mod(longitudeIndex, longitudeZoneCount(true, latitudeOdd)) + float64(lngCprOdd)/131072.0)

	return longitudeOdd, longitudeEven
}
