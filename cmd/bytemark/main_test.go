package main

import (
	"bytemark.co.uk/client/mocks"
)

func baseTestSetup() (config *mocks.Config, client *mocks.Client) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	global.Client = client
	global.Config = config
	global.Error = nil

	baseAppSetup()
	return
}
