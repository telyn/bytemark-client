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
			os.Exit(-1)
		}

	}()

	flags := MakeCommonFlagSet()

	flags.Parse(os.Args[1:])
	configDirFlag := flags.Lookup("config-dir")

	configDir := ""
	if configDirFlag != nil {
		configDir = configDirFlag.Value.String()
	}
	config := NewConfig(configDir, flag.CommandLine)

	dispatch := NewDispatcher(config)

	dispatch.Do(flag.Args())
	os.Exit(0)
}
