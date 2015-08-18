package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// HelpForDebug outputs usage information for the debug command.
func (commands *CommandSet) HelpForDebug() ExitCode {
	fmt.Println("go-bigv debug")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    go-bigv debug [--junk-token] [--auth] GET <path>")
	fmt.Println("    go-bigv debug [--junk-token] [--auth] DELETE <path>")
	fmt.Println("    go-bigv debug [--junk-token] [--auth] PUT <path>")
	fmt.Println("    go-bigv debug [--junk-token] [--auth] POST <path>")
	fmt.Println()
	fmt.Println("GET sends an HTTP GET request with an optional valid authorization header to the given path on the BigV endpoint and pretty-prints the received json.")
	fmt.Println("The rest do similar, but PUT and POST")
	fmt.Println("The --junk-token flag sets the token to empty - useful if you want to ensure that credential-auth is working, or you want to do something as another user")
	fmt.Println("The --auth token tells the client to gain valid auth and send the auth header on that request.")
	fmt.Println()
	return E_USAGE_DISPLAYED

}

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>
func (commands *CommandSet) Debug(args []string) ExitCode {
	flags := MakeCommonFlagSet()
	junkToken := flags.Bool("junk-token", false, "")
	shouldAuth := flags.Bool("auth", false, "")
	flags.Parse(args)
	args = commands.config.ImportFlags(flags)

	if *junkToken {
		commands.config.Set("token", "", "FLAG junk-token")
	}

	if len(args) < 1 {
		commands.HelpForDebug()
		return E_SUCCESS
	}

	switch args[0] {
	case "GET", "PUT", "POST", "DELETE":
		// BUG(telyn): don't panic
		if !strings.HasPrefix(args[1], "/") {
			args[1] = "/" + args[1]
		}
		if *shouldAuth {
			err := commands.EnsureAuth()
			if err != nil {
				return processError(err)
			}
		}

		requestBody := ""
		err := error(nil)
		if args[0] == "PUT" || args[0] == "POST" {
			buf := bufio.NewReader(os.Stdin)
			requestBody, err = buf.ReadString(byte(uint8(14)))
			if err != nil {
				// BUG(telyn): deal with EOFs properly
				return processError(err)
			}
		}
		body, err := commands.bigv.RequestAndRead(*shouldAuth, args[0], args[1], requestBody)
		if err != nil {
			return processError(err)
		}

		buf := new(bytes.Buffer)
		json.Indent(buf, body, "", "    ")
		fmt.Printf("%s", buf)
	case "config":
		vars, err := commands.config.GetAll()
		if err != nil {
			return processError(err)
		}
		indented, _ := json.MarshalIndent(vars, "", "    ")
		fmt.Printf("%s", indented)
	default:
		commands.HelpForDebug()
	}
	return E_SUCCESS
}
