package main

import (
	"bufio"
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"bytes"
	"encoding/json"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "debug",
		Usage:     "Test out the Bytemark API",
		UsageText: "bytemark debug [--junk-token] [--auth] [--use-billing] GET|POST|PUT|DELETE <path>",
		Description: `GET sends an HTTP GET request with an optional valid authorization header to the given path on the API endpoint and pretty-prints the received json.
The rest do similar, but PUT and POST
The --junk-token flag sets the token to empty - useful if you want to ensure that credential-auth is working, or you want to do something as another user
The --auth flag tells the client to gain valid auth and send the auth header on that request.
The --use-billing flag tells the client to send the request to the billing endpoint instead of the brain.`,
		Action: func(c *cli.Context) {
			flags := util.MakeCommonFlagSet()
			junkToken := flags.Bool("junk-token", false, "")
			shouldAuth := flags.Bool("auth", false, "")
			billing := flags.Bool("use-billing", false, "")
			flags.Parse(c.Args())
			args := global.Config.ImportFlags(flags)

			endpoint := lib.EP_BRAIN
			if *billing {
				endpoint = lib.EP_BILLING
			}

			if *junkToken {
				global.Config.Set("token", "", "FLAG junk-token")
			}

			if len(args) < 1 {
				global.Error = &util.PEBKACError{}
				return
			}

			switch args[0] {
			case "GET", "PUT", "POST", "DELETE":
				method := args[0]
				if len(args) < 2 {
					global.Error = &util.PEBKACError{}
					return
				}
				url := args[1]
				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				if *shouldAuth {
					err := EnsureAuth()
					if err != nil {
						global.Error = err
						return
					}
				}

				err := error(nil)
				reader := io.Reader(nil)
				if method == "PUT" || method == "POST" {
					reader = bufio.NewReader(os.Stdin)
					// read until an eof
				}
				req, err := global.Client.BuildRequest(method, endpoint, url)
				if !*shouldAuth {
					req, err = global.Client.BuildRequestNoAuth(method, endpoint, url)
				}
				if err != nil {
					global.Error = err
					return
				}

				statusCode, body, err := req.Run(reader, nil)
				if err != nil {
					global.Error = err
					return
				}
				reqURL := req.GetURL()
				log.Logf("%s %s: %d\r\n", method, reqURL.String(), statusCode)

				buf := new(bytes.Buffer)
				json.Indent(buf, body, "", "    ")
				log.Log(buf.String())
			default:
				return
			}
		},
	})

}

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>
