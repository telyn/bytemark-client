package cmd

import (
	"fmt"
	"strings"
)

func (dispatch *Dispatcher) ShowAccount(args []string) {
	name := ParseAccountName(args[0])

	acc, err := dispatch.BigV.GetAccount(name)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Account %d: %s", acc.Id, acc.Name)

}
