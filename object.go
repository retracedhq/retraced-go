package retraced

// Object represents an object that an action was taken on
type Object struct {
	// ID is the id of the object
	ID string `json:"id"`

	// Name can be used to represent the name of the object
	Name string `json:"name"`

	// Type describes the type of object
	Type string `json:"type"`

	// URL is a reference to the object
	URL string `json:"url"`

	// Fields can store any additional data on the object
	Fields map[string]interface{} `json:"fields,omitempty"`
}
