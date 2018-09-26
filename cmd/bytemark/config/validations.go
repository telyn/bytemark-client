package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

func (c *config) validateEndpoint(endpoint string) error {
	url, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return errors.New("The endpoint URL should start with http:// or https://")
	}
	if url.Host == "" {
		return errors.New("The endpoint URL should have a hostname")
	}
	return nil
}

func (c *config) validateAccount(client lib.Client, name string) (err error) {
	_, err = client.GetAccount(name)
	if err != nil {
		if _, ok := err.(lib.NotFoundError); ok {
			return fmt.Errorf("No such account %s - check your typing and specify --yubikey if necessary", name)
		}
		return err
	}
	return
}

func (c *config) validateGroup(client lib.Client, name string) (err error) {
	groupName := lib.ParseGroupName(name, c.GetGroup())
	_, err = client.GetGroup(groupName)
	if err != nil {
		if _, ok := err.(lib.NotFoundError); ok {
			return fmt.Errorf("No such group %v - check your typing and specify --yubikey if necessary", groupName)
		}
		return err
	}
	return
}

func (c *config) Validate(client lib.Client, varname string, value string) error {
	switch varname {
	case "endpoint", "api-endpoint", "billing-endpoint", "spp-endpoint", "auth-endpoint":
		return c.validateEndpoint(value)
	case "account":
		return c.validateAccount(client, value)
	case "group":
		return c.validateGroup(client, value)
	case "debug-level":
		_, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return errors.New("debug-level must be an integer")
		}
	}
	return nil
}
