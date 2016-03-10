package lib

//go:generate ../gen/version.go.sh

import (
	"fmt"
)

type Version struct {
	MajorVersion int
	MinorVersion int
	Revision     int
	BuildNumber  int
	GitCommit    string
	GitBranch    string
	BuildDate    string
}

func GetVersion() *Version {
	return &Version{
		MajorVersion: majorversion,
		MinorVersion: minorversion,
		Revision:     revision,
		BuildNumber:  buildnumber,
		GitCommit:    gitcommit,
		GitBranch:    gitbranch,
		BuildDate:    builddate,
	}
}

func (v *Version) FullString() string {
	return fmt.Sprintf("%d.%d.%d\r\nBuild number: %d\r\nBuild date:%s\r\nGit commit: %s %s", v.MajorVersion, v.MinorVersion, v.Revision, v.BuildNumber, v.BuildDate, v.GitBranch, v.GitCommit)
}

func (v *Version) String() string {
	if v.BuildNumber == 0 {
		return fmt.Sprintf("%d.%d.%d~%s (%s) manual build from %s", v.MajorVersion, v.MinorVersion, v.Revision, v.GitCommit, v.BuildDate, v.GitBranch)
	} else {
		if v.GitBranch == "master" {
			return fmt.Sprintf("%d.%d.%d", v.MajorVersion, v.MinorVersion, v.Revision)
		} else {
			return fmt.Sprintf("%d.%d.%d~%s (%s) automated build #%d from %s", v.MajorVersion, v.MinorVersion, v.Revision, v.GitCommit, v.BuildDate, v.BuildNumber, v.GitBranch)
		}
	}
}
