package modes

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
