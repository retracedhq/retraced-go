package auditable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (c *Client) GetViewerToken(foreignActorId string, foreignTeamId string) (*ViewerToken, error) {
	params := url.Values{}
	params.Add("actor_id", foreignActorId)
	params.Add("team_id", foreignTeamId)

	fmt.Printf(fmt.Sprintf("%s/v1/project/%s/viewertoken", c.Endpoint, c.ProjectId))
	u, err := url.Parse(fmt.Sprintf("%s/v1/project/%s/viewertoken", c.Endpoint, c.ProjectId))
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response from auditable api: %d", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	viewerToken := ViewerToken{}
	if err := json.Unmarshal(contents, &viewerToken); err != nil {
		return nil, err
	}

	return &viewerToken, nil
}
