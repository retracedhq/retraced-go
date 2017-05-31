package retraced

// Target represents an object that an action was taken on
type Target struct {
	// Id is the id of the target
	ID string `json:"id"`

	// Name can be used to represent the name of the target
	Name string `json:"name"`

	// Type describes the type of target
	Type string `json:"type"`

	// Href is a reference to the target
	Href string `json:"href"`

	// Fields can store any additional data on the target
	Fields Fields `json:"fields,omitempty"`
}
