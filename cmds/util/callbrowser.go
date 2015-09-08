// +build !windows

package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func openCommand() string {
	switch runtime.GOOS {
	case "darwin":
		return "open"
	default:
		return "xdg-open"
	}
}

func CallBrowser(url string) error {
	command := openCommand()
	bin, err := exec.LookPath(command)
	if err != nil {

		command = "/usr/bin/x-www-browser"
		bin, err = exec.LookPath(command)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stderr, "Running a browser to open %s...\r\n", url)

	var attr os.ProcAttr
	proc, err := os.StartProcess(bin, []string{bin, url}, &attr)
	subprocErr := SubprocessFailedError{Args: []string{bin, url}}
	if err != nil {
		subprocErr.Err = err
		return &subprocErr
	}
	state, err := proc.Wait()

	subprocErr.Err = err
	waitStatus, ok := state.Sys().(syscall.WaitStatus)
	if ok {
		subprocErr.ExitCode = waitStatus.ExitStatus()
	}
	return &subprocErr

}
