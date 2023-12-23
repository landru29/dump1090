package model

const shortSquitterName = "short squitter"

// ShortSquitter is a short squitter message.
type ShortSquitter struct {
	ModeS
}

// AircraftAddress implements the Squitter interface.
func (s ShortSquitter) AircraftAddress() ICAOAddr {
	return ICAOAddr(0)
}

// Name implements the Squitter interface.
func (s ShortSquitter) Name() string {
	return shortSquitterName
}
