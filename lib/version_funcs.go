package lib

import (
	"fmt"
)

type Version struct {
	MajorVersion int
	MinorVersion int
	BuildNumber  int
	GitCommit    string
	BuildDate    string
}

func GetVersion() *Version {
	return &Version{
		MajorVersion: majorversion,
		MinorVersion: minorversion,
		BuildNumber:  buildnumber,
		GitCommit:    gitcommit,
		BuildDate:    builddate,
	}
}

func (v *Version) String() string {
	if v.BuildNumber == 0 {
		return fmt.Sprintf("%d.%d~%s (%s) custom build", v.MajorVersion, v.MinorVersion, v.GitCommit, v.BuildDate)
	} else {
		return fmt.Sprintf("%d.%d.%d", v.MajorVersion, v.MinorVersion, v.BuildNumber)
	}
}
