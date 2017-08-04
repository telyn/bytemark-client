package main

import (
	"testing"

	"github.com/cheekybits/is"
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
		collected := collectArgs(argstr)

		is.Equal(len(argsli), len(collected))
		for i := range collected {
			is.Equal(argsli[i], collected[i])
		}
	}
}
