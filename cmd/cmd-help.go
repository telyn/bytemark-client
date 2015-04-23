package cmd

import (
	"fmt"
	"strings"
)

func (cmds *CommandSet) HelpForHelp() {
	fmt.Println("bigv command-line client (the new, cool one)")
	fmt.Println()
	fmt.Println("There would be some usage output here if I had actually written any.")
	fmt.Println()
}

func (cmds *CommandSet) HelpForExitCodes() {
	fmt.Println(`bigv exit code list:

Exit code ranges:
    All these ranges are inclusive (i.e. 1-99 means numbers from 1 to 99, including 1 and 99.)

      0- 99: local problems
    100-199: problem talking to auth.
    200-299: problem talking to BigV.

    Errors in 100-299 with the same tens and units have the same meaning.

0 - 99 Errors:
    0
	Nothing went wrong and I feel great!

    3
	Couldn't read file from config directory

    4
	Couldn't write file to config directory

    5
	Trapped an interrupt signal, so exited.

100-299 Errors:

    100 / 200
        Unable to establish a connection to auth/BigV endpoint
    
    101 / 201
        Auth endpoint reported an internal error
    
    102 / 202
        Unable to parse output from auth endpoint (probably implies a protocol mismatch - try updating go-bigv)

    103
	Your credentials were rejected for containing invalid characters or fields.

    104
	Your credentials did not match any user on file - check you entered them correctly

    205
	Your user account doesn't have authorisation to perform that action

    206
        Something couldn't be found on BigV. This could be due to the following reasons:
            * It doesn't exist
	    * Your user account doesn't have authorisation to see it
	    * Protocol mismatch between the BigV endpoint and go-bigv.
`)
}

func (cmds *CommandSet) Help(args []string) {
	if len(args) == 0 {
		cmds.HelpForHelp()
		return
	}

	// please try and keep these in alphabetical order
	switch strings.ToLower(args[0]) {
	case "config":
		cmds.HelpForConfig()
	case "debug":
		cmds.HelpForDebug()
	case "exit":
		cmds.HelpForExitCodes()
	case "exit-codes":
		cmds.HelpForExitCodes()
	case "show":
		cmds.HelpForShow()
	}

}
