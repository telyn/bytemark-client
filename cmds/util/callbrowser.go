// +build !windows

package util

import (
	"bigv.io/client/util/log"
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
	command := openCommand()
	bin, err := exec.LookPath(command)
	if err != nil {

		command = "/usr/bin/x-www-browser"
		bin, err = exec.LookPath(command)
		if err != nil {
			return err
		}
	}

	log.Logf("Running a browser to open %s...\r\n", url)

	var attr os.ProcAttr
	proc, err := os.StartProcess(bin, []string{bin, url}, &attr)
	log.Debugf(1, "Executing %s \"%s\"", bin, url)
	if err != nil {
		return err
	}
	_, err = proc.Wait()
	return err

}
