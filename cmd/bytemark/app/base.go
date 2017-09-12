package app

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/urfave/cli"
)

func BaseAppSetup(flags []cli.Flag, config util.ConfigManager, commands []cli.Command) (app *cli.App, err error) {
	app = cli.NewApp()
	app.Version = lib.Version
	app.Flags = flags
	app.Commands = commands

	/* TODO(telyn): move this clump over to main, probably.
	   Phil - if you see this, I've messed up
	// last minute alterations to commands
	// used for modifying help descriptions, mostly.
	for idx, cmd := range app.Commands {
		switch cmd.Name {
		case "admin":
			app.Commands[idx].Description = cmd.Description + "\r\n\r\n" + GenerateCommandsHelp(adminCommands)
		case "commands":
			app.Commands[idx].Description = cmd.Description + "\r\n\r\n" + GenerateCommandsHelp(app.Commands)
		}
	}
	*/
	app.Commands = cliutil.CreateMultiwordCommands(app.Commands)
	return

}
