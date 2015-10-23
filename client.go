package auditable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	ProjectId string
	Token     string
	Endpoint  string
}

func NewClient(projectId string, apiToken string) (*Client, error) {
	return &Client{
		ProjectId: projectId,
		Token:     apiToken,
		Endpoint:  "https://api.auditable.io",
	}, nil
}

func (c *Client) ReportEvent(event *Event) error {
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/project/%s/event", c.Endpoint, c.ProjectId), bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected response from auditable api: %d", resp.StatusCode)
	}

	return nil
}
