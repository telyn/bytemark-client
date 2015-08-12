package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
	"net/http"
)

// Client provides the interface which all BigV API clients should implement.
type Client interface {
	// Getters
	//

	// GetEndpoint returns the BigV endpoint currently in use.
	GetEndpoint() string

	// GetSessionToken returns the token for the current auth session - note that this may cause panics at this time.
	GetSessionToken() string

	//
	// Setters
	//

	// SetDebugLevel sets the debug level / verbosity of the BigV client. 0 (default) is silent.
	SetDebugLevel(int)

	//
	// Auth
	//

	// AuthWithToken attempts to read sessiondata from auth for the given token. Returns nil on success or an error otherwise.
	AuthWithToken(string) error
	// AuthWithCredentials attempts to authenticate with the given credentials. Returns nil on success or an error otherwise.
	AuthWithCredentials(auth3.Credentials) error

	//
	// Requests
	//

	// RequestAndUnmarshal performs a request (with no body) and unmarshals the result into output - which should be a pointer to something cool
	RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error

	// RequestAndRead makes a request to the URL specified, giving the token stored in the auth.Client, returning the entirety of the response body.
	// Use RequestAndUnmarshal unless you know that the given URL doesn't return JSON - this method may be unexported in future releases.
	RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error)

	// Request makes an HTTP request and then request it, returning the request object, response object and any errors
	// For use by Client.RequestAndRead, do not use externally except for testing
	Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error)

	//
	// Parsers
	//

	ParseVirtualMachineName(string) VirtualMachineName
	ParseGroupName(string) GroupName
	ParseAccountName(string) string

	//
	// ACCOUNTS
	//

	// GetAccount takes an account name or ID and returns a filled-out Account object
	GetAccount(name string) (account *Account, err error)

	//
	// GROUPS
	//

	// CreateGroup sends a request to the BigV server to create a group with the given name.
	CreateGroup(name GroupName) error
	DeleteGroup(name GroupName) error
	GetGroup(name GroupName) (*Group, error)

	//
	// VIRTUAL MACHINES
	//

	// CreateVirtualMachine creates a virtual machine with a given specification in the given group.
	// returns nil on success or an error otherwise.
	CreateVirtualMachine(group GroupName, vm VirtualMachineSpec) (*VirtualMachine, error)

	// DeleteVirtualMachine deletes the named virtual machine.
	// returns nil on success or an error otherwise.
	DeleteVirtualMachine(name VirtualMachineName, purge bool) error

	// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
	GetVirtualMachine(name VirtualMachineName) (*VirtualMachine, error)

	// ResetVirtualMachine resets the named virtual machine. This is like pressing the reset
	// button on a physical computer. This does not cause a new process to be started, so does not apply any pending hardware changes.
	// returns nil on success or an error otherwise.
	ResetVirtualMachine(name VirtualMachineName) (err error)

	// RestartVirtualMachine restarts the named virtual machine. This is
	// returns nil on success or an error otherwise.
	RestartVirtualMachine(name VirtualMachineName) (err error)

	// StartVirtualMachine starts the named virtual machine.
	// returns nil on success or an error otherwise.
	StartVirtualMachine(name VirtualMachineName) (err error)

	// StopVirtualMachine starts the named virtual machine.
	// returns nil on success or an error otherwise.
	StopVirtualMachine(name VirtualMachineName) (err error)

	// ShutdownVirtualMachine sends an ACPI shutdown to the VM. This will cause a graceful shutdown of the machine
	// returns nil on success or an error otherwise.
	ShutdownVirtualMachine(name VirtualMachineName, stayoff bool) (err error)

	// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
	// Return nil on success, an error otherwise.
	UndeleteVirtualMachine(name VirtualMachineName) error

	// SetVirtualMachineHardwareProfile specifies the hardware profile on a VM. Optionally locks or unlocks h. profile
	// Return nil on success, an error otherwise.
	SetVirtualMachineHardwareProfile(name VirtualMachine, profile string, locked ...bool) (err error)

	// SetVirtualMachineHardwareProfileLock locks or unlocks the hardware profile of a VM.
	// Return nil on success, an error otherwise.
	SetVirtualMachineHardwareProfileLock(name VirtualMachine, locked bool) (err error)
}
