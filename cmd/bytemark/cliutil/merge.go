package cliutil

import "github.com/urfave/cli"

// MergeCommand merges src into dst, only copying non-nil fields of src,
// and calling mergeCommands upon the .Subcommands
// and appending all .Flags
func MergeCommand(dst *cli.Command, src cli.Command) {
	if src.Usage != "" {
		dst.Usage = src.Usage
	}
	if src.UsageText != "" {
		dst.UsageText = src.UsageText
	}
	if src.Description != "" {
		dst.Description = src.Description
	}
	if src.Action != nil {
		dst.Action = src.Action
	}
	if src.Flags != nil {
		dst.Flags = append(dst.Flags, src.Flags...)
	}
	if src.Subcommands != nil {
		dst.Subcommands = MergeCommands(dst.Subcommands, src.Subcommands)
	}
}

// MergeCommands copies over all the commands from base to result,
// then puts all the commands from extras in too, overwriting any provided fields.
func MergeCommands(base []cli.Command, extras []cli.Command) (result []cli.Command) {
	result = make([]cli.Command, len(base))
	copy(result, base)

	for _, cmd := range extras {
		found := false
		for idx := range result {
			if result[idx].Name == cmd.Name {
				MergeCommand(&result[idx], cmd)
				found = true
			}
		}
		if !found {
			result = append(result, cmd)
		}
	}
	return
}
