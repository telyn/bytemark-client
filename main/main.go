package main

import (
	"bigv.io/client/cmd"
	bigv "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
)

var (
	configDir = flag.String("config", "", "Location of go-bigv's config store - defaults to ~/.go-bigv")
)

func main() {
	flag.Parse()
	config := cmd.NewConfig(configDir)

	os.Exit(0)
}
