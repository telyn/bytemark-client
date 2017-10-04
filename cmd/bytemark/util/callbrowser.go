// +build !windows

package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

func openCommand() string {
	switch runtime.GOOS {
	case "darwin":
		return "open"
	default:
		return "xdg-open"
	}
}

// CallBrowser opens the user's desktop browser to the given URL.
// It tries really hard - first trying open on mac or xdg-open on other systems.
// If xdg-open couldn't be used, it attempts to use /usr/bin/x-www-browser
func CallBrowser(url string) error {
	command := openCommand()
	bin, err := exec.LookPath(command)
	if err != nil {
		if runtime.GOOS != "darwin" {

			command = "/usr/bin/x-www-browser"
			bin, err = exec.LookPath(command)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var attr os.ProcAttr

	log.Debugf(log.LvlOutline, "Executing %s \"%s\"", bin, url)

	proc, err := os.StartProcess(bin, []string{bin, url}, &attr)
	if err != nil {
		return SubprocessFailedError{Args: []string{bin, url}, Err: err}
	}
	state, err := proc.Wait()
	if err != nil {
		return SubprocessFailedError{
			Args: []string{bin, url},
			Err:  err,
		}
	}

	waitStatus, ok := state.Sys().(syscall.WaitStatus)
	if ok {
		exitCode := waitStatus.ExitStatus()
		if exitCode != 0 {
			return SubprocessFailedError{
				Args:     []string{bin, url},
				Err:      fmt.Errorf("subprocess failed with exit code %d", exitCode),
				ExitCode: exitCode,
			}
		}
	}
	return nil

}
