package brain

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
	"strings"
)

// PrivilegeLevel is a type to represent different privilege levels.
// since privilege levels in the brain are just strings, they're just a string type here too.
type PrivilegeLevel string

const (
	// ClusterAdminPrivile allows a user to administer the cluster managed by the brain, and do things like create/delete VMs on accounts they have no explicit right on, grant others AccountAdminPrivilege, and set disc iops_limit
	ClusterAdminPrivilege PrivilegeLevel = "cluster_admin"
	// AccountAdminPrivilege allows a user to create, modify & delete groups and servers in an account.
	AccountAdminPrivilege = "account_admin"
	// GroupAdminPrivilege allows a user to create, modify & delete servers in a specific group.
	GroupAdminPrivilege = "group_admin"
	// VMAdminPrivilege allows a user to modify & administer a server, including increasing the performance (and hence the price on the uk0 cluster) and accessing the console.
	VMAdminPrivilege = "vm_admin"
	// VMConsolePrivilege allows a user to access the console for a particular server.
	VMConsolePrivilege = "vm_console"
)

// String returns the privilege level cast to a string.
func (pl PrivilegeLevel) String() string {
	return string(pl)
}

// Privilege represents a privilege on the brain.
// A user may have multiple privileges, and multiple privileges may be granted on the same object.
// At the moment we're not worried about the extra fields that privileges have on the brain (IP restrictions) because they're unused
type Privilege struct {
	// ID is the numeric ID used mostly by the brain
	ID int `json:"id,omitempty"`
	// Username is the user who the privilege is granted to
	Username string `json:"username,omitempty"`
	// VirtualMachineID is the ID of the virtual machine the privilege is granted on
	VirtualMachineID int `json:"virtual_machine_id,omitempty"`
	// AccountID is the ID of the account the privilege is granted on
	AccountID int `json:"account_id,omitempty"`
	// GroupID is the ID of the group the privilege is granted on
	GroupID int `json:"group_id,omitempty"`
	// Level is the PrivilegeLevel they have
	Level PrivilegeLevel `json:"level,omitempty"`
	// YubikeyRequired is true if the user should have to authenticate with a yubikey in order to use this privilege. Only set it to true if you're sure the user has a yubikey set up on their account, and that they know where it is!
	YubikeyRequired bool `json:"yubikey_required,omitempty"`
	// YubikeyOTPMaxAge should set how long (in seconds) a yubikey one-time-password would be accepted for, but it might not be used?
	YubikeyOTPMaxAge int `json:"yubikey_otp_max_age,omitempty"`
}

func (p Privilege) targetType() string {
	return strings.Split(string(p.Level), "_")[0]
}

// String returns a string representation of the Privilege in English.
// Privileges are a little tricky to represent in English because the Privilege itself doesn't know if it exists on a user or if it has just been removed from a user, nor does it now anything about the target it's been granted on/revoked from other than a numerical ID. So we do the best we can.
func (p Privilege) String() string {
	switch p.targetType() {
	case "vm":
		return fmt.Sprintf("%s on VM #%d for %s")
	case "group":
		return fmt.Sprintf("%s on group #%d for %s")
	case "account":
		return fmt.Sprintf("%s on account #%d for %s")
	}
	return fmt.Sprintf("%s on the whole cluster for %s")
}

// PrettyPretty nicely formats the Privilege and sends it to the given writer.
// At the moment, the detail parameter is ignored.
func (p Privilege) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	_, err = wr.Write([]byte(p.String()))
	return
}

type Privileges []*Privilege

func (ps Privileges) IndexOf(priv Privilege) int {
	if priv.Username == "" || priv.Level == "" {
		return -1
	}
	for i, p := range ps {
		if p.VirtualMachineID == priv.VirtualMachineID &&
			p.GroupID == priv.GroupID && p.AccountID == priv.AccountID &&
			p.YubikeyRequired == priv.YubikeyRequired &&
			p.Level == priv.Level && p.Username == priv.Username {
			return i
		}
	}
	return -1
}
