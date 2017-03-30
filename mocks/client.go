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

func (c *Client) SetDiscIopsLimit(name *lib.VirtualMachineName, id string, size int) error {
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

func (c *Client) CreateBackup(server lib.VirtualMachineName, discLabelOrID string) (brain.Backup, error) {
	r := c.Called(server, discLabelOrID)
	snap, _ := r.Get(0).(brain.Backup)
	return snap, r.Error(1)
}
func (c *Client) DeleteBackup(server lib.VirtualMachineName, discLabelOrID string, backupLabelOrID string) error {
	r := c.Called(server, discLabelOrID, backupLabelOrID)
	return r.Error(0)
}
func (c *Client) CreateBackupSchedule(server lib.VirtualMachineName, discLabelOrID string, start string, interval int) (brain.BackupSchedule, error) {
	r := c.Called(server, discLabelOrID, start, interval)
	sched, _ := r.Get(0).(brain.BackupSchedule)
	return sched, r.Error(1)
}
func (c *Client) DeleteBackupSchedule(server lib.VirtualMachineName, discLabelOrID string, id int) error {
	r := c.Called(server, discLabelOrID, id)
	return r.Error(0)
}
func (c *Client) GetBackups(server lib.VirtualMachineName, discLabelOrID string) (brain.Backups, error) {
	r := c.Called(server, discLabelOrID)
	snaps, _ := r.Get(0).(brain.Backups)
	return snaps, r.Error(1)
}
func (c *Client) RestoreBackup(server lib.VirtualMachineName, discLabelOrID string, backupLabelOrID string) (brain.Backup, error) {
	r := c.Called(server, discLabelOrID, backupLabelOrID)
	snap, _ := r.Get(0).(brain.Backup)

	return snap, r.Error(1)
}

func (c *Client) GetPrivileges(username string) (privs brain.Privileges, err error) {
	r := c.Called(username)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GetPrivilegesForAccount(accountName string) (privs brain.Privileges, err error) {
	r := c.Called(accountName)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GetPrivilegesForGroup(group lib.GroupName) (privs brain.Privileges, err error) {
	r := c.Called(group)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GetPrivilegesForVirtualMachine(vm lib.VirtualMachineName) (privs brain.Privileges, err error) {
	r := c.Called(vm)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GrantPrivilege(priv brain.Privilege) (err error) {
	r := c.Called(priv)
	return r.Error(0)
}
func (c *Client) RevokePrivilege(priv brain.Privilege) (err error) {
	r := c.Called(priv)
	return r.Error(0)
}
func (c *Client) GetVLANs() ([]*brain.VLAN, error) {
	r := c.Called()
	vlans, _ := r.Get(0).([]*brain.VLAN)
	return vlans, r.Error(1)
}
func (c *Client) GetVLAN(num int) (*brain.VLAN, error) {
	r := c.Called(num)
	vlans, _ := r.Get(0).(*brain.VLAN)
	return vlans, r.Error(1)
}
func (c *Client) GetIPRanges() ([]*brain.IPRange, error) {
	r := c.Called()
	ipRanges, _ := r.Get(0).([]*brain.IPRange)
	return ipRanges, r.Error(1)
}
func (c *Client) GetIPRange(id int) (*brain.IPRange, error) {
	r := c.Called(id)
	ipRange, _ := r.Get(0).(*brain.IPRange)
	return ipRange, r.Error(1)
}
func (c *Client) GetHeads() ([]*brain.Head, error) {
	r := c.Called()
	heads, _ := r.Get(0).([]*brain.Head)
	return heads, r.Error(1)
}
func (c *Client) GetHead(idOrLabel string) (*brain.Head, error) {
	r := c.Called(idOrLabel)
	head, _ := r.Get(0).(*brain.Head)
	return head, r.Error(1)
}
func (c *Client) GetTails() ([]*brain.Tail, error) {
	r := c.Called()
	tails, _ := r.Get(0).([]*brain.Tail)
	return tails, r.Error(1)
}
func (c *Client) GetTail(idOrLabel string) (*brain.Tail, error) {
	r := c.Called(idOrLabel)
	tail, _ := r.Get(0).(*brain.Tail)
	return tail, r.Error(1)
}
func (c *Client) GetStoragePools() ([]*brain.StoragePool, error) {
	r := c.Called()
	storagePools, _ := r.Get(0).([]*brain.StoragePool)
	return storagePools, r.Error(1)
}
func (c *Client) GetStoragePool(idOrLabel string) (*brain.StoragePool, error) {
	r := c.Called(idOrLabel)
	storagePool, _ := r.Get(0).(*brain.StoragePool)
	return storagePool, r.Error(1)
}
func (c *Client) GetMigratingVMs() ([]*brain.VirtualMachine, error) {
	r := c.Called()
	vms, _ := r.Get(0).([]*brain.VirtualMachine)
	return vms, r.Error(1)
}
func (c *Client) GetStoppedEligibleVMs() ([]*brain.VirtualMachine, error) {
	r := c.Called()
	vms, _ := r.Get(0).([]*brain.VirtualMachine)
	return vms, r.Error(1)
}
func (c *Client) GetRecentVMs() ([]*brain.VirtualMachine, error) {
	r := c.Called()
	vms, _ := r.Get(0).([]*brain.VirtualMachine)
	return vms, r.Error(1)
}
func (c *Client) MigrateDisc(disc int, newStoragePool string) error {
	r := c.Called(disc, newStoragePool)
	return r.Error(0)
}
func (c *Client) MigrateVirtualMachine(vmName *lib.VirtualMachineName, newHead string) error {
	r := c.Called(vmName, newHead)
	return r.Error(0)
}
func (c *Client) ReapVMs() error {
	r := c.Called()
	return r.Error(0)
}
func (c *Client) DeleteVLAN(id int) error {
	r := c.Called()
	return r.Error(0)
}
func (c *Client) AdminCreateGroup(name *lib.GroupName, vlanNum int) error {
	r := c.Called(name, vlanNum)
	return r.Error(0)
}
func (c *Client) CreateIPRange(ipRange string, vlanNum int) error {
	r := c.Called(ipRange, vlanNum)
	return r.Error(0)
}
func (c *Client) CancelDiscMigration(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) CancelVMMigration(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) EmptyStoragePool(idOrLabel string) error {
	r := c.Called(idOrLabel)
	return r.Error(0)
}
func (c *Client) EmptyHead(idOrLabel string) error {
	r := c.Called(idOrLabel)
	return r.Error(0)
}
func (c *Client) ReifyDisc(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) ApproveVM(name *lib.VirtualMachineName, powerOn bool) error {
	r := c.Called(name, powerOn)
	return r.Error(0)
}
func (c *Client) RejectVM(name *lib.VirtualMachineName, reason string) error {
	r := c.Called(name, reason)
	return r.Error(0)
}
func (c *Client) RegradeDisc(disc int, newGrade string) error {
	r := c.Called(disc, newGrade)
	return r.Error(0)
}
