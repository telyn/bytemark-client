package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// AccountNameFlag is used for all --account flags, including the global one.
type AccountNameFlag string

// Set runs lib.Client.ParseAccountName using the global.Client to make sure we get just the 'pure' account name; no cluster / endpoint details
func (name *AccountNameFlag) Set(value string) error {
	*name = AccountNameFlag(global.Client.ParseAccountName(value, global.Config.GetIgnoreErr("account")))
	return nil
}

// String returns the AccountNameFlag as a string.
func (name *AccountNameFlag) String() string {
	return string(*name)
}

// GroupNameFlag is used for all --account flags, including the global one.
type GroupNameFlag lib.GroupName

// Set runs lib.Client.ParseGroupName using the global.Client to make sure we have a valid group name
func (name *GroupNameFlag) Set(value string) error {
	gp := global.Client.ParseGroupName(value, global.Config.GetGroup())
	*name = GroupNameFlag(*gp)
	return nil
}

// String returns the GroupNameFlag as a string.
func (name *GroupNameFlag) String() string {
	return lib.GroupName(*name).String()
}
