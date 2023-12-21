// Package serialize describes how to serialize aircraft.
package serialize

// Serializer is the aircraft serializer.
type Serializer interface {
	Serialize(ac ...any) ([]byte, error)
	MimeType() string
	String() string
}
