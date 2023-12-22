package adsbmodel

const (
	// TypeCodeAircraftIdentification is the aircraft identification Type Code.
	TypeCodeAircraftIdentification TypeCode = 4
	// TypeCodeSurfacePosition is the surface position Type Code.
	TypeCodeSurfacePosition TypeCode = 8
	// TypeCodeAirbornePositionBaroAltitude is the airborne position (w/Baro Altitude) Type Code.
	TypeCodeAirbornePositionBaroAltitude TypeCode = 18
	// TypeCodeAirborneVelocities is the airborne velocities Type Code.
	TypeCodeAirborneVelocities TypeCode = 19
	// TypeCodeAirbornePositionGNSSHeight is the airborne position (w/GNSS Height) Type Code.
	TypeCodeAirbornePositionGNSSHeight TypeCode = 22
	// TypeCodeReserved is the reserved Type Code.
	TypeCodeReserved TypeCode = 27
	// TypeCodeAircraftStatus is the aircraft status Type Code.
	TypeCodeAircraftStatus TypeCode = 28
	// TypeCodeTargetStateAndStatusInformation is the target state and status information Type Code.
	TypeCodeTargetStateAndStatusInformation TypeCode = 29
	// TypeCodeAircraftOperationStatus is the aircraft operation status Type Code.
	TypeCodeAircraftOperationStatus TypeCode = 31
)

// TypeCode is the type of the exetended squitter.
type TypeCode byte

// Code is the normalized Type Code.
func (c TypeCode) Code() TypeCode {
	switch {
	case c == 0:
		return 0
	case c <= TypeCodeAircraftIdentification:
		return TypeCodeAircraftIdentification
	case c <= TypeCodeSurfacePosition:
		return TypeCodeSurfacePosition
	case c <= TypeCodeAirbornePositionBaroAltitude:
		return TypeCodeAirbornePositionBaroAltitude
	case c <= TypeCodeAirborneVelocities:
		return TypeCodeAirborneVelocities
	case c <= TypeCodeAirbornePositionGNSSHeight:
		return TypeCodeAirbornePositionGNSSHeight
	case c <= TypeCodeReserved:
		return TypeCodeReserved
	case c <= TypeCodeAircraftStatus:
		return TypeCodeAircraftStatus
	case c <= TypeCodeTargetStateAndStatusInformation:
		return TypeCodeTargetStateAndStatusInformation
	case c <= TypeCodeAircraftOperationStatus:
		return TypeCodeAircraftOperationStatus
	}

	return 0
}
