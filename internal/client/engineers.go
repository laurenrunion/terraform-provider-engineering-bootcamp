package client

import (
	"encoding/json"
	"errors"
	"strings"
)

// ErrNotFound means the API has no engineer with that id.
// (delete this line if you already declared ErrNotFound in client.go)
var ErrNotFound = errors.New("engineer not found")

// isNotFound checks the error text, because this API returns 400 instead of 404.
func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "No engineer with that ID")
}


// CreateEngineer creates a new engineer in the API.
// It passes "POST", the path /engineers, and the engineer's name and email to doRequest.
// It is called by Create in engineer_resource.go (the provider package).
// If doRequest returns an error, it passes the error back up to the provider.
// If the request works, it turns the JSON response into an Engineer, which now
// has the id the API generated, and returns it.
func (c *Client) CreateEngineer(e Engineer) (*Engineer, error) {
	body, err := c.doRequest("POST", "/engineers", e)
	if err != nil {
		return nil, err
	}
	var out Engineer
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Read GetEngineer sends GET /engineers/id/<id>.
func (c *Client) GetEngineer(id string) (*Engineer, error) {
	body, err := c.doRequest("GET", "/engineers/id/"+id, nil)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var out Engineer
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateEngineer sends PUT /engineers/<id>, then re-reads the engineer.
// I re-read because I don't know exactly what PUT returns.
func (c *Client) UpdateEngineer(id string, e Engineer) (*Engineer, error) {
	_, err := c.doRequest("PUT", "/engineers/"+id, e)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c.GetEngineer(id)
}

// DeleteEngineer sends DELETE /engineers/<id>.
func (c *Client) DeleteEngineer(id string) error {
	_, err := c.doRequest("DELETE", "/engineers/"+id, nil)
	if err != nil {
		if isNotFound(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
