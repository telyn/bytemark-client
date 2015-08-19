// +build !windows

package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	fmt.Fprintf(os.Stderr, "Running a browser to open %s...", url)
	command := openCommand()
	bin, err := exec.LookPath(command)
	if err != nil {

		command = "/usr/bin/x-www-browser"
		bin, err = exec.LookPath(command)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(os.Stderr, "%s %s\r\n", bin, url)

	var attr os.ProcAttr
	proc, err := os.StartProcess(bin, []string{bin, url}, &attr)
	if err != nil {
		return err
	}
	_, err = proc.Wait()
	return err

}
