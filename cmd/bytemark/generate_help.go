package main

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/urfave/cli"
)

func generateHelp([]cli.Command) {
	for idx, cmd := range commands {
		switch cmd.Name {
		case "admin":
			commands[idx].Description = cmd.Description + app.GenerateCommandsHelp(admin.Commands)
		case "commands":
			commands[idx].Description = cmd.Description + app.GenerateCommandsHelp(commands)
		}
	}
}
