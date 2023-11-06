package modes

import "github.com/landru29/dump1090/internal/source"

type Aircraft struct {
	Identification *Identification
}

func (a *Aircraft) Unmarshal(data []byte) error {
	extendedSquitter := &ExtendedSquitter{}

	err := extendedSquitter.Unmarshal(data)
	if err != nil {
		return err
	}

	switch extendedSquitter.Type {
	case MessageTypeAircraftIdentification:
		id, err := extendedSquitter.Identification()
		if err != nil {
			return err
		}

		a.Identification = id
	}

	return nil
}

func (a Aircraft) ToSource() source.Aircraft {
	return source.Aircraft{}
}
