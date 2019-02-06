package lib

import (
	auth3 "github.com/BytemarkHosting/auth-client"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
)

// Client provides the interface which all API clients should implement.
type Client interface {
	// Getters
	//

	// GetEndpoint returns the API endpoint currently in use.
	GetEndpoint() string

	// GetSessionFactors returns the factors provided when the current auth session was set up
	GetSessionFactors() []string
	// GetSessionToken returns the token for the current auth session
	GetSessionToken() string

	// GetSessionUser returns the user's name for the current auth session.
	GetSessionUser() string

	//
	// Setters
	//

	// SetDebugLevel sets the debug level / verbosity of the API client. 0 (default) is silent.
	SetDebugLevel(int)

	//
	// Ensurers
	//

	EnsureVirtualMachineName(vm *pathers.VirtualMachineName) error
	EnsureGroupName(group *pathers.GroupName) error
	EnsureAccountName(account *pathers.AccountName) error

	//
	// Auth
	//

	// AuthWithToken attempts to read sessiondata from auth for the given token. Returns nil on success or an error otherwise.
	AuthWithToken(string) error
	// AuthWithCredentials attempts to authenticate with the given credentials. Returns nil on success or an error otherwise.
	AuthWithCredentials(auth3.Credentials) error
	// Impersonate requests impersonation of a given user. Authentication must have already occurred.
	Impersonate(user string) error

	//
	// Requests
	//

	AllowInsecureRequests()
	BuildRequestNoAuth(method string, endpoint Endpoint, path string, parts ...string) (Request, error)
	BuildRequest(method string, endpoint Endpoint, path string, parts ...string) (Request, error)

	///////////////////////
	////// SPP STUFF //////
	///////////////////////

	CreateCreditCard(spp.CreditCard) (string, error)
	CreateCreditCardWithToken(spp.CreditCard, string) (string, error)

	///////////////////////
	//// BILLING STUFF ////
	///////////////////////

	// CreateAccount(*Account) (*Account, error) // TODO(telyn): figure out if CreateAccount is needed/useful
	GetSPPToken(cc spp.CreditCard, owner billing.Person) (string, error)
	RegisterNewAccount(acc Account) (Account, error)

	////////////////////
	//// BRAIN STUFF ///
	////////////////////

	//
	// DEFINITIONS
	//
	ReadDefinitions() (Definitions, error)

	//
	// ACCOUNTS
	//

	// GetAccount takes an account name or ID and returns a filled-out Account object
	GetAccount(name string) (account Account, err error)
	// GetDefaultAccount gets the most-likely default account for the user.
	GetDefaultAccount() (account Account, err error)
	// GetAccounts gets all the accounts the logged-in user can see.
	GetAccounts() (accounts Accounts, err error)

	//
	// DISCS
	//

	CreateDisc(vm pathers.VirtualMachineName, disc brain.Disc) error
	DeleteDisc(vm pathers.VirtualMachineName, idOrLabel string) error
	GetDisc(vm pathers.VirtualMachineName, idOrLabel string) (brain.Disc, error)
	GetDiscByID(id int) (brain.Disc, error)
	ResizeDisc(vm pathers.VirtualMachineName, idOrLabel string, size int) error
	SetDiscIopsLimit(vm pathers.VirtualMachineName, idOrLabel string, iopsLimit int) error

	//
	// GROUPS
	//

	// CreateGroup sends a request to the API server to create a group with the given name.
	CreateGroup(name pathers.GroupName) error
	DeleteGroup(name pathers.GroupName) error
	GetGroup(name pathers.GroupName) (brain.Group, error)

	//
	// NICS
	//

	AddIP(name pathers.VirtualMachineName, ipcr brain.IPCreateRequest) (brain.IPs, error)

	//
	// PRIVILEGES
	//
	// username is allowed to be empty
	GetPrivileges(username string) (brain.Privileges, error)
	GetPrivilegesForAccount(account string) (brain.Privileges, error)
	GetPrivilegesForGroup(group pathers.GroupName) (brain.Privileges, error)
	GetPrivilegesForVirtualMachine(vm pathers.VirtualMachineName) (brain.Privileges, error)
	GrantPrivilege(p brain.Privilege) error
	RevokePrivilege(p brain.Privilege) error

	//
	// BACKUPS
	//

	CreateBackup(server pathers.VirtualMachineName, discLabelOrID string) (brain.Backup, error)
	DeleteBackup(server pathers.VirtualMachineName, discLabelOrID string, backupLabelOrID string) error
	GetBackups(server pathers.VirtualMachineName, discLabelOrID string) (brain.Backups, error)
	RestoreBackup(server pathers.VirtualMachineName, discLabelOrID string, backupLabelOrID string) (brain.Backup, error)

	//
	// BACKUP SCHEDULES
	//

	CreateBackupSchedule(server pathers.VirtualMachineName, discLabel string, startDate string, intervalSeconds int) (brain.BackupSchedule, error)
	DeleteBackupSchedule(server pathers.VirtualMachineName, discLabel string, id int) error

	//
	// USERS
	//

	GetUser(name string) (brain.User, error)

	//
	// VIRTUAL MACHINES
	//

	// CreateVirtualMachine creates a virtual machine with a given specification in the given group.
	// returns nil on success or an error otherwise.
	CreateVirtualMachine(group pathers.GroupName, vm brain.VirtualMachineSpec) (brain.VirtualMachine, error)

	// DeleteVirtualMachine deletes the named virtual machine.
	// returns nil on success or an error otherwise.
	DeleteVirtualMachine(name pathers.VirtualMachineName, purge bool) error

	// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
	GetVirtualMachine(name pathers.VirtualMachineName) (brain.VirtualMachine, error)

	// MoveVirtualMachine moves a server from one pathers.VirtualMachineName to another.
	MoveVirtualMachine(old pathers.VirtualMachineName, new pathers.VirtualMachineName) error

	// ReimageVirtualMachine reimages the named virtual machine. This will wipe everything on the first disk in the vm and install a new OS on top of it.
	// Note that the machine in question must already be powered off. Once complete, according to the API docs, the vm will be powered on but its autoreboot_on will be false.
	ReimageVirtualMachine(name pathers.VirtualMachineName, image brain.ImageInstall) (err error)

	// ResetVirtualMachine resets the named virtual machine. This is like pressing the reset
	// button on a physical computer. This does not cause a new process to be started, so does not apply any pending hardware changes.
	// returns nil on success or an error otherwise.
	ResetVirtualMachine(name pathers.VirtualMachineName) (err error)

	// RestartVirtualMachine restarts the named virtual machine. This is
	// returns nil on success or an error otherwise.
	RestartVirtualMachine(name pathers.VirtualMachineName) (err error)

	// StartVirtualMachine starts the named virtual machine.
	// returns nil on success or an error otherwise.
	StartVirtualMachine(name pathers.VirtualMachineName) (err error)

	// StopVirtualMachine starts the named virtual machine.
	// returns nil on success or an error otherwise.
	StopVirtualMachine(name pathers.VirtualMachineName) (err error)

	// ShutdownVirtualMachine sends an ACPI shutdown to the VM. This will cause a graceful shutdown of the machine
	// returns nil on success or an error otherwise.
	ShutdownVirtualMachine(name pathers.VirtualMachineName, stayoff bool) (err error)

	// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
	// Return nil on success, an error otherwise.
	UndeleteVirtualMachine(name pathers.VirtualMachineName) error

	// SetVirtualMachineHardwareProfile specifies the hardware profile on a VM. Optionally locks or unlocks h. profile
	// Return nil on success, an error otherwise.
	SetVirtualMachineHardwareProfile(name pathers.VirtualMachineName, profile string, locked ...bool) (err error)

	// SetVirtualMachineHardwareProfileLock locks or unlocks the hardware profile of a VM.
	// Return nil on success, an error otherwise.
	SetVirtualMachineHardwareProfileLock(name pathers.VirtualMachineName, locked bool) (err error)

	// SetVirtualMachineMemory sets the RAM available to a virtual machine in megabytes
	// Return nil on success, an error otherwise.
	SetVirtualMachineMemory(name pathers.VirtualMachineName, memory int) (err error)

	// SetVirtualMachineCores sets the number of CPUs available to a virtual machine
	// Return nil on success, an error otherwise.
	SetVirtualMachineCores(name pathers.VirtualMachineName, cores int) (err error)

	// SetVirtualMachineCDROM sets the URL of a CD to attach to a virtual machine. Set url to "" to remove the CD.
	// Returns nil on success, an error otherwise.
	SetVirtualMachineCDROM(name pathers.VirtualMachineName, url string) (err error)

	//
	// ADMIN
	//

	GetVLANs() (brain.VLANs, error)
	GetVLAN(num int) (brain.VLAN, error)
	GetIPRanges() (brain.IPRanges, error)
	GetIPRange(idOrCIDR string) (brain.IPRange, error)
	GetHeads() (brain.Heads, error)
	GetHead(idOrLabel string) (brain.Head, error)
	GetTails() (brain.Tails, error)
	GetTail(idOrLabel string) (brain.Tail, error)
	GetStoragePools() (brain.StoragePools, error)
	GetStoragePool(idOrLabel string) (brain.StoragePool, error)
	GetMigratingDiscs() (brain.Discs, error)
	GetMigratingVMs() (brain.VirtualMachines, error)
	GetStoppedEligibleVMs() (brain.VirtualMachines, error)
	GetRecentVMs() (brain.VirtualMachines, error)
	MigrateDisc(disc int, newStoragePool string) error
	MigrateVirtualMachine(vmName pathers.VirtualMachineName, newHead string) error
	DeleteVLAN(id int) error
	AdminCreateGroup(name pathers.GroupName, vlanNum int) error
	CreateIPRange(ipRange string, vlanNum int) error
	CancelDiscMigration(id int) error
	CancelVMMigration(id int) error
	EmptyStoragePool(idOrLabel string) error
	EmptyHead(idOrLabel string) error
	RegradeDisc(disc int, newGrade string) error
	UpdateVMMigration(name pathers.VirtualMachineName, speed *int64, downtime *int) error
	CreateUser(username string, privilege string) error
	UpdateHead(idOrLabel string, options UpdateHead) error
	UpdateTail(idOrLabel string, options UpdateTail) error
	UpdateStoragePool(idOrLabel string, options brain.StoragePool) error
}
