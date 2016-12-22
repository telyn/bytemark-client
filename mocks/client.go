package mocks

import (
	"fmt"
	auth3 "github.com/BytemarkHosting/auth-client"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	mock "github.com/maraino/go-mock"
	"net/http"
	"net/url"
)

type Client struct {
	mock.Mock
}

func (c *Client) AllowInsecureRequests() {
	c.Called()
}
func (c *Client) BuildRequestNoAuth(method string, endpoint lib.Endpoint, path string, parts ...string) (*lib.Request, error) {
	r := c.Called(method, endpoint, path, parts)
	req, _ := r.Get(0).(*lib.Request)
	return req, r.Error(1)
}
func (c *Client) BuildRequest(method string, endpoint lib.Endpoint, path string, parts ...string) (*lib.Request, error) {
	r := c.Called(method, endpoint, path, parts)
	req, _ := r.Get(0).(*lib.Request)
	return req, r.Error(1)

}
func (c *Client) NewRequestNoAuth(method string, url *url.URL) *lib.Request {
	r := c.Called(method, url)
	req, _ := r.Get(0).(*lib.Request)
	return req
}
func (c *Client) NewRequest(method string, url *url.URL) *lib.Request {
	r := c.Called(method, url)
	req, _ := r.Get(0).(*lib.Request)
	return req
}

func (c *Client) GetEndpoint() string {
	r := c.Called()
	return r.String(0)
}
func (c *Client) GetSessionToken() string {
	r := c.Called()
	return r.String(0)
}
func (c *Client) GetSessionUser() string {
	r := c.Called()
	return r.String(0)
}
func (c *Client) GetSessionFactors() []string {
	r := c.Called()
	ar := r.Get(0)
	return ar.([]string)
}
func (c *Client) GetSPPToken(cc spp.CreditCard, owner *billing.Person) (string, error) {
	r := c.Called(cc, owner)
	return r.String(0), r.Error(1)
}
func (c *Client) SetDebugLevel(level int) {
	c.Called(level)
}
func (c *Client) AuthWithToken(token string) error {
	r := c.Called(token)
	return r.Error(0)
}
func (c *Client) AuthWithCredentials(credents auth3.Credentials) error {
	r := c.Called(credents)
	return r.Error(0)
}
func (c *Client) RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error {
	r := c.Called(auth, method, path, requestBody, output)
	return r.Error(0)
}
func (c *Client) RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error) {
	r := c.Called(auth, method, path, requestBody)
	return r.Bytes(0), r.Error(1)
}

func (c *Client) Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
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

func (c *Client) ReadDefinitions() (*lib.Definitions, error) {
	r := c.Called()
	defs, _ := r.Get(0).(*lib.Definitions)
	return defs, r.Error(1)
}

func (c *Client) AddIP(name *lib.VirtualMachineName, spec *brain.IPCreateRequest) (brain.IPs, error) {
	r := c.Called(name, spec)
	ips, _ := r.Get(0).(brain.IPs)
	return ips, r.Error(1)
}

func (c *Client) AddUserAuthorizedKey(name, key string) error {
	r := c.Called(name, key)
	return r.Error(0)
}

func (c *Client) DeleteUserAuthorizedKey(name, key string) error {
	r := c.Called(name, key)
	return r.Error(0)
}

func (c *Client) GetUser(name string) (*brain.User, error) {
	r := c.Called(name)
	u, _ := r.Get(0).(*brain.User)
	return u, r.Error(1)
}

func (c *Client) CreateCreditCard(cc *spp.CreditCard) (string, error) {
	r := c.Called(cc)
	return r.String(0), r.Error(1)
}
func (c *Client) CreateCreditCardWithToken(cc *spp.CreditCard, token string) (string, error) {
	r := c.Called(cc, token)
	return r.String(0), r.Error(1)
}
func (c *Client) CreateAccount(acc *lib.Account) (*lib.Account, error) {
	r := c.Called(acc)
	a, _ := r.Get(0).(*lib.Account)
	return a, r.Error(1)
}

func (c *Client) RegisterNewAccount(acc *lib.Account) (*lib.Account, error) {
	r := c.Called(acc)
	a, _ := r.Get(0).(*lib.Account)
	return a, r.Error(1)
}

func (c *Client) GetAccount(name string) (account *lib.Account, err error) {
	r := c.Called(name)
	acc, _ := r.Get(0).(*lib.Account)
	return acc, r.Error(1)
}

