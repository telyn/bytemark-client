package mocks

import (
	bigv "bigv.io/client/lib"
	auth3 "bytemark.co.uk/auth3/client"
	"fmt"
	mock "github.com/maraino/go-mock"
	"net/http"
)

type BigVClient struct {
	mock.Mock
}

func (c *BigVClient) GetEndpoint() string {
	r := c.Called()
	return r.String(0)
}
func (c *BigVClient) GetSessionToken() string {
	r := c.Called()
	return r.String(0)
}
func (c *BigVClient) GetSessionUser() string {
	r := c.Called()
	return r.String(0)
}
func (c *BigVClient) SetDebugLevel(level int) {
	c.Called(level)
}
func (c *BigVClient) AuthWithToken(token string) error {
	r := c.Called(token)
	return r.Error(0)
}
func (c *BigVClient) AuthWithCredentials(credents auth3.Credentials) error {
	r := c.Called(credents)
	return r.Error(0)
}
func (c *BigVClient) RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error {
	r := c.Called(auth, method, path, requestBody, output)
	return r.Error(0)
}
func (c *BigVClient) RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error) {
	r := c.Called(auth, method, path, requestBody)
	return r.Bytes(0), r.Error(1)
}

func (c *BigVClient) Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
	r := c.Called(auth, method, location, requestBody)
	req, ok := r.Get(0).(*http.Request)
	if !ok {
		panic(fmt.Sprintf("Couldn't typecast request object because it was a %t", r.Get(0)))
	}
	res, ok = r.Get(1).(*http.Response)
	if !ok {
		panic(fmt.Sprintf("Couldn't typecast response object because it was a %t", r.Get(1)))
	}
	return req, res, r.Error(2)
}

func (c *BigVClient) GetAccount(name string) (account *bigv.Account, err error) {
	r := c.Called(name)
	acc, _ := r.Get(0).(*bigv.Account)
	return acc, r.Error(1)
}

func (c *BigVClient) GetAccounts() (accounts []*bigv.Account, err error) {
	r := c.Called()
	acc, _ := r.Get(0).([]*bigv.Account)
	return acc, r.Error(1)
}

func (c *BigVClient) CreateDisc(name bigv.VirtualMachineName, disc bigv.Disc) error {
	r := c.Called(name, disc)
	return r.Error(0)
}

func (c *BigVClient) CreateGroup(name bigv.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *BigVClient) GetGroup(name bigv.GroupName) (*bigv.Group, error) {
	r := c.Called(name)
	group, _ := r.Get(0).(*bigv.Group)
	return group, r.Error(1)
}

func (c *BigVClient) DeleteDisc(name bigv.VirtualMachineName, disc int) error {
	r := c.Called(name, disc)
	return r.Error(0)
}

func (c *BigVClient) DeleteGroup(name bigv.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *BigVClient) DeleteVirtualMachine(name bigv.VirtualMachineName, purge bool) error {
	r := c.Called(name, purge)
	return r.Error(0)
}

func (c *BigVClient) CreateVirtualMachine(group bigv.GroupName, vm bigv.VirtualMachineSpec) (*bigv.VirtualMachine, error) {
	r := c.Called(group, vm)
	rvm, _ := r.Get(0).(*bigv.VirtualMachine)
	return rvm, r.Error(1)
}

func (c *BigVClient) GetVirtualMachine(name bigv.VirtualMachineName) (vm *bigv.VirtualMachine, err error) {
	r := c.Called(name)
	vm, _ = r.Get(0).(*bigv.VirtualMachine)
	return vm, r.Error(1)
}

func (c *BigVClient) ParseVirtualMachineName(name string) (bigv.VirtualMachineName, error) {
	r := c.Called(name)
	n, _ := r.Get(0).(bigv.VirtualMachineName)
	return n, r.Error(1)
}

func (c *BigVClient) ParseGroupName(name string) bigv.GroupName {
	r := c.Called(name)
	n, _ := r.Get(0).(bigv.GroupName)
	return n
}

func (c *BigVClient) ParseAccountName(name string) string {
	r := c.Called(name)
	return r.String(0)
}

func (c *BigVClient) ResetVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *BigVClient) ResizeDisc(name bigv.VirtualMachineName, id int, size int) error {
	r := c.Called(name, id, size)
	return r.Error(0)
}

func (c *BigVClient) RestartVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *BigVClient) SetVirtualMachineCores(name bigv.VirtualMachineName, cores int) error {
	r := c.Called(name, cores)
	return r.Error(0)
}
func (c *BigVClient) SetVirtualMachineHardwareProfile(name bigv.VirtualMachineName, hwprofile string, locked ...bool) error {
	r := c.Called(name, hwprofile, locked)
	return r.Error(0)
}
func (c *BigVClient) SetVirtualMachineHardwareProfileLock(name bigv.VirtualMachineName, locked bool) error {
	r := c.Called(name, locked)
	return r.Error(0)
}
func (c *BigVClient) SetVirtualMachineMemory(name bigv.VirtualMachineName, memory int) error {
	r := c.Called(name, memory)
	return r.Error(0)
}
func (c *BigVClient) StartVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *BigVClient) StopVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *BigVClient) ShutdownVirtualMachine(name bigv.VirtualMachineName, stayoff bool) error {
	r := c.Called(name, stayoff)
	return r.Error(0)
}

func (c *BigVClient) UndeleteVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
