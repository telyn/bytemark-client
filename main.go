package main

import (
	"bigv.io/client/cmd"
	//	bigv "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

var (
	configDir  = flag.String("config", "", "Location of go-bigv's config store - defaults to ~/.go-bigv")
	help       = flag.Bool("help", false, "Display usage information")
	debugLevel = flag.Int("debug-level", 0, "How much debugging output to display - 0 is none, other values are 1 and 2.")
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			fmt.Printf("\r\nCaught an interrupt - exiting.\r\n")
			os.Exit(-1)
		}

	}()

	flag.Parse()
	config := cmd.NewConfig(*configDir, flag.CommandLine)

	dispatch := cmd.NewDispatcher(config)

	dispatch.Do(flag.Args())
	os.Exit(0)
}
