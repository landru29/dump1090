// Package text is the text serializer.
package text

import (
	"strings"

	"github.com/landru29/dump1090/internal/model"
)

// Serializer is the text serializer.
type Serializer struct{}

// Serialize implements the Serialize.Serializer interface.
func (s Serializer) Serialize(ac any) ([]byte, error) {
	if ac == nil {
		return nil, nil
	}

	switch aircraft := ac.(type) {
	case model.Aircraft:
		return s.Serialize([]*model.Aircraft{&aircraft})
	case *model.Aircraft:
		return s.Serialize([]*model.Aircraft{aircraft})
	case []model.Aircraft:
		out := make([]*model.Aircraft, len(aircraft))
		for idx := range aircraft {
			out[idx] = &aircraft[idx]
		}

		return s.Serialize(out)
	case []*model.Aircraft:
		out := make([]string, len(aircraft))
		for idx, plane := range aircraft {
			out[idx] = plane.String()
		}

		return []byte(strings.Join(out, "\n")), nil
	}

	return nil, nil
}

// MimeType implements the Serialize.Serializer interface.
func (s Serializer) MimeType() string {
	return "text/plain"
}

// String implements the Serialize.Serializer interface.
func (s Serializer) String() string {
	return "text"
}
