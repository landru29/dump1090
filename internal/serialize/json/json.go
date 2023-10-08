// Package json is the json serializer.
package json

import (
	"encoding/json"
)

// Serializer is the json serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	return json.Marshal(ac)
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "application/json"
}
