package model

// Position is a GPS position.
type Position struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}
