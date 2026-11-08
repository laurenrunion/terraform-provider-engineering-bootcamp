package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// default address of the DevOps API when run locally
const DefaultHost = "http://localhost:8080"

// Client holds what we need to talk to the API: where it is, and goes built in HTTP client
type Client struct {
	HostURL    string
	HTTPClient *http.Client
}

// NewClient builds a Client. If no host is given, use the default.
// (No username/password like HashiCups, because our API has no auth)
func NewClient(host *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    DefaultHost,
	}
	if host != nil && *host != "" {
		c.HostURL = *host
	}
	return &c, nil
}

// Engineer is a plain Go type (NOT a Terraform type). It is used to build and read the JSON.
// If your curl showed "id" as a number, change ID to int.
type Engineer struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

/*all i have to do is call doRequest("POST", "/engineers"(endpt), engineer(data tf)) and it will*/
/*build the struct into json*/
/*builds a request for http://localhost:8080/engineers */
/*sends and the api replies if status is bad it errors*/
//does error checking the whole time
// doRequest sends one request to the API and returns the response body.
// Every CRUD function in Exercise 4 will use this.
func (c *Client) doRequest(method, path string, payload interface{}) ([]byte, error) {
    //turns paload into json
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
    //Build the request method + host + path + JSON
	req, err := http.NewRequest(method, c.HostURL+path, body)
    //if last step failed stop and report the problem. called a bunch as you can see
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

//send the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
//read the reply and check status
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// treat any non-2xx status as an error
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d: %s", res.StatusCode, string(data))
	}
	return data, nil
}
