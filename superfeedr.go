// Package superfeedr provides basic methods for the superfeedr api.
package superfeedr

import (
	"time"
)

// Feed contains the title and multiple Items
type Feed struct {
	Title string  `json:"title"`
	Items []*Item `json:"items"`
}

// Item contains all infos of an item.
type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	Published time.Time `json:"published"`
	Updated   time.Time `json:"updated"`
	Links     []*Link
	// standardLinks interface{} `json:"standardLinks"`
}

// Link represents a structured link.
type Link struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"alternate"`
}

// Superfeedr represents the object used to work with.
type Superfeedr struct {
	username string
	password string
}

// NewSuperfeedr creates and sets the basic attributes.
func NewSuperfeedr(username string, password string) *Superfeedr {
	return &Superfeedr{username: username, password: password}
}

// Retrieve entries for the given topic. You must be a subscriber of the
// given topic.
func (f *Superfeedr) Retrieve(topic string) (Feed, error) {
	return Feed{}, nil
}
