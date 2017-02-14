package retraced

// Group represents a tenant in the application. Audit logs will be scoped by this value.
type Group struct {
	// Id is the unique id (in the environment) for this group/team
	Id string `json:"id"`

	// Name is the display name for this group/team.
	Name string `json:"name"`
}
