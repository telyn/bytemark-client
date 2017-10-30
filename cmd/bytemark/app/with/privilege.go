package with

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// normalisePrivilegeLevel makes sure the level provided is actually a valid PrivilegeLevel and provides a couple of aliases.
func normalisePrivilegeLevel(l brain.PrivilegeLevel) (level brain.PrivilegeLevel, ok bool) {
	level = brain.PrivilegeLevel(strings.ToLower(string(l)))
	switch level {
	case "cluster_admin", "account_admin", "group_admin", "vm_admin", "vm_console":
		ok = true
	case "server_admin", "server_console":
		level = brain.PrivilegeLevel(strings.Replace(string(level), "server", "vm", 1))
		ok = true
	case "console":
		level = "vm_console"
		ok = true
	}
	return
}

// Privilege gets the named PrivilegeFlag from the context, then resolves its target to an ID if needed to create a brain.Privilege, then attaches that to the context
func Privilege(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		pf := c.PrivilegeFlag(flagName)
		level, ok := normalisePrivilegeLevel(pf.Level)
		if !ok && !c.Bool("force") {
			return fmt.Errorf("Unexpected privilege level '%s' - expecting account_admin, group_admin, vm_admin or vm_console", pf.Level)
		}
		c.Privilege = brain.Privilege{
			Username: pf.Username,
			Level:    level,
		}
		err = preflight(c)
		if err != nil {
			return
		}
		switch c.Privilege.TargetType() {
		case brain.PrivilegeTargetTypeVM:
			var vm brain.VirtualMachine
			vm, err = c.Client().GetVirtualMachine(*pf.VirtualMachineName)
			c.Privilege.VirtualMachineID = vm.ID
		case brain.PrivilegeTargetTypeGroup:
			var group brain.Group
			group, err = c.Client().GetGroup(*pf.GroupName)
			c.Privilege.GroupID = group.ID
		case brain.PrivilegeTargetTypeAccount:
			var acc lib.Account
			acc, err = c.Client().GetAccount(pf.AccountName)
			if err != nil {
				if _, ok := err.(lib.BillingAccountNotFound); ok {
					err = nil
				}
			}
			c.Privilege.AccountID = acc.BrainID
		}
		return
	}
}

// User gets a username from the given flag, then gets the corresponding User from the brain, and attaches it to the app.Context.
func User(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		if c.User != nil {
			return
		}
		if err = preflight(c); err != nil {
			return
		}
		username := c.String(flagName)
		if username == "" {
			username = c.Client().GetSessionUser()
		}

		user, err := c.Client().GetUser(username)
		if err != nil {
			return
		}
		c.User = &user
		return
	}
}

// VirtualMachine gets a VirtualMachineName from a flag, then gets the named VirtualMachine from the brain and attaches it to the app.Context.
func VirtualMachine(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		if c.VirtualMachine != nil {
			return
		}
		err = preflight(c)
		if err != nil {
			return
		}
		vmName := c.VirtualMachineName(flagName)
		vm, err := c.Client().GetVirtualMachine(vmName)
		if err != nil {
			return
		}
		c.VirtualMachine = &vm
		return
	}
}
