// Package transport manages how the data is transported.
package transport

import "github.com/landru29/dump1090/internal/source"

// Transporter is the data transporter.
type Transporter interface {
	Transport(ac *source.Aircraft) error
}
