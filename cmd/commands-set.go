package cmd

import (
	"fmt"
	"strings"
)

// Set provides the bigv set command, which sets variables in the user's config
// It's slightly more user friendly than echo "value" > ~/.go-bigv/
func (dispatch *Dispatcher) Set(args []string) {
	variable := strings.ToLower(args[0])

	oldVar := dispatch.Config.Get(variable)

	// TODO(telyn): input validation ha ha ha
	dispatch.Config.SetPersistent(variable, args[1])

	if oldVar != "" {
		fmt.Printf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", variable, oldVar, args[1])
	} else {
		fmt.Printf("%s has been set. \r\nNew value: %s\r\n", variable, args[1])
	}

}

func (dispatch *Dispatcher) Unset(args []string) {
	variable := strings.ToLower(args[0])

	oldVar := dispatch.Config.Get(variable)

	dispatch.Config.Unset(variable)
	fmt.Printf("%s has been unset.\r\nOld value: %s\r\n", variable, oldVar)

}
