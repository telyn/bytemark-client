package util

import (
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
