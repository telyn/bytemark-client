package main

import (
	"bigv.io/client/cmd"
	//	bigv "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
)

var (
	configDir  = flag.String("config", "", "Location of go-bigv's config store - defaults to ~/.go-bigv")
	debugLevel = flag.Int("debug-level", 0, "How much debugging output to display - 0 is none, other values are 1 and 2.")
)

func main() {
	flag.Parse()
	config := cmd.NewConfig(*configDir, flag.CommandLine)

	fmt.Fprintf(os.Stderr, "Using config in %s \r\n", config.Dir)

	dispatch := cmd.NewDispatcher(config)

	dispatch.Do(flag.Args())
	os.Exit(0)
}
