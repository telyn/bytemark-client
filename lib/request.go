package lib

import (
	"bigv.io/client/util/log"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// BuildURL pieces together a URL from parts, escaping as necessary..
func BuildURL(format string, args ...string) string {
	arr := make([]interface{}, len(args), len(args))
	for i, str := range args {
		arr[i] = url.QueryEscape(str)
	}
	return fmt.Sprintf(format, arr...)
}

// RequestAndUnmarshal performs a request (with no body) and unmarshals the result into output - which should be a pointer to something cool
func (bigv *bigvClient) RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error {

	data, err := bigv.RequestAndRead(auth, method, path, requestBody)

	if bigv.debugLevel >= 3 {
		buf := new(bytes.Buffer)
		json.Indent(buf, data, "", "    ")
		log.Debugf(3, "%s", buf)
	}

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, output)
	return err

}

// RequestAndRead makes a request to the URL specified, giving the token stored in the auth.Client, returning the entirety of the response body.
// This is intended as the low-level work-horse of the libary, but may be deprecated in favour of MakeRequest in order to use a streaming JSON parser.
func (bigv *bigvClient) RequestAndRead(auth bool, method, location, requestBody string) (responseBody []byte, err error) {
	_, res, err := bigv.Request(auth, method, location, requestBody)
	defer res.Body.Close()

	responseBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	log.Debugf(2, "response body: '%s'\r\n", string(responseBody))

	return responseBody, err
}

// Request makes an HTTP request and then request it, returning the request object, response object and any errors
// For use by Client.RequestAndRead, do not use externally except for testing
func (bigv *bigvClient) Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
	url := location

	if strings.HasPrefix(location, "/") {
		url = bigv.endpoint + location
	}
	cli := &http.Client{}

	req, err = http.NewRequest(method, url, bytes.NewBufferString(requestBody))
	if err != nil {
		return req, nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if auth {
		req.Header.Add("Authorization", "Bearer "+bigv.authSession.Token)
	}

	res, err = cli.Do(req)
	if err != nil {
		return req, res, err
	}

	log.Debugf(1, "%s %s: %d\r\n", method, req.URL, res.StatusCode)
	log.Debugf(3, "request body: '%s'\r\n", requestBody)

	baseErr := BigVError{
		Method:       method,
		URL:          req.URL,
		StatusCode:   res.StatusCode,
		RequestBody:  requestBody,
		ResponseBody: "",
	}

	switch res.StatusCode {
	case 400:
		err = BadRequestError{baseErr}

	case 403:
		err = NotAuthorizedError{baseErr}
	case 404:
		err = NotFoundError{baseErr}
	default:
		if 200 <= res.StatusCode && res.StatusCode <= 299 {
			break
		}
		err = UnknownStatusCodeError{baseErr}
	}
	return req, res, err
}
