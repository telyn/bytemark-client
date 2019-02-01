package util

import (
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// SubprocessFailedError is returned when a process run by bytemark-client (e.g. open/xdg-open to open a browser) failed
type SubprocessFailedError struct {
	Args     []string
	ExitCode int
	Err      error
}

func (e SubprocessFailedError) Error() string {
	return fmt.Sprintf("Running %s failed - %s", e.Args[0], e.Err.Error())
}

// UsageDisplayedError is returned by commands when the user entered wrong info and the help was output
type UsageDisplayedError struct {
	TheProblem string
	Command    string
}

func (e UsageDisplayedError) Error() string {
	return e.TheProblem + "\r\n\r\nFor more information, see `bytemark help " + e.Command + "`"
}

// WontDeleteGroupWithVMsError is returned when 'delete group' was called on a group with stuff in, without --recursive being specified
type WontDeleteGroupWithVMsError struct {
	Group pathers.GroupName
}

func (e WontDeleteGroupWithVMsError) Error() string {
	return fmt.Sprintf("Group %s contains servers, will not be deleted without --recursive\r\n", e.Group)
}

// RecursiveDeleteGroupError is returned by delete group when called with --recursive, when deleting VMs.
type RecursiveDeleteGroupError struct {
	Group pathers.GroupName
	// Map of VirtualMachine names to the error that occurred when trying to delete them.
	// N.B. that this will not contain nil errors for VMs that were successfully deleted.
	Errors map[string]error
}

func (e RecursiveDeleteGroupError) Error() string {
	strs := make([]string, 0, len(e.Errors))
	for vm, err := range e.Errors {
		strs = append(strs, fmt.Sprintf("%s: %s", vm, err))
	}

	return fmt.Sprintf("Errors occurred while deleting VMs in group %s: \n\t%s", e.Group, strings.Join(strs, "\n\t"))
}
