package main

import (
	"time"

	"github.com/urfave/cli"
)

type innerContext interface {
	Args() cli.Args
	Bool(name string) bool
	BoolT(name string) bool
	Duration(name string) time.Duration
	FlagNames() (names []string)
	Float64(name string) float64
	Generic(name string) interface{}
	GlobalBool(name string) bool
	GlobalBoolT(name string) bool
	GlobalDuration(name string) time.Duration
	GlobalFlagNames() (names []string)
	GlobalFloat64(name string) float64
	GlobalGeneric(name string) interface{}
	GlobalInt(name string) int
	GlobalInt64(name string) int64
	GlobalInt64Slice(name string) []int64
	GlobalIntSlice(name string) []int
	GlobalIsSet(name string) bool
	GlobalSet(name, value string) error
	GlobalString(name string) string
	GlobalStringSlice(name string) []string
	GlobalUint(name string) uint
	GlobalUint64(name string) uint64
	Int(name string) int
	Int64(name string) int64
	Int64Slice(name string) []int64
	IntSlice(name string) []int
	IsSet(name string) bool
	NArg() int
	NumFlags() int
	Parent() *cli.Context
	Set(name, value string) error
	String(name string) string
	StringSlice(name string) []string
	Uint(name string) uint
	Uint64(name string) uint64

	App() *cli.App
	Command() cli.Command
}

type cliContextWrapper struct {
	*cli.Context
}

func (c cliContextWrapper) App() *cli.App {
	return c.Context.App
}

func (c cliContextWrapper) Command() cli.Command {
	return c.Context.Command
}
