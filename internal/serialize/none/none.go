// Package none is the empty serializer.
package none

// Serializer is the none serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	return nil, nil
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return ""
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "none"
}
