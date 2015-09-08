package main

import (
	util "bigv.io/client/cmds/util"
	"bigv.io/client/util/log"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			log.Error("\r\nCaught an interrupt - exiting.\r\n")
			os.Exit(int(util.E_TRAPPED_INTERRUPT))
		}

	}()

	flags := util.MakeCommonFlagSet()

	flags.Parse(os.Args[1:])

	configDir := ""
	value := flags.Lookup("config-dir").Value
	if value != nil {
		configDir = value.String()
	}

	config, err := util.NewConfig(configDir, flags)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}

	dispatch, err := NewDispatcher(config)
	if err != nil {
		os.Exit(int(util.ProcessError(err)))
	}

	os.Exit(int(dispatch.Do(flags.Args())))
}
