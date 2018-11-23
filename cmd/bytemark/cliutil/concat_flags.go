package cliutil

import (
	"github.com/urfave/cli"
)

// ConcatFlags combines multiple slices of flags together without modification
func ConcatFlags(flagSets ...[]cli.Flag) (res []cli.Flag) {
	// saves us re-allocating each time.. Probably doesn't save a lot of
	// memory/GC time in real life but eh. Feels good to be efficient.
	totalLen := 0
	for i := range flagSets {
		totalLen += len(flagSets[i])
	}
	res = make([]cli.Flag, 0, totalLen)

	for i := range flagSets {
		res = append(res, flagSets[i]...)
	}
	return
}
