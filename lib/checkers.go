package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// checkDiscPather makes sure the pather has all the necessary fields
// in it. This only matters for DiscName
func (c *bytemarkClient) checkDiscPather(discPather brain.DiscPather) (brain.DiscPather, error) {
	if discName, ok := discPather.(DiscName); ok {
		var newVMPather brain.VirtualMachinePather
		if newVMPather, err := c.checkVirtualMachinePather(discName.VirtualMachine); err != nil {
			return discName, err
		}
		discName.VirtualMachine = newVMPather
		return discName, nil
	}
	return discPather, nil
}

// checkVirtualMachinePather makes sure the pather has all the necessary fields
// in it. This only matters for VirtualMachineNames - IPs and IDs we'll just
// send to the brain no worries.
func (c *bytemarkClient) checkVirtualMachinePather(vmPather brain.VirtualMachinePather) (brain.VirtualMachinePather, error) {
	if vm, ok := vmPather.(VirtualMachineName); ok {
		if vm.Account == "" {
			acc, err := c.checkAccountPather(AccountName(vm.Account))
			vm.Account = string(acc.(AccountName))
			if err != nil {
				return vm, err
			}
		}
		if vm.Group == "" {
			vm.Group = DefaultGroup
		}

		if vm.VirtualMachine == "" {
			return vm, BadNameError{Type: "virtual machine", ProblemField: "name", ProblemValue: vm.VirtualMachine}
		}
		return vm, nil
	}
	return vmPather, nil
}

func (c *bytemarkClient) checkGroupPather(groupPather brain.GroupPather) (brain.GroupPather, error) {
	if group, ok := groupPather.(GroupName); ok {
		if group.Account == "" {
			if accountPather, err := c.checkAccountPather(AccountName(group.Account)); err != nil {
				if account, ok := accountPather.(AccountName); ok {
					group.Account = string(account)
				}
				return group, err
			}
		}
		if group.Group == "" {
			group.Group = DefaultGroup
		}
		return group, nil
	}
	return groupPather, nil
}

func (c *bytemarkClient) checkAccountPather(accountPather brain.AccountPather) (brain.AccountPather, error) {
	if account, ok := accountPather.(AccountName); ok {
		if account == "" && c.authSession != nil {
			log.Debug(log.LvlArgs, "CheckAccountPather called with empty name and a valid auth session - will try to figure out the default by talking to APIs.")
			if c.urls.Billing == "" {
				log.Debug(log.LvlArgs, "CheckAccountPather - there's no Billing endpoint, so we're most likely on a cluster devoid of bmbilling. Requesting account list from bigv...")
				brainAccs, err := c.getBrainAccounts()
				if err != nil {
					return account, err
				}
				log.Debugf(log.LvlArgs, "CheckAccountPather found %d accounts\r\n", len(brainAccs))
				if len(brainAccs) > 0 {
					log.Debugf(log.LvlArgs, "CheckAccountPather using the first account returned from bigv (%s) as the default\r\n", brainAccs[0].Name)
					account = AccountName(brainAccs[0].Name)
				}
			} else {
				log.Debug(log.LvlArgs, "CheckAccountPather finding the default billing account")
				billAcc, err := c.getDefaultBillingAccount()
				if err == nil && billAcc.IsValid() {
					log.Debugf(log.LvlArgs, "CheckAccountPather found the default billing account - %s\r\n", billAcc.Name)
					account = AccountName(billAcc.Name)
				} else if err != nil {
					return nil, err
				}
			}
		}
		if account == "" {
			return nil, NoDefaultAccountError{}
		}
		return AccountName(account), nil
	}
	return accountPather, nil
}