func (c *Client) GetDefaultAccount() (account *lib.Account, err error) {
	r := c.Called()
	acc, _ := r.Get(0).(*lib.Account)
	return acc, r.Error(1)
}

func (c *Client) GetAccounts() (accounts []*lib.Account, err error) {
	r := c.Called()
	acc, _ := r.Get(0).([]*lib.Account)
	return acc, r.Error(1)
}

func (c *Client) CreateDisc(name *lib.VirtualMachineName, disc brain.Disc) error {
	r := c.Called(name, disc)
	return r.Error(0)
}

func (c *Client) GetDisc(name *lib.VirtualMachineName, discId string) (disc *brain.Disc, err error) {
	r := c.Called(name, discId)
	disc, _ = r.Get(0).(*brain.Disc)
	return disc, r.Error(1)
}

func (c *Client) CreateGroup(name *lib.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) GetGroup(name *lib.GroupName) (*brain.Group, error) {
	r := c.Called(name)
	group, _ := r.Get(0).(*brain.Group)
	return group, r.Error(1)
}

func (c *Client) DeleteDisc(name *lib.VirtualMachineName, disc string) error {
	r := c.Called(name, disc)
	return r.Error(0)
}

func (c *Client) DeleteGroup(name *lib.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) DeleteVirtualMachine(name *lib.VirtualMachineName, purge bool) error {
	r := c.Called(name, purge)
	return r.Error(0)
}

func (c *Client) CreateVirtualMachine(group *lib.GroupName, vm brain.VirtualMachineSpec) (*brain.VirtualMachine, error) {
	r := c.Called(group, vm)
	rvm, _ := r.Get(0).(*brain.VirtualMachine)
	return rvm, r.Error(1)
}

func (c *Client) GetVirtualMachine(name *lib.VirtualMachineName) (vm *brain.VirtualMachine, err error) {
	r := c.Called(name)
	vm, _ = r.Get(0).(*brain.VirtualMachine)
	return vm, r.Error(1)
}

func (c *Client) MoveVirtualMachine(oldName *lib.VirtualMachineName, newName *lib.VirtualMachineName) error {
	r := c.Called(oldName, newName)
	return r.Error(0)
}

func (c *Client) ParseVirtualMachineName(name string, defaults ...*lib.VirtualMachineName) (*lib.VirtualMachineName, error) {
	r := c.Called(name, defaults)
	n, _ := r.Get(0).(*lib.VirtualMachineName)
	return n, r.Error(1)
}

func (c *Client) ParseGroupName(name string, defaults ...*lib.GroupName) *lib.GroupName {
	r := c.Called(name, defaults)
	n, _ := r.Get(0).(*lib.GroupName)
	return n
}

func (c *Client) ParseAccountName(name string, defaults ...string) string {
	r := c.Called(name, defaults)
	return r.String(0)
}

func (c *Client) ReimageVirtualMachine(name *lib.VirtualMachineName, image *brain.ImageInstall) error {
	r := c.Called(name, image)
	return r.Error(0)
}

func (c *Client) ResetVirtualMachine(name *lib.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) ResizeDisc(name *lib.VirtualMachineName, id string, size int) error {
	r := c.Called(name, id, size)
	return r.Error(0)
}

func (c *Client) RestartVirtualMachine(name *lib.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineCDROM(name *lib.VirtualMachineName, url string) error {
	r := c.Called(name, url)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineCores(name *lib.VirtualMachineName, cores int) error {
	r := c.Called(name, cores)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineHardwareProfile(name *lib.VirtualMachineName, hwprofile string, locked ...bool) error {
	r := c.Called(name, hwprofile, locked)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineHardwareProfileLock(name *lib.VirtualMachineName, locked bool) error {
	r := c.Called(name, locked)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineMemory(name *lib.VirtualMachineName, memory int) error {
	r := c.Called(name, memory)
	return r.Error(0)
}
func (c *Client) StartVirtualMachine(name *lib.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *Client) StopVirtualMachine(name *lib.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *Client) ShutdownVirtualMachine(name *lib.VirtualMachineName, stayoff bool) error {
	r := c.Called(name, stayoff)
	return r.Error(0)
}

func (c *Client) UndeleteVirtualMachine(name *lib.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
