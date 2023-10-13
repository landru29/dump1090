// Package text is the text serializer.
package text

import (
	"strings"

	"github.com/landru29/dump1090/internal/dump"
)

// Serializer is the text serializer.
type Serializer struct {
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	if ac == nil {
		return nil, nil
	}

	switch aircraft := ac.(type) {
	case dump.Aircraft:
		return s.Serialize([]*dump.Aircraft{&aircraft})
	case *dump.Aircraft:
		return s.Serialize([]*dump.Aircraft{aircraft})
	case []dump.Aircraft:
		out := make([]*dump.Aircraft, len(aircraft))
		for idx := range aircraft {
			out[idx] = &aircraft[idx]
		}
		return s.Serialize(out)
	case []*dump.Aircraft:
		out := make([]string, len(aircraft))
		for idx, plane := range aircraft {
			out[idx] = plane.String()
		}

		return []byte(strings.Join(out, "\n")), nil
	}

	return nil, nil
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "text/plain"
}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "text"
}
