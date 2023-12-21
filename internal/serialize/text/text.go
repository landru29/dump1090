// Package text is the text serializer.
package text

import (
	"bytes"

	"github.com/landru29/dump1090/internal/model"
)

// Serializer is the text serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(planes ...any) ([]byte, error) {
	output := [][]byte{}

	for _, ac := range planes {
		switch aircraft := ac.(type) {
		case model.Aircraft:
			data, err := s.Serialize(&aircraft)
			if err != nil {
				return nil, err
			}

			output = append(output, data)

		case *model.Aircraft:
			if aircraft != nil {
				output = append(output, []byte(aircraft.String()))
			}
		case []model.Aircraft:
			data, err := s.Serialize(model.UntypeArray(aircraft)...)
			if err != nil {
				return nil, err
			}

			output = append(output, data)
		case []*model.Aircraft:
			data, err := s.Serialize(model.UntypeArray(aircraft)...)
			if err != nil {
				return nil, err
			}

			output = append(output, data)
		}
	}

	return bytes.Join(output, []byte("\n")), nil
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "text/plain"
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "text"
}
