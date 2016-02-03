package lib

import (
	"bytemark.co.uk/client/util/log"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	authenticate    bool
	client          *Client
	url             *url.URL
	method          string
	body            []byte
	insecureAllowed bool
	hasRun          bool
}

func (c *bytemarkClient) NewRequestNoAuth(method, path string) (*Request, error) {
	u := url.Parse(path)
	return &Request{
		client: c,
		url:    u,
		method: method,
	}
}

func (c *bytemarkClient) NewRequest(method, path string) (*Request, error) {
	u := url.Parse(path)
	return &Request{
		authenticate: true,
		client:       c,
		url:          u,
		method:       method,
	}

}

func (r *Request) AllowInsecure() {
	r.insecureAllowed = true
}

/*
func (r *Request) mkHTTPClient() (c *http.Client) {
	if r.url.Host == "staging.bigv.io" {
		c.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return c
}
*/

func (r *Request) mkHTTPRequest() (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, bytes.NewBufferString(requestBody))
	if err != nil {
		return req, nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if auth {
		if c.authSession == nil {
			return nil, nil, &NilAuthError{}
		}
		req.Header.Add("Authorization", "Bearer "+c.authSession.Token)
	}
}

func (r *Request) Run(body []byte, responseObject interface{}) (statusCode int, response []byte, err error) {
	if r.hasRun {
		return []byte{}, RequestAlreadyRunError{r}
	}
	r.hasRun = true

	if !r.insecureAllowed && r.url.Scheme == "http" {
		return []byte{}, InsecureConnectionError{r}
	}

	cli := mkHTTPClient()

	req, err := r.mkHTTPRequest(body)
	if err != nil {
		return
	}

	res, err = cli.Do(req)
	if err != nil {
		return
	}

	statusCode = res.StatusCode

	log.Debugf(log.DBG_OUTLINE, "%s %s: %d\r\n", r.method, req.URL, res.StatusCode)
	log.Debugf(log.DBG_HTTPDATA, "request body: '%s'\r\n", string(body))

	baseErr := APIError{
		Method:       method,
		URL:          req.URL,
		StatusCode:   res.StatusCode,
		RequestBody:  string(body),
		ResponseBody: "",
	}

	response, err = ioutil.ReadAll(res.Body)
	log.Debugf(log.DBG_HTTPDATA, "response body: '%s'\r\n", responseBody)
	if err != nil {
		return
	}
	baseErr.ResponseBody = string(responseBody)

	switch res.StatusCode {
	case 400:
		err := BadRequestError{APIError: baseErr, Problems: make(map[string][]string)}
		jsonErr = json.Unmarshal(responseBody, &err.Problems)
		if jsonErr != nil {
			log.Debug(log.DBG_OUTLINE, jsonErr)
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
			jsonErr := json.Unmarshal(response, responseObject)
			if jsonErr != nil {
				return statusCode, response, jsonErr
			}
			break
		}
		err = UnknownStatusCodeError{baseErr}
	}
	return
}

// BuildURL pieces together a URL from parts, escaping as necessary..
func BuildURL(endpoint Endpoint, format string, args ...string) (*url.URL, error) {
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
	default:
		return "", UnsupportedEndpointError{Selection: endpoint}
	}
	return fmt.Sprintf(format, arr...)
}
