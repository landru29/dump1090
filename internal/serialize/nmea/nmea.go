// Package nmea is the nmea serializer.
package nmea

// Serializer is the nmea serializer.
type Serializer struct {
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	return nil, nil
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "text/plain"
}
