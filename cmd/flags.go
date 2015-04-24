package main

import "flag"

func MakeCommonFlagSet() *flag.FlagSet {
	flags := flag.NewFlagSet("go-bigv", flag.ExitOnError)

	// because I'm creating my own help functions I'm not going to supply usage info. Neener neener.
	flags.Bool("help", false, "")
	flags.Bool("force", false, "")
	flags.Int("debug-level", 0, "")
	flags.String("user", "", "")
	flags.String("account", "", "")
	flags.String("endpoint", "", "")
	flags.String("auth-endpoint", "", "")
	flags.String("config-dir", "", "")

	return flags
}
