package main

import (
	//	bigv "bigv.io/client/lib"
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

	config, err := NewConfig(configDir, flags)
	if err != nil {
		os.Exit(int(processError(err)))
	}

	dispatch, err := NewDispatcher(config)
	if err != nil {
		os.Exit(int(processError(err)))
	}

	os.Exit(int(dispatch.Do(flags.Args())))
}
