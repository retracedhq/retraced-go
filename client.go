package retraced

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	apiVersion      = 2
	defaultEndpoint = "http://localhost:3000/auditlog"
)

// Client represents a client that can send events into the retraced service.
type Client struct {
	projectID string
	token     string
	// Endpoint is the retraced api base url, default is `https://api.retraced.io`
	Endpoint string
	// Component is an identifier for a specific component of a vendor app platform
	Component string
	// Version is an identifier for the specific version of this component, usually a git SHA
	Version string
	// ViewLogAction is the action logged when a Viewer Token is used, default is 'audit.log.view'
	ViewLogAction string
	//
	HttpClient *http.Client
}

// NewClient creates a new retraced api client that can be used to send events
func NewClient(endpoint string, projectID string, apiToken string) (*Client, error) {
	ep := defaultEndpoint
	if endpoint != "" {
		ep = endpoint
	}

	return &Client{
		projectID:  projectID,
		token:      apiToken,
		Endpoint:   ep,
		HttpClient: http.DefaultClient,
	}, nil
}

// NewClientWithVersion Same as NewClient, but includes params for specifying the
// Component and Version of the Retraced client application
func NewClientWithVersion(endpoint string, projectID string, apiToken string, component string, version string) (*Client, error) {
	ep := defaultEndpoint
	if endpoint != "" {
		ep = endpoint
	}

	return &Client{
		projectID:  projectID,
		token:      apiToken,
		Endpoint:   ep,
		Component:  component,
		Version:    version,
		HttpClient: http.DefaultClient,
	}, nil
}

// NewEventRecord is returned from the Retraced API when an event is created
type NewEventRecord struct {
	ID   string `json:"id"`
	Hash string `json:"hash"`
}

// ReportEvent is the method to call to send a new event.
func (c *Client) ReportEvent(event *Event) (*NewEventRecord, error) {
	event.apiVersion = apiVersion
	if event.Version == "" {
		event.Version = c.Version
	}
	if event.Component == "" {
		event.Component = c.Component
	}

	encoded, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/publisher/v1/project/%s/event", c.Endpoint, c.projectID), bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.token))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected response from retraced api endpoint %s: %d", req.URL.String(), resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reqResp NewEventRecord
	if err := json.Unmarshal(bodyBytes, &reqResp); err != nil {
		return nil, err
	}

	if err := event.VerifyHash(&reqResp); err != nil {
		return nil, err
	}

	return &reqResp, nil
}

// GetViewerToken will return a one-time use token that can be used to view a group's audit log.
func (c *Client) GetViewerToken(groupID string, isAdmin bool, actorID string, targetID string) (*ViewerToken, error) {
	params := url.Values{}
	params.Add("group_id", groupID)
	params.Add("is_admin", strconv.FormatBool(isAdmin))
	params.Add("actor_id", actorID)
	params.Add("view_log_action", c.ViewLogAction)

	if targetID != "" {
		params.Add("target_id", targetID)
	}

	u, err := url.Parse(fmt.Sprintf("%s/publisher/v1/project/%s/viewertoken", c.Endpoint, c.projectID))
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.token))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK { // There's a pending PR in the retraced API to match this.
		return nil, fmt.Errorf("unexpected response from retraced api endpoint %s: %d", req.URL.String(), resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	viewerToken := ViewerToken{}
	if err := json.Unmarshal(contents, &viewerToken); err != nil {
		return nil, err
	}

	return &viewerToken, nil
}

// DeleteViewerSessions will delete all viewer sessions for the given actor in the given group.
func (c *Client) DeleteViewerSessions(groupID string, actorID string) error {
	url := fmt.Sprintf("%s/v1/project/%s/group/%s/actor/%s/viewersessions", c.Endpoint, c.projectID, groupID, actorID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response from retraced api: %d", resp.StatusCode)
	}

	return nil
}

// Query searches for events using the Publisher API's GraphQL endpoint.
func (c *Client) Query(sq *StructuredQuery, mask *EventNodeMask, pageSize int) (EventsPager, error) {
	url := fmt.Sprintf("%s/publisher/v1/project/%s/graphql", c.Endpoint, c.projectID)
	ec := &EventsConnection{
		url:             url,
		authorization:   fmt.Sprintf("Token token=%s", c.token),
		structuredQuery: sq,
		mask:            mask,
		pageSize:        pageSize,
		httpClient:      c.HttpClient,
	}

	err := ec.call()
	if err != nil {
		return nil, err
	}

	return ec, nil
}
