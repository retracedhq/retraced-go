package retraced

// Actor represents an entity that performs an action
type Actor struct {
	// Id is the unique id (in the environment) for this actor
	ID string `json:"id"`

	// Name is the display name for this actor. It can be email
	Name string `json:"name"`

	// Href represents a URL to the actor
	Href string `json:"href"`

	// Fields are any additional data to store with the actor
	Fields map[string]string `json:"fields,omitempty"`
}
