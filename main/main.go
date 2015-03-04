package main

import (
	"bigv.io/client/cmd"
	//	bigv "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
)

var (
	configDir = flag.String("config", "", "Location of go-bigv's config store - defaults to ~/.go-bigv")
)

func main() {
	flag.Parse()
	config := cmd.NewConfig(*configDir, flag.CommandLine)

	// this line is just to make it build, will be removed
	fmt.Printf("Using configuration in %s\n", config.Dir)
	os.Exit(0)
}
