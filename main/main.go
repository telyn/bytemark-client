package main

import (
	//	bigv "bigv.io/client/lib"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			fmt.Printf("\r\nCaught an interrupt - exiting.\r\n")
			os.Exit(int(E_TRAPPED_INTERRUPT))
		}

	}()

	flags := MakeCommonFlagSet()

	flags.Parse(os.Args[1:])
	configDir := ""
	value := flags.Lookup("config-dir").Value
	if value != nil {
		configDir = value.String()
	}
	config := NewConfig(configDir, flag.CommandLine)

	dispatch := NewDispatcher(config)

	dispatch.Do(flags.Args())
	os.Exit(0)
}
