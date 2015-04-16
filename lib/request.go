package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// RequestAndUnmarshal performs a request (with no body) and unmarshals the result into output - which should be a pointer to something cool
func (bigv *BigVClient) RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error {

	data, err := bigv.RequestAndRead(auth, method, path, requestBody)

	if bigv.DebugLevel >= 4 {
		fmt.Printf("'%s'\r\n", data)
	}

	if err != nil {
		//TODO(telyn): good error handling here (need to see more errors first)
		return err
	}

	err = json.Unmarshal(data, output)
	if err != nil {
		// BUG(telyn): this is a bad error message and you should feel bad
		fmt.Printf("Data returned was not the right type.\r\n")
		fmt.Printf("%+v\r\n", output)

		return err
	}

	if bigv.DebugLevel >= 3 {
		buf := new(bytes.Buffer)
		json.Indent(buf, data, "", "    ")
		fmt.Printf("%s", buf)
	}

	return nil

}

// Request makes a request to the URL specified, giving the token stored in the auth.Client, returning the entirety of the response body.
// This is intended as the low-level work-horse of the libary, but may be deprecated in favour of MakeRequest in order to use a streaming JSON parser.
func (bigv *BigVClient) RequestAndRead(auth bool, method, location, requestBody string) (responseBody []byte, err error) {
	req, res, err := bigv.Request(auth, method, location, requestBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if bigv.DebugLevel > 1 {
		fmt.Printf("%s %s: %d\r\n", method, req.URL, res.StatusCode)
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
// For use by Client.RequestAndRead, do not use externally except for testing
func (bigv *BigVClient) Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
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
	if auth {
		req.Header.Add("Authorization", "Bearer "+bigv.AuthSession.Token)
	}

	res, err = cli.Do(req)
	if err != nil {
		return req, res, err
	}
	return req, res, nil
}
