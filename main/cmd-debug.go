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
func (commands *CommandSet) HelpForDebug() {
	fmt.Println("bigv debug")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    bigv debug [--junk-token] [--auth] GET <path>")
	fmt.Println("    bigv debug [--junk-token] [--auth] DELETE <path>")
	fmt.Println("    bigv debug [--junk-token] [--auth] PUT <path>")
	fmt.Println("    bigv debug [--junk-token] [--auth] POST <path>")
	fmt.Println()
	fmt.Println("GET sends an HTTP GET request with an optional valid authorization header to the given path on the BigV endpoint and pretty-prints the received json.")
	fmt.Println("The rest do similar, but PUT and POST")
	fmt.Println("The --junk-token flag sets the token to empty - useful if you want to ensure that credential-auth is working, or you want to do something as another user")
	fmt.Println("The --auth token tells the client to gain valid auth and send the auth header on that request.")
	fmt.Println()
}

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>
func (commands *CommandSet) Debug(args []string) {
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
		return
	}

	switch args[0] {
	case "GET", "PUT", "POST", "DELETE":
		if !strings.HasPrefix(args[1], "/") {
			args[1] = "/" + args[1]
		}
		if *shouldAuth {
			commands.EnsureAuth()
		}

		requestBody := ""
		err := error(nil)
		if args[0] == "PUT" || args[0] == "POST" {
			buf := bufio.NewReader(os.Stdin)
			requestBody, err = buf.ReadString(byte(uint8(14)))
			if err != nil {
				Exit(err)
			}
		}
		body, err := commands.bigv.RequestAndRead(*shouldAuth, args[0], args[1], requestBody)
		if err != nil {
			exit(err)
		}

		buf := new(bytes.Buffer)
		json.Indent(buf, body, "", "    ")
		fmt.Printf("%s", buf)
	case "config":
		indented, _ := json.MarshalIndent(commands.config.GetAll(), "", "    ")
		fmt.Printf("%s", indented)
	default:
		commands.HelpForDebug()
	}
}
