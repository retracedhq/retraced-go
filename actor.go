package retraced

import (
	"time"
)

type Actor struct {
	ID         string    `json:"id"`
	ForeignID  string    `json:"foreign_id"`
	Name       string    `json:"name"`
	Created    time.Time `json:"created"`
	LastActive time.Time `json:"last_active"`
	EventCount int64     `json:"event_count"`
}
