package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (dispatch *Dispatcher) Debug(args []string) {
	dispatch.BigV.DebugLevel = 1

	body, err := dispatch.BigV.Request(args[0], args[1], "")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	json.Indent(buf, body, "", "    ")
	fmt.Printf("%s", buf)
}
