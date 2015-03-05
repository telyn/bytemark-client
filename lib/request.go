package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Make a request to the URL specified, giving the token stored in the auth.Client, returning the entirety of the response body.
// This is intended as the low-level work-horse of the libary, but may be deprecated in favour of MakeRequest in order to use a streaming JSON parser.
func (bigv *Client) Request(method string, location string, requestBody string) (responseBody []byte, err error) {
	req, res, err := bigv.FireRequest(method, location, requestBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if bigv.DebugLevel > 1 {
		fmt.Printf("%s: %s: %d\r\n", req.URL, res.StatusCode)
	}

	responseBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if bigv.DebugLevel > 2 {
		fmt.Printf(string(responseBody))
	}

	return responseBody, nil
}

// Make an HTTP request and then request it, returning the request object, response object and any errors
// For use by Client.MakeAndReadRequest, do not use externally except for testing
func (bigv *Client) FireRequest(method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
	url := location

	if strings.HasPrefix(location, "/") {
		url = bigv.Endpoint + location
	}
	cli := &http.Client{}

	req, err = http.NewRequest(method, url, bytes.NewBufferString(requestBody))
	if err != nil {
		return req, nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bigv.AuthSession.Token)

	res, err = cli.Do(req)
	if err != nil {
		return req, res, err
	}
	return req, res, nil
}
