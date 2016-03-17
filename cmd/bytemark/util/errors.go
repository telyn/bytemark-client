package util

import (
	"bytemark.co.uk/client/lib"
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

type WontDeleteNonEmptyGroupError struct {
	Group *lib.GroupName
}

func (e WontDeleteNonEmptyGroupError) Error() string {
	return fmt.Sprintf("Group %s contains servers, will not be deleted without --recursive\r\n", e.Group)
}
