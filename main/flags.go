package main

import "flag"

// MakeCommonFlagSet creates a FlagSet which provides the flags shared between the main command and the sub-commands.
func MakeCommonFlagSet() *flag.FlagSet {
	flags := flag.NewFlagSet("go-bigv", flag.ContinueOnError)

	// because I'm creating my own help functions I'm not going to supply usage info. Neener neener.
	flags.Bool("help", false, "")
	flags.Bool("force", false, "")
	flags.Bool("silent", false, "")
	flags.Int("debug-level", 0, "")
	flags.String("user", "", "")
	flags.String("account", "", "")
	flags.String("endpoint", "", "")
	flags.String("auth-endpoint", "", "")
	flags.String("config-dir", "", "")

	return flags
}
