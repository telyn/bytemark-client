// Client library for auth3. Thin wrapper around net/http that supports creating
// a session and reading it back

package client


import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	sessionEndpoint *url.URL
	HTTP            http.Client
}

// n-factor auth. We expect factor => credential, i.e:
// {"password" => "foo", "yubikey" => "cccbar"}
type Credentials map[string]string

type CreateSessionError struct {
	httpErr error
}

func (e CreateSessionError) Error() string {
	return fmt.Sprintf("Failed to create session: %v", e.httpErr)
}

// Data in the session. We expect it to look like this.
type SessionData struct {
	Token      string // not actually in the session, but communicate it here
	Username   string
	Factors  []string

	// The groups this user is a member of
	GroupMemberships []string `json:"group_memberships"`
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	rsp, rspErr := c.HTTP.Do(req)
	if rspErr != nil {
		return nil, rspErr
	}
	defer rsp.Body.Close()

	body := make([]byte, rsp.ContentLength)
	for count := 0; count < len(body); {
		n, err := rsp.Body.Read(body[count:])
		count = count + n
		// We shouldn't get EOF, we're not trying to read past the end
		if err != nil && err != io.EOF {
			return nil, err // no interest in partial bodies
		}
	}

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		 if len(body) == 0 {
			return nil, errors.New(rsp.Status)
		}
		return nil, errors.New(string(body))
	}

	return body, nil
}

func New(endpoint string) (*Client, error) {
	// ensure we end up with a string like "https://example.com/session"
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, fmt.Errorf("Endpoint scheme must be http or https, got: %s", parsed.Scheme)
	}

	if len(parsed.Path) > 0 {
		if parsed.Path[len(parsed.Path):] != "/" {
			parsed.Path = parsed.Path + "/"
		}
	} else {
		parsed.Path = "/"
	}
	parsed.Path = parsed.Path + "session"

	return &Client{sessionEndpoint: parsed}, nil
}

func (c *Client) ReadSession(token string) (*SessionData, error) {
	x := *c.sessionEndpoint // shallow copy. Don't touch UserInfo
	x.Path = x.Path + "/" + token
	req, reqErr := http.NewRequest("GET", x.String(), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Add("Accept", "application/json")
	body, bodyErr := c.doRequest(req)
	if bodyErr != nil {
		return nil, bodyErr
	}
	if len(body) == 0 {
		return nil, errors.New("Empty body returned reading session")
	}

	out := &SessionData{Token: token} // not included in session data
	jsonErr := json.Unmarshal(body, out)
	if jsonErr != nil {
		return nil, jsonErr
	}


	return out, nil
}

// Creates a session and returns the token
func (c *Client) CreateSessionToken(credentials Credentials) (string, error) {
	data, marshalErr := json.Marshal(credentials)
	if marshalErr != nil {
		return "", marshalErr
	}

	req, reqErr := http.NewRequest("POST", c.sessionEndpoint.String(), bytes.NewBuffer(data))
	if reqErr != nil {
		return "", reqErr
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "text/plain")

	body, err := c.doRequest(req)
	if err != nil || len(body) == 0 {
		return "", CreateSessionError{err}
	}

	// FIXME: auth should really put the token in an Authorization: header.
	// TODO:  It could then return the session data in the response body.
	return string(body), nil
}

// Creates a session, returning the session data rather than just the token
func (c *Client) CreateSession(credentials Credentials) (*SessionData, error) {

	token, createErr := c.CreateSessionToken(credentials)
	if createErr != nil {
		return nil, createErr
	}

	sessionData, getErr := c.ReadSession(token)
	if getErr != nil {
		return nil, fmt.Errorf("Created session but couldn't read it: %s", getErr)
	}

	return sessionData, nil
}

// TODO: func (c *Client) CreateUser() {}
// TODO: func (c *Client) ReadUser() {}
// TODO: func (c *Client) IsUsernameAvailable(username string) {}
// TODO: func (c *Client) ResetUserPassword() {}

