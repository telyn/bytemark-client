package cmds

import (
	"bufio"
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/util/log"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
)

// HelpForDebug outputs usage information for the debug command.
func (commands *CommandSet) HelpForDebug() util.ExitCode {
	log.Log("bytemark debug")
	log.Log()
	log.Log("Usage:")
	log.Log("    bytemark debug [--junk-token] [--auth] [--use-billing] GET <path>")
	log.Log("    bytemark debug [--junk-token] [--auth] [--use-billing] DELETE <path>")
	log.Log("    bytemark debug [--junk-token] [--auth] [--use-billing] PUT <path>")
	log.Log("    bytemark debug [--junk-token] [--auth] [--use-billing] POST <path>")
	log.Log()
	log.Log("GET sends an HTTP GET request with an optional valid authorization header to the given path on the API endpoint and pretty-prints the received json.")
	log.Log("The rest do similar, but PUT and POST")
	log.Log("The --junk-token flag sets the token to empty - useful if you want to ensure that credential-auth is working, or you want to do something as another user")
	log.Log("The --auth token tells the client to gain valid auth and send the auth header on that request.")
	log.Log()
	return util.E_USAGE_DISPLAYED

}

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>
func (commands *CommandSet) Debug(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	junkToken := flags.Bool("junk-token", false, "")
	shouldAuth := flags.Bool("auth", false, "")
	billing := flags.Bool("use-billing", false, "")
	flags.Parse(args)
	args = commands.config.ImportFlags(flags)

	if *junkToken {
		commands.config.Set("token", "", "FLAG junk-token")
	}

	if len(args) < 1 {
		return commands.HelpForDebug()
	}

	switch args[0] {
	case "GET", "PUT", "POST", "DELETE":
		if !strings.HasPrefix(args[1], "/") {
			args[1] = "/" + args[1]
		}
		if *shouldAuth {
			err := commands.EnsureAuth()
			if err != nil {
				return util.ProcessError(err)
			}
		}

		requestBody := ""
		err := error(nil)
		if args[0] == "PUT" || args[0] == "POST" {
			buf := bufio.NewReader(os.Stdin)
			requestBody, err = buf.ReadString(byte(uint8(14)))
			if err != nil && err != io.EOF {
				return util.ProcessError(err)
			}
		}
		body, err := commands.client.RequestAndRead(*shouldAuth, args[0], args[1], requestBody)
		if err != nil {
			return util.ProcessError(err)
		}

		buf := new(bytes.Buffer)
		json.Indent(buf, body, "", "    ")
		log.Log(buf.String())
	case "config":
		vars, err := commands.config.GetAll()
		if err != nil {
			return util.ProcessError(err)
		}
		indented, _ := json.MarshalIndent(vars, "", "    ")
		log.Log(string(indented))
	default:
		commands.HelpForDebug()
	}
	return util.E_SUCCESS
}
