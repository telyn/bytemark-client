package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
	"net/http"
)

// Client provides the interface which all BigV API clients should implement.
type Client interface {
	// Getters
	GetEndpoint() string
	GetSessionToken() string

	// Setters
	SetDebugLevel(int)

	// Auth
	AuthWithToken(string) error
	AuthWithCredentials(auth3.Credentials) error

	// Requests
	RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error
	RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error)
	Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error)

	// Parsers
	ParseVirtualMachineName(string) VirtualMachineName
	ParseGroupName(string) GroupName
	ParseAccountName(string) string

	// ACCOUNTS
	GetAccount(name string) (account *Account, err error)

	// GROUPS
	CreateGroup(name GroupName) error

	// VIRTUAL MACHINES
	DeleteVirtualMachine(name VirtualMachineName, purge bool) error
	GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error)
	UndeleteVirtualMachine(name VirtualMachineName) error
}
