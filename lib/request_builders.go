//go:generate go run ../gen/request_builders.go -t request_builders.go.inc -o request_builders_generated.go
package lib

// BuildRequestNoAuth creates a new Request with the intention of not authenticating.
func (c *bytemarkClient) BuildRequestNoAuth(method string, endpoint Endpoint, path string, parts ...string) (r Request, err error) {
	url, err := c.BuildURL(endpoint, path, parts...)
	if err != nil {
		return
	}
	return &internalRequest{
		client:        c,
		endpoint:      endpoint,
		url:           url,
		method:        method,
		allowInsecure: c.allowInsecure,
	}, nil
}

// BuildRequest builds a request that will be authenticated by the endpoint given.
func (c *bytemarkClient) BuildRequest(method string, endpoint Endpoint, path string, parts ...string) (r Request, err error) {
	url, err := c.BuildURL(endpoint, path, parts...)
	if err != nil {
		return
	}
	return &internalRequest{
		authenticate:  true,
		client:        c,
		endpoint:      endpoint,
		url:           url,
		method:        method,
		allowInsecure: c.allowInsecure,
	}, nil
}
