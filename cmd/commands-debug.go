package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TODO(telyn): does the URL really have to start with /?

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>
// URL probably needs to start with a /
func (dispatch *Dispatcher) Debug(args []string) {
	dispatch.BigV.DebugLevel = 1

	// make sure the command is well-formed

	body, err := dispatch.BigV.RequestAndRead(args[0], args[1], "")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	json.Indent(buf, body, "", "    ")
	fmt.Printf("%s", buf)
}
