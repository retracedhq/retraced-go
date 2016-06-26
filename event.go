package retraced

import (
	"time"
)

type Event struct {
	ID          string         `json:"id"`
	Created     time.Time      `json:"created"`
	Title       string         `json:"title"`
	Action      string         `json:"action"`
	Description string         `json:"description"`
	SourceIP    string         `json:"source_ip"`
	Actor       *Actor         `json:"actor"`
	Location    *EventLocation `json:"location,omitempty"`
	Teams       []string       `json:"teams,omitempty"`
}
