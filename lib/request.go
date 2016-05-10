package lib

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type RequestAlreadyRunError struct {
	Request *Request
}
type InsecureConnectionError struct {
	Request *Request
}

func (e RequestAlreadyRunError) Error() string {
	return "A Request was Run twice"
}

func (e InsecureConnectionError) Error() string {
	return "A Request to an insecure endpoint was attempted when AllowInsecure had not been called."
}

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

func (r *Request) GetURL() url.URL {
	if r.url == nil {
		return url.URL{}
	}
	return *r.url
}

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

func (r *Request) AllowInsecure() {
	r.allowInsecure = true
}

func (r *Request) mkHTTPClient() (c *http.Client) {
	c = new(http.Client)
	if r.url.Host == "staging.bigv.io" {
		c.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return c
}

func (r *Request) mkHTTPRequest(body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(r.method, r.url.String(), body)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("User-Agent", "bytemark-client"+GetVersion().String())

	if r.endpoint == EP_SPP {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}
	if r.authenticate {
		if r.client.GetSessionToken() == "" {
			return nil, &NilAuthError{}
		}
		// if we could settle on a single standard
		// rather than two basically-identical ones that'd be cool
		if r.endpoint == EP_BILLING {
			req.Header.Add("Authorization", "Token token="+r.client.GetSessionToken())
		} else {
			req.Header.Add("Authorization", "Bearer "+r.client.GetSessionToken())
		}
	}
	return
}

func (r *Request) Run(body io.Reader, responseObject interface{}) (statusCode int, response []byte, err error) {
	if r.hasRun {
		err = RequestAlreadyRunError{r}
		return
	}
	r.hasRun = true

	if !r.allowInsecure && r.url.Scheme == "http" {
		err = InsecureConnectionError{r}
		return
	}
	rb := make([]byte, 0)
	if body != nil {

		rb, err = ioutil.ReadAll(body)
		if err != nil {
			return 0, nil, err
		}
		log.Debugf(log.DBG_HTTPDATA, "request body: '%s'\r\n", string(rb))
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

	log.Debugf(log.DBG_OUTLINE, "%s %s: %d\r\n", r.method, req.URL, res.StatusCode)

	baseErr := APIError{
		Method:      r.method,
		URL:         req.URL,
		StatusCode:  res.StatusCode,
		RequestBody: string(rb),
	}

	response, err = ioutil.ReadAll(res.Body)
	log.Debugf(log.DBG_HTTPDATA, "response body: '%s'\r\n", response)
	if err != nil {
		return
	}
	baseErr.ResponseBody = string(response)

	switch res.StatusCode {
	case 400:
		err := BadRequestError{APIError: baseErr, Problems: make(map[string][]string)}
		jsonErr := json.Unmarshal(response, &err.Problems)
		if jsonErr != nil {
			log.Debug(log.DBG_OUTLINE, "Couldn't parse 400 response into JSON, so bunging it into a single Problem in the BadRequestError")
			err.Problems["The problem"] = []string{baseErr.ResponseBody}
		}
	case 403:
		err = NotAuthorizedError{baseErr}
	case 404:
		err = NotFoundError{baseErr}
	case 500:
		err = InternalServerError{baseErr}
	default:
		if 200 <= res.StatusCode && res.StatusCode <= 299 {
			if responseObject != nil {
				jsonErr := json.Unmarshal(response, responseObject)
				if jsonErr != nil {
					return statusCode, response, jsonErr
				}
			}
			break
		}
		err = UnknownStatusCodeError{baseErr}
	}
	return
}

// BuildURL pieces together a URL from parts, escaping as necessary..
func (c *bytemarkClient) BuildURL(endpoint Endpoint, format string, args ...string) (*url.URL, error) {
	arr := make([]interface{}, len(args), len(args))
	for i, str := range args {
		arr[i] = url.QueryEscape(str)
	}
	endpointUrl := ""
	switch endpoint {
	case EP_BRAIN:
		endpointUrl = c.brainEndpoint
	case EP_BILLING:
		endpointUrl = c.billingEndpoint
	case EP_SPP:
		endpointUrl = c.sppEndpoint
	default:
		return nil, UnsupportedEndpointError(endpoint)
	}
	if !strings.HasPrefix(format, "/") {
		return nil, UnsupportedEndpointError(-1)
	}
	return url.Parse(endpointUrl + fmt.Sprintf(format, arr...))
}
