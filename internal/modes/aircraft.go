package modes

import "github.com/landru29/dump1090/internal/source"

// Aircraft is an aircraft.
type Aircraft struct {
	Identification *Identification
}

// Unmarshal is the mode-s unmarshaler.
func (a *Aircraft) Unmarshal(data []byte) error {
	extendedSquitter := &ExtendedSquitter{}

	err := extendedSquitter.Unmarshal(data)
	if err != nil {
		return err
	}

	switch extendedSquitter.Type { //nolint: gocritic, exhaustive
	case MessageTypeAircraftIdentification:
		id, err := extendedSquitter.Identification()
		if err != nil {
			return err
		}

		a.Identification = id
	}

	return nil
}

// ToSource ...
func (a Aircraft) ToSource() source.Aircraft {
	return source.Aircraft{}
}
