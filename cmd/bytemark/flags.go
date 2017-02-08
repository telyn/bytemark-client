package main

type AccountNameFlag string

func (name *AccountNameFlag) Set(value string) error {
	*name = AccountNameFlag(global.Client.ParseAccountName(value, global.Config.GetIgnoreErr("account")))
	return nil
}

func (name *AccountNameFlag) String() string {
	return string(*name)
}
