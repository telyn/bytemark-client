package lib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// RequestAlreadyRunError is returned if the Run method was already called for this Request.
type RequestAlreadyRunError struct {
	Request *Request
}

// InsecureConnectionError is returned if the endpoint isn't https but AllowInsecure was not called.
type InsecureConnectionError struct {
	Request *Request
}

func (e RequestAlreadyRunError) Error() string {
	return "A Request was Run twice"
}

func (e InsecureConnectionError) Error() string {
	return "A Request to an insecure endpoint was attempted when AllowInsecureRequests had not been called."
}

// Request is the workhorse of the bytemark-client/lib - it builds up a request, then Run can be called to get its results.
type Request struct {
	authenticate  bool
	client        Client
	endpoint      Endpoint
	url           *url.URL
	method        string
	body          []byte
	allowInsecure bool
	hasRun        bool
}

// GetURL returns the URL that the Request is for.
func (r *Request) GetURL() url.URL {
	if r.url == nil {
		return url.URL{}
	}
	return *r.url
}

// BuildRequestNoAuth creates a new Request with the intention of not authenticating.
func (c *bytemarkClient) BuildRequestNoAuth(method string, endpoint Endpoint, path string, parts ...string) (r *Request, err error) {
	url, err := c.BuildURL(endpoint, path, parts...)
	if err != nil {
		return
	}
	return &Request{
		client:        c,
		endpoint:      endpoint,
		url:           url,
		method:        method,
		allowInsecure: c.allowInsecure,
	}, nil
}

// BuildRequest builds a request that will be authenticated by the endpoint given.
func (c *bytemarkClient) BuildRequest(method string, endpoint Endpoint, path string, parts ...string) (r *Request, err error) {
	url, err := c.BuildURL(endpoint, path, parts...)
	if err != nil {
		return
	}
	return &Request{
		authenticate:  true,
		client:        c,
		endpoint:      endpoint,
		url:           url,
		method:        method,
		allowInsecure: c.allowInsecure,
	}, nil
}

// AllowInsecure tells the Request that it's ok if the endpoint isn't communicated with over HTTPS.
func (r *Request) AllowInsecure() {
	r.allowInsecure = true
}

// mkHTTPClient creates an http.Client for this request. If the staging endpoint is used, InsecureSkipVerify is used because I guess we don't have a good cert for that brain.
func (r *Request) mkHTTPClient() (c *http.Client) {
	c = new(http.Client)
	if r.url.Host == "staging.bigv.io" {
		c.Transport = &http.Transport{
			// disable gas lint for this line (gas looks for insecure TLS settings, among other things)
			/* #nosec */
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return c
}

// mkHTTPRequest assembles an http.Request for this Request, adding Authorization headers as needed, setting the Content-Type correctly for whichever endpoint it's talking to.
func (r *Request) mkHTTPRequest(body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(r.method, r.url.String(), body)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("User-Agent", "bytemark-client-"+Version)

	if r.endpoint == SPPEndpoint {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}
	if r.authenticate {
		if r.client.GetSessionToken() == "" {
			return nil, NilAuthError{}
		}
		// if we could settle on a single standard
		// rather than two basically-identical ones that'd be cool
		if r.endpoint == BillingEndpoint {
			req.Header.Add("Authorization", "Token token="+r.client.GetSessionToken())
		} else {
			req.Header.Add("Authorization", "Bearer "+r.client.GetSessionToken())
		}
	}
	return
}

// MarshalAndRun marshals the 'in' object and passes that and 'out' to Run as the body and responseObject, respectively.
func (r *Request) MarshalAndRun(in interface{}, out interface{}) (statusCode int, responseBody []byte, err error) {
	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(in)
	if err != nil {
		return
	}

	return r.Run(&b, out)
}

// Run performs the request with the given body, and attempts to unmarshal a successful response into responseObject
func (r *Request) Run(body io.Reader, responseObject interface{}) (statusCode int, responseBody []byte, err error) {
	if r.hasRun {
		err = RequestAlreadyRunError{r}
		return
	}
	r.hasRun = true

	if !r.allowInsecure && r.url.Scheme == "http" {
		err = InsecureConnectionError{r}
		return
	}
	var rb []byte
	if body != nil {

		rb, err = ioutil.ReadAll(body)
		if err != nil {
			return 0, nil, err
		}
		log.Debugf(log.LvlHTTPData, "request body: '%s'\r\n", string(rb))
	}

	cli := r.mkHTTPClient()

	req, err := r.mkHTTPRequest(bytes.NewBuffer(rb))
	if err != nil {
		return
	}
	if len(rb) > 0 {
		req.Header.Add("Content-Length", fmt.Sprintf("%d", len(rb)))
	}
	res, err := cli.Do(req)
	if err != nil {
		return
	}

	statusCode = res.StatusCode

	log.Debugf(log.LvlOutline, "%s %s: %d\r\n", r.method, req.URL, res.StatusCode)

	responseBody, err = r.handleResponse(req, rb, res, responseObject)
	return
}

// handleResponse deals with the response coming back from the request - creating an error if required, unmarshalling responseObject if necessary
func (r *Request) handleResponse(req *http.Request, requestBody []byte, res *http.Response, responseObject interface{}) (body []byte, err error) {
	body, err = ioutil.ReadAll(res.Body)
	log.Debugf(log.LvlHTTPData, "response body: '%s'\r\n", body)
	if err != nil {
		return
	}

	baseErr := APIError{
		Method:      r.method,
		URL:         req.URL,
		StatusCode:  res.StatusCode,
		RequestBody: string(requestBody),
	}

	baseErr.ResponseBody = string(body)

	switch res.StatusCode {
	case http.StatusBadRequest:
		// because we need to reference fields specific to BadRequestError later
		err = newBadRequestError(baseErr, body)
	case http.StatusUnauthorized:
		err = UnauthorizedError{baseErr}
	case http.StatusForbidden:
		err = ForbiddenError{baseErr}
	case http.StatusNotFound:
		err = NotFoundError{baseErr}
	case http.StatusInternalServerError:
		err = InternalServerError{baseErr}
	case http.StatusServiceUnavailable:
		err = ServiceUnavailableError{baseErr}
	default:
		if 200 <= res.StatusCode && res.StatusCode <= 299 {
			if responseObject != nil {
				jsonErr := json.Unmarshal(body, responseObject)
				if jsonErr != nil {
					return body, jsonErr
				}
			}
		} else {
			err = UnknownStatusCodeError{baseErr}
		}
	}
	return
}

// BuildURL pieces together a URL from parts, escaping as necessary..
func (c *bytemarkClient) BuildURL(endpoint Endpoint, format string, args ...string) (*url.URL, error) {
	arr := make([]interface{}, len(args))
	for i, str := range args {
		arr[i] = url.QueryEscape(str)
	}
	endpointURL := ""
	switch endpoint {
	case APIEndpoint:
		endpointURL = c.urls.API
	case BrainEndpoint:
		endpointURL = c.urls.Brain
	case BillingEndpoint:
		endpointURL = c.urls.Billing
	case SPPEndpoint:
		endpointURL = c.urls.SPP
	default:
		return nil, UnsupportedEndpointError(endpoint)
	}
	if !strings.HasPrefix(format, "/") {
		return nil, UnsupportedEndpointError(-1)
	}
	return url.Parse(endpointURL + fmt.Sprintf(format, arr...))
}
