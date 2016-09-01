package util

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func CallBrowser(url string) error {
	log.Logf("Running a browser to open %s...\r\n", url)

	var attr os.ProcAttr
	attr.Sys = &syscall.SysProcAttr{HideWindow: true}
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	path, err := exec.LookPath("cmd")
	if err != nil {
		return err
	}
	// so on windows when you're using cmd you have to escape ampersands with the ^ character.
	// ¯\(º_o)/¯
	url = strings.Replace(url, "&", "^&", -1)

	log.Debugf(log.LvlOutline, "Executing %s /C start \"%s\"", path, url)
	proc, err := os.StartProcess(path, []string{path, "/C", "start", url}, &attr)
	if err != nil {

		return err
	}

	_, err = proc.Wait()
	return err
}
