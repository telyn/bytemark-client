package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Make a request to the URL specified, giving the token stored in the auth.Client.
// This is intended as the low-level work-horse of the libary

func (bigv *Client) MakeRequest(method string, location string, requestBody string) (responseBody string, err error) {
	url := location

	if strings.HasPrefix(location, "/") {
		url = bigv.Endpoint + location
	}
	cli := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bigv.AuthSession.Token)

	res, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if bigv.DebugLevel > 1 {
		fmt.Printf("%s: %s: %d\r\n", url, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if bigv.DebugLevel > 2 {
		fmt.Printf(string(body))
	}

	return string(body), nil
}
