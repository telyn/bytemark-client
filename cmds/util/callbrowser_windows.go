package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func CallBrowser(url string) error {
	fmt.Fprintf(os.Stderr, "Running a browser to open %s...\r\n", url)

	var attr os.ProcAttr
	attr.Sys = &syscall.SysProcAttr{HideWindow: false}
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	path, err := exec.LookPath("cmd")
	if err != nil {
		return err
	}
	// so on windows when you're using cmd you have to escape ampersands with the ^ character.
	// ¯\(º_o)/¯
	url = strings.Replace(url, "&", "^&", -1)
	proc, err := os.StartProcess(path, []string{path, "/C", "start", url }, &attr)
	if err != nil {

		return err
	}

	_, err = proc.Wait()
	return err
}
