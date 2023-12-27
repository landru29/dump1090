// Package transport manages how the data is transported.
package transport

import "github.com/landru29/dump1090/internal/model"

// Transporter is the data transporter.
type Transporter interface {
	Transport(ac *model.Aircraft) error
	String() string
}
