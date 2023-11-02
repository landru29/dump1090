// Package none is the empty serializer.
package none

// Serializer is the none serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(_ any) ([]byte, error) {
	return nil, nil
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return ""
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "none"
}
