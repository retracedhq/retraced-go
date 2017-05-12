package retraced

import "log"

// Initialize a new client with your projectID and API key and then configure options.
func ExampleClient() {
	client, err := NewClient("f4228ca2220d4d0a89d39a93f9987658", "ce6eba2ba9534e94ad48624079bcccf6")
	if err != nil {
		log.Fatal(err)
	}
	client.Component = "Web Dashboard"
	client.Version = "0.3.0"
	client.ViewLogAction = "audit.log.view"
}
