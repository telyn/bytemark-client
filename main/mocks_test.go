package main

import (
	bigv "bigv.io/client/lib"
	auth3 "bytemark.co.uk/auth3/client"
	"flag"
	"fmt"
	mock "github.com/maraino/go-mock"
	"net/http"
)

// mock Config

type mockConfig struct {
	mock.Mock
}

func (c *mockConfig) EndpointName() string {
	ret := c.Called()
	return ret.String(0)
}
func (c *mockConfig) Force() bool {
	ret := c.Called()
	return ret.Bool(0)
}

func (c *mockConfig) Get(name string) (string, error) {
	ret := c.Called(name)
	return ret.String(0), ret.Error(1)
}
func (c *mockConfig) GetIgnoreErr(name string) string {
	ret := c.Called(name)
	return ret.String(0)
}

func (c *mockConfig) GetBool(name string) (bool, error) {
	ret := c.Called(name)
	return ret.Bool(0), ret.Error(1)
}

func (c *mockConfig) GetV(name string) (ConfigVar, error) {
	ret := c.Called(name)
	return ret.Get(0).(ConfigVar), ret.Error(1)
}

func (c *mockConfig) GetAll() ([]ConfigVar, error) {
	ret := c.Called()
	return ret.Get(0).([]ConfigVar), ret.Error(1)
}

func (c *mockConfig) PanelURL() string {
	ret := c.Called()
	return ret.String(0)
}

func (c *mockConfig) Set(name, value, source string) {
	c.Called(name, value, source)
	return
}

func (c *mockConfig) SetPersistent(name, value, source string) error {
	ret := c.Called(name, value, source)
	return ret.Error(0)
}
func (c *mockConfig) Silent() bool {
	ret := c.Called()
	return ret.Bool(0)
}

func (c *mockConfig) Unset(name string) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *mockConfig) GetDebugLevel() int {
	ret := c.Called()
	return ret.Int(0)
}

func (c *mockConfig) ImportFlags(*flag.FlagSet) []string {
	ret := c.Called()
	if arr, ok := ret.Get(0).([]string); ok {
		return arr
	}
	return nil
}

// mock CommandSet

type mockCommands struct {
	mock.Mock
}

func (cmds *mockCommands) Console(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) Debug(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) DeleteVM(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}
func (cmds *mockCommands) DeleteGroup(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) Help(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) Config(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}
func (cmds *mockCommands) CreateGroup(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}
func (cmds *mockCommands) CreateVM(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) ResetVM(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}
func (cmds *mockCommands) Restart(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))

}
func (cmds *mockCommands) Shutdown(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}
func (cmds *mockCommands) Stop(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}
func (cmds *mockCommands) Start(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) ShowAccount(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) ShowGroup(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) ShowVM(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) UndeleteVM(args []string) ExitCode {
	r := cmds.Called(args)
	return ExitCode(r.Int(0))
}

func (cmds *mockCommands) EnsureAuth() error {
	r := cmds.Called()
	return r.Error(0)
}

func (cmds *mockCommands) HelpForConfig() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForCreate() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForDebug() {
	cmds.Called()
}
func (cmds *mockCommands) HelpForDelete() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForHelp() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForPower() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForShow() {
	cmds.Called()
}

type mockBigVClient struct {
	mock.Mock
}

func (c *mockBigVClient) GetEndpoint() string {
	r := c.Called()
	return r.String(0)
}
func (c *mockBigVClient) GetSessionToken() string {
	r := c.Called()
	return r.String(0)
}
func (c *mockBigVClient) GetSessionUser() string {
	r := c.Called()
	return r.String(0)
}
func (c *mockBigVClient) SetDebugLevel(level int) {
	c.Called(level)
}
func (c *mockBigVClient) AuthWithToken(token string) error {
	r := c.Called(token)
	return r.Error(0)
}
func (c *mockBigVClient) AuthWithCredentials(credents auth3.Credentials) error {
	r := c.Called(credents)
	return r.Error(0)
}
func (c *mockBigVClient) RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error {
	r := c.Called(auth, method, path, requestBody, output)
	return r.Error(0)
}
func (c *mockBigVClient) RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error) {
	r := c.Called(auth, method, path, requestBody)
	return r.Bytes(0), r.Error(1)
}

func (c *mockBigVClient) Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
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

func (c *mockBigVClient) GetAccount(name string) (account *bigv.Account, err error) {
	r := c.Called(name)
	acc, _ := r.Get(0).(*bigv.Account)
	return acc, r.Error(1)
}

func (c *mockBigVClient) CreateGroup(name bigv.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) GetGroup(name bigv.GroupName) (*bigv.Group, error) {
	r := c.Called(name)
	group, _ := r.Get(0).(*bigv.Group)
	return group, r.Error(1)
}
func (c *mockBigVClient) DeleteGroup(name bigv.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) DeleteVirtualMachine(name bigv.VirtualMachineName, purge bool) error {
	r := c.Called(name, purge)
	return r.Error(0)
}
func (c *mockBigVClient) CreateVirtualMachine(group bigv.GroupName, vm bigv.VirtualMachineSpec) (*bigv.VirtualMachine, error) {
	r := c.Called(group, vm)
	rvm, _ := r.Get(0).(*bigv.VirtualMachine)
	return rvm, r.Error(1)
}
func (c *mockBigVClient) GetVirtualMachine(name bigv.VirtualMachineName) (vm *bigv.VirtualMachine, err error) {
	r := c.Called(name)
	vm, _ = r.Get(0).(*bigv.VirtualMachine)
	return vm, r.Error(1)
}
func (c *mockBigVClient) UndeleteVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) ResetVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) RestartVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) SetVirtualMachineHardwareProfile(name bigv.VirtualMachineName, hwprofile string, locked ...bool) error {
	r := c.Called(name, hwprofile, locked)
	return r.Error(0)
}
func (c *mockBigVClient) SetVirtualMachineHardwareProfileLock(name bigv.VirtualMachineName, locked bool) error {
	r := c.Called(name, locked)
	return r.Error(0)
}
func (c *mockBigVClient) StartVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) StopVirtualMachine(name bigv.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *mockBigVClient) ShutdownVirtualMachine(name bigv.VirtualMachineName, stayoff bool) error {
	r := c.Called(name, stayoff)
	return r.Error(0)
}

func (c *mockBigVClient) ParseVirtualMachineName(name string) bigv.VirtualMachineName {
	r := c.Called(name)
	n, _ := r.Get(0).(bigv.VirtualMachineName)
	return n
}
func (c *mockBigVClient) ParseGroupName(name string) bigv.GroupName {
	r := c.Called(name)
	n, _ := r.Get(0).(bigv.GroupName)
	return n
}
func (c *mockBigVClient) ParseAccountName(name string) string {
	r := c.Called(name)
	return r.String(0)
}
