package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
	"strings"
)

const privilegeText = `PRIVILEGES

Privileges are listed in descending order of privilege. Users can only create or modify privileges which have a lower level than themselves. For example, an account_admin can only create or modify group_admin and lower privileges. This means that to grant or revoke another use account_admin on an account you have account_admin on, you must contact support.

cluster_admin - full administration rights to the whole cluster
account_admin - full administration rights to a whole account
group_admin - create, modify & delete servers within a group
vm_admin - modify & administer a single server
vm_console - access to a server's console`

func init() {
	adminCommands = append(adminCommands, cli.Command{
		Name:        "grant",
		Usage:       "grant privileges on bytemark self-service objects to other users",
		UsageText:   "bytemark grant <privilege> [on] <object> [from|to] <user>\r\nbytemark grant cluster_admin [to] <user>",
		Description: "Grant a privilege to a user for a particular bytemark object\r\n\r\n" + privilegeText,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "yubikey-required",
				Usage: "Set if the privilege should require a yubikey.",
			},
		},
		Action: With(AuthProvider, func(c *Context) (err error) {
			priv, _, err := parsePrivilege(c)
			if err != nil {
				return
			}
			err = global.Client.GrantPrivilege(priv)
			if err == nil {
				log.Outputf("Granted %s to %s\r\n", priv.Level, priv.Username)
			}
			return
		}),
	})
}

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

// fillPrivilegeTarget adds the object to the privilege, trying to use it as a VM, Group or Account name depending on what PrivilegeLevel the Privilege is for. The target is expected to be the NextArg at this point in the Context
func fillPrivilegeTarget(c *Context, p *brain.Privilege) (targetName string, err error) {
	name, err := c.NextArg()
	if err != nil {
		return
	}
	if strings.HasPrefix(string(p.Level), "vm") {
		vmName, err := lib.ParseVirtualMachineName(name, global.Config.GetVirtualMachine())
		if err != nil {
			return "", err
		}
		vm, err := global.Client.GetVirtualMachine(vmName)
		if err != nil {
			return "", err
		}
		targetName = vm.Hostname
		p.VirtualMachineID = vm.ID
	} else if strings.HasPrefix(string(p.Level), "group") {
		groupName := global.Client.ParseGroupName(name, global.Config.GetGroup())
		group, err := global.Client.GetGroup(groupName)
		if err != nil {
			return "", err
		}
		targetName = groupName.String()
		p.GroupID = group.ID
	} else if strings.HasPrefix(string(p.Level), "account") {
		accountName := global.Client.ParseAccountName(name, global.Config.GetIgnoreErr("account"))
		account, err := global.Client.GetAccount(accountName)
		if err != nil {
			return "", err
		}
		targetName = account.Name
		p.AccountID = account.BrainID
	}
	return
}

// creates a Privilege from the arguments in the Context.
func parsePrivilege(c *Context) (p brain.Privilege, targetName string, err error) {
	p.YubikeyRequired = c.Bool("yubikey-required")
	var level string
	level, err = c.NextArg()
	p.Level = brain.PrivilegeLevel(level)
	if err != nil {
		return
	}
	var ok bool
	p.Level, ok = normalisePrivilegeLevel(p.Level)
	if !ok && !c.Bool("force") {
		err = fmt.Errorf("Invalid privilege level %s", p.Level)
		return
	}
	if c.Args()[0] == "on" {
		_, _ = c.NextArg()
	}
	targetName, err = fillPrivilegeTarget(c, &p)
	if err != nil {
		return
	}
	if c.Args()[0] == "to" || c.Args()[0] == "from" {
		_, _ = c.NextArg()
	}
	p.Username, err = c.NextArg()
	return
}
