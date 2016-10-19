// Package superfeedr provides basic methods for the superfeedr api.
package superfeedr

// Superfeedr represents the object used to work with.
type Superfeedr struct {
	username string
	password string
}

// NewSuperfeedr creates and sets the basic attributes.
func NewSuperfeedr(username string, password string) *Superfeedr {
	return &Superfeedr{username: username, password: password}
}
