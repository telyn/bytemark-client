package lib

//go:generate ../gen/version.go.sh

import (
	"fmt"
	"strings"
)

type Version struct {
	MajorVersion int
	MinorVersion int
	BuildNumber  int
	GitCommit    string
	GitBranch    string
	BuildDate    string
}

func GetVersion() *Version {
	return &Version{
		MajorVersion: majorversion,
		MinorVersion: minorversion,
		BuildNumber:  buildnumber,
		GitCommit:    gitcommit,
		GitBranch:    gitbranch,
		BuildDate:    builddate,
	}
}

func (v *Version) String() string {
	if v.BuildNumber == 0 {
		return fmt.Sprintf("%d.%d~%s (%s) custom build from %s", v.MajorVersion, v.MinorVersion, v.GitCommit, v.BuildDate, v.GitBranch)
	} else {
		if strings.HasPrefix(v.GitBranch, "release-") {
			return fmt.Sprintf("%d.%d.%d", v.MajorVersion, v.MinorVersion, v.BuildNumber)
		} else {
			return fmt.Sprintf("%d.%d.%d~%s (%s) automated build from %s", v.MajorVersion, v.MinorVersion, v.BuildNumber, v.GitCommit, v.BuildDate, v.GitBranch)
		}
	}
}
