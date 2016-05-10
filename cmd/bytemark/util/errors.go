package util

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"fmt"
)

type SubprocessFailedError struct {
	Args     []string
	ExitCode int
	Err      error
}

func (e *SubprocessFailedError) Error() string {
	return fmt.Sprintf("Running %s failed - %s", e.Args[0], e.Err.Error())
}

type NotEnoughArgumentsError struct {
}

func (e NotEnoughArgumentsError) Error() string {
	return "Not enough arguments passed to the command!"
}

// UsageDisplayedError is returned by commands when the user entered wrong info and the help was output
type UsageDisplayedError struct {
	TheProblem string
}

func (e UsageDisplayedError) Error() string {
	return e.TheProblem
}

type WontDeleteNonEmptyGroupError struct {
	Group *lib.GroupName
}

func (e WontDeleteNonEmptyGroupError) Error() string {
	return fmt.Sprintf("Group %s contains servers, will not be deleted without --recursive\r\n", e.Group)
}
