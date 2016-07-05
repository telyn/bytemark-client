package main

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestCollectArgs(t *testing.T) {
	is := is.New(t)
	tests := map[string][]string{
		"":    {},
		"-tt": {"-tt"},
		"-i 'muh_identities.rsa'":   {"-i", "muh_identities.rsa"},
		"-i 'a long path/id_rsa'":   {"-i", "a long path/id_rsa"},
		"-i \"a long path/id_rsa\"": {"-i", "a long path/id_rsa"},
	}

	for argstr, argsli := range tests {
		is.Equal(argsli, collect_args(argstr))
	}
}
