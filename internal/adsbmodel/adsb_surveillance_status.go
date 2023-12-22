package adsbmodel

// SurveillanceStatus is the surveillance status.
type SurveillanceStatus byte

const (
	// SurveillanceStatusNoCondition is no condition.
	SurveillanceStatusNoCondition SurveillanceStatus = iota

	// SurveillanceStatusPermanentAlert is permanent alert.
	SurveillanceStatusPermanentAlert

	// SurveillanceStatusTemporaryAlert is temporary alert.
	SurveillanceStatusTemporaryAlert

	// SurveillanceStatusSpiCondition is SPI condition.
	SurveillanceStatusSpiCondition
)
