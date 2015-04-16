package cmd

import (
	"fmt"
	"strings"
)

func (dispatch *Dispatcher) HelpForSet() {
	fmt.Println("bigv set")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("    bigv set <variable> <value>")
	fmt.Println("")
	fmt.Println("Sets a variable by writing to your bigv config (usually ~/.go-bigv)")
	fmt.Println("")
	fmt.Println("Available variables:")
	fmt.Println("    endpoint - the BigV endpoint to connect to. https://uk0.bigv.io is the default")
	fmt.Println("    auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.")
	fmt.Println("    debug-level - the default debug level. Set to 0 unless you like lots of output")
	fmt.Println("")
	fmt.Println("")
}

func (dispatch *Dispatcher) HelpForUnset() {
	fmt.Println("bigv unset")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("    bigv unset <variable>")
	fmt.Println("")
	fmt.Println("Unsets a variable by removing data from bigv config (usually ~/.go-bigv)")
	fmt.Println("See the set command for the list of available variables")
	fmt.Println("")
}

// Set provides the bigv set command, which sets variables in the user's config
// It's slightly more user friendly than echo "value" > ~/.go-bigv/
func (dispatch *Dispatcher) Set(args []string) {
	if len(args) != 2 {
		dispatch.HelpForSet()
		return
	}

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
	if len(args) != 2 {
		dispatch.HelpForUnset()
		return
	}

	variable := strings.ToLower(args[0])

	oldVar := dispatch.Config.Get(variable)

	dispatch.Config.Unset(variable)
	fmt.Printf("%s has been unset.\r\nOld value: %s\r\n", variable, oldVar)

}
