package util

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func callBrowser(url string) error {
	fmt.Fprintf(os.Stderr, "Running a browser to open %s...\r\n", url)

	var attr os.ProcAttr
	attr.Sys = &syscall.SysProcAttr{HideWindow: true}

	path, err := exec.LookPath("cmd")
	if err != nil {
		return err
	}
	proc, err := os.StartProcess(path, []string{"start", url}, &attr)
	if err != nil {

		return err
	}

	_, err = proc.Wait()
	return err
}
