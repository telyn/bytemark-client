package app

import (
	"io"
	"os"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

// BaseAppSetup sets up a cli.App for the given commands and config
func BaseAppSetup(flags []cli.Flag, commands []cli.Command) (app *cli.App, err error) {
	app = cli.NewApp()
	app.Version = lib.Version
	app.Flags = flags
	app.Commands = commands
	app.Writer = io.MultiWriter(
		log.LogFile,
		os.Stdout,
	)
	app.ErrWriter = io.MultiWriter(
		log.LogFile,
		os.Stderr,
	)

	for idx, cmd := range app.Commands {
		switch cmd.Name {
		case "admin":
			app.Commands[idx].Description = cmd.Description + "\r\n\r\n" + GenerateCommandsHelp(admin.Commands)
		case "commands":
			app.Commands[idx].Description = cmd.Description + "\r\n\r\n" + GenerateCommandsHelp(app.Commands)
		}
	}
	app.Commands = cliutil.CreateMultiwordCommands(app.Commands)
	return

}

// SetClientAndConfig adds the client and config to the given app.
// it abstracts away setting the Metadata on the app. Mostly so that we get some type-checking.
// without it - it's just assigning to an interface{} which will always succeed,
// and which would near-inevitably result in hard-to-debug null pointer errors down the line.
func SetClientAndConfig(app *cli.App, client lib.Client, config util.ConfigManager) {
	if app.Metadata == nil {
		app.Metadata = make(map[string]interface{})
	}
	app.Metadata["client"] = client
	app.Metadata["config"] = config
}
