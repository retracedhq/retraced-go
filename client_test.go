package retraced

import "log"

// Initialize a new client with your projectID and API key and then configure options.
func ExampleClient() {
	client, err := NewClient("", "dev", "dev")
	if err != nil {
		log.Fatal(err)
	}
	client.Component = "Web Dashboard"
	client.Version = "0.3.0"
	client.ViewLogAction = "audit.log.view"
}
