//
package main

import (
	auth "bytemark.co.uk/auth3/client"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func new_session(token_path string, auth_client *auth.Client) string {
	bpass, err := ioutil.ReadFile("pass")
	check(err)

	pass := strings.TrimSpace(string(bpass))

	creds := make(auth.Credentials)
	creds["username"] = "myuser"
	creds["password"] = pass

	session, err := auth_client.CreateSession(creds)
	check(err)

	token := session.Token
	err = ioutil.WriteFile(token_path, []byte(session.Token), 0700)
	check(err)

	return token
}

func main() {
	config_dir := filepath.Join(os.Getenv("HOME"), ".bigv-go")
	token_path := filepath.Join(config_dir, "token")
	url := "https://uk0.bigv.io/accounts/telyn/groups/work/virtual_machines"

	err := os.MkdirAll(config_dir, 0700)
	check(err)

	auth_client, err := auth.New("https://auth.bytemark.co.uk")
	check(err)

	token_bytes, err := ioutil.ReadFile(token_path)
	token := ""

	if err != nil && os.IsNotExist(err) {
		// get a token
		new_session(token_path, auth_client)
	} else if err != nil {
		check(err)
	} else {
		token = string(token_bytes)
		session, err := auth_client.ReadSession(token)
		if err != nil {
			// really the response should depend on what kind of error we get but atm I don't know auth3 well enough.
			token = new_session(token_path, auth_client)
		} else {
			token = session.Token
		}
	}

	cli := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token+"")

	res, err := cli.Do(req)
	defer res.Body.Close()

	fmt.Printf("%s returned a %d\r\n", url, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	fmt.Printf(string(body))

	os.Exit(0)
}
