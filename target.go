package retraced

// Target represents an object that an action was taken on
type Target struct {
	// ID is the id of the target
	ID string `json:"id"`

	// Name can be used to represent the name of the target
	Name string `json:"name"`

	// Type describes the type of target
	Type string `json:"type"`

	// URL is a reference to the target
	URL string `json:"url"`

	// Fields can store any additional data on the target
	Fields map[string]interface{} `json:"fields,omitempty"`
}
