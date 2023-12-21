// Package json is the json serializer.
package json

import (
	"encoding/json"

	"github.com/landru29/dump1090/internal/model"
)

// Serializer is the json serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(planes ...any) ([]byte, error) {
	output := []model.Aircraft{}

	for _, ac := range planes {
		switch aircraft := ac.(type) {
		case model.Aircraft:
			output = append(output, aircraft)

		case *model.Aircraft:
			output = append(output, *aircraft)

		case []model.Aircraft:
			output = append(output, aircraft...)

		case []*model.Aircraft:
			for _, elt := range aircraft {
				output = append(output, *elt)
			}
		}
	}

	return json.Marshal(output)
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/json"
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "json"
}
