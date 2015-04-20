package cmd

import (
	"fmt"
	"strings"
)

func (cmds *CommandSet) HelpForConfig() {
	fmt.Println("bigv config set")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    bigv config")
	fmt.Println("        Outputs the current values of all variables and what source they were derived from")
	fmt.Println()
	fmt.Println("    bigv config set <variable> <value>")
	fmt.Println("        Sets a variable by writing to your bigv config (usually ~/.go-bigv)")
	fmt.Println()
	fmt.Println("    bigv config unset <variable>")
	fmt.Println("        Unsets a variable by removing data from bigv config (usually ~/.go-bigv)")
	fmt.Println()
	fmt.Println("Available variables:")
	fmt.Println("    endpoint - the BigV endpoint to connect to. https://uk0.bigv.io is the default")
	fmt.Println("    auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.")
	fmt.Println("    debug-level - the default debug level. Set to 0 unless you like lots of output")
	fmt.Println()
}

// Set provides the bigv set command, which sets variables in the user's config
// It's slightly more user friendly than echo "value" > ~/.go-bigv/
func (cmds *CommandSet) Config(args []string) {
	if len(args) == 0 {
		for _, v := range cmds.config.GetAll() {
			fmt.Println("%s: '%s' (%s)", v.Name, v.Value, v.Source)
		}
		return
	} else if len(args) == 1 {
		cmds.HelpForConfig()
		return
	}

	switch strings.ToLower(args[0]) {
	case "set":
		variable := strings.ToLower(args[1])

		oldVar := cmds.config.GetV(variable)

		if len(args) == 2 {
			fmt.Printf("%s: '%s' (%s)", oldVar.Name, oldVar.Value, oldVar.Source)
		}

		// TODO(telyn): input validation ha ha ha
		cmds.config.SetPersistent(variable, args[2], "CMD set")

		if oldVar.Source == "config" {
			fmt.Printf("%s has been changed.\r\nOld value: %s\r\nNew value: %s\r\n", variable, oldVar.Value, args[1])
		} else {
			fmt.Printf("%s has been set. \r\nNew value: %s\r\n", variable, args[1])
		}

	case "unset":

	}

}
