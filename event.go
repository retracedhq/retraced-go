package retraced

import (
	"time"
)

// Event represents a single audited event.
type Event struct {
	// Created is a timestamp representing when the event took place
	Created time.Time `json:"created"`

	// CRUD is a list of the most basic verbs that describe the type of action
	CRUD string `json:"crud"`

	// Action is a short, readable word to describe the action
	Action string `json:"action"`

	// Object represents the item that had an action performed on it
	Object *Object `json:"object,omitempty"`

	// Title is the title of the event to display in the audit log
	Title string `json:"title"`

	// Description is a string with the full details of the event
	Description string `json:"description"`

	// SourceIP is the client ip address where the event was performed
	SourceIP string `json:"source_ip"`

	// Actor represents the entity that performed the action
	Actor *Actor `json:"actor"`

	// TeamID is the team that will be able to see this event in the audit log
	Group *Group `json:"group"`

	// Fields are any additional data to store with the event
	Fields map[string]interface{} `json:"fields,omitempty"`

	// IsFailure is an optional flag that, when set, indicates that this audited event is a failed use of privileges
	IsFailure bool `json:"is_failure"`

	// IsAnonymous is an optional flag that, when set, indicates that this is an anonymous event
	IsAnonymous bool `json:"is_anonymous"`

	// apiVersion is set here to allow updates to this model without breaking the API server
	apiVersion int
}
