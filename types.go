// Package xchedule provides ...
package xchedule

import (
	"time"
)

// Event is the unit of schedule
type Event struct {
	Title     string
	Time      Time
	Locations []Location
	Members   []Member
	Schedule  []Event
	Alert     Alert
	Notes     []string
	root      bool
}

// SetRoot to indicate root event or not
func (e *Event) SetRoot(isRoot bool) {
	e.root = isRoot
}

// IsRoot return true if this event is root event
func (e *Event) IsRoot() bool {
	return e.root
}

// Time indicates the time or period
type Time struct {
	Start  time.Time
	End    time.Time
	Period bool
}

// Location records event
type Location struct {
	Name    string
	Address string
	LatLng  string
	MapURL  string
}

// Member who attend this event
type Member struct {
	Name string
}

// Alert Info
type Alert struct {
	Times []time.Time
}
