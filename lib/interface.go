package lib

import (
	"net/http"
)

type Client interface {
	RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error
	RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error)
	Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error)

	// ACCOUNTS
	GetAccount(name string) (account *Account, err error)

	// GROUPS

	// VIRTUAL MACHINES
	GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error)
}
