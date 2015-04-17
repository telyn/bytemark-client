package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (commands *CommandSet) HelpForDebug() {
	fmt.Println("bigv debug")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    bigv debug GET <path>")
	fmt.Println()
	fmt.Println("GET sends an HTTP GET request with a valid authorization header to the given path on the BigV endpoint and pretty-prints the received json.")
	fmt.Println()
}

// TODO(telyn): does the URL really have to start with /?

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>
// URL probably needs to start with a /
func (commands *CommandSet) Debug(args []string) {
	if len(args) < 2 {
		commands.HelpForDebug()
		return
	}
	// TODO(telyn): add a flag to disable auth
	// TODO(telyn): add a flag to junk the token
	shouldAuth := true
	commands.EnsureAuth()
	commands.bigv.SetDebugLevel(1)

	// make sure the command is well-formed

	body, err := commands.bigv.RequestAndRead(shouldAuth, args[0], args[1], "")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	json.Indent(buf, body, "", "    ")
	fmt.Printf("%s", buf)
}
