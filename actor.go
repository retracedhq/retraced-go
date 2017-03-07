package retraced

// Actor represents an entity that performs an action
type Actor struct {
	// Id is the unique id (in the environment) for this actor
	ID string `json:"id"`

	// Name is the display name for this actor. It can be email
	Name string `json:"name"`

	// Type represents the type of actor. This is often "user" or "token"
	Type string `json:"type"`

	// URL represents a URL to the actor
	URL string `json:"url"`
}
