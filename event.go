package auditable

import (
	"time"
)

type Event struct {
	Id          string         `json:"id"`
	Created     time.Time      `json:"created"`
	Title       string         `json:"title"`
	Action      string         `json:"action"`
	Description string         `json:"description"`
	SourceIp    string         `json:"source_ip"`
	Actor       *Actor         `json:"actor"`
	Location    *EventLocation `json:"location,omitempty"`
	Teams       []string       `json:"teams,omitempty"`
}
