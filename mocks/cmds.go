package mocks

import (
	"bytemark.co.uk/client/cmds/util"
	mock "github.com/maraino/go-mock"
)

// mock CommandSet

type Commands struct {
	mock.Mock
}

func (cmds *Commands) Console(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) Debug(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) Delete(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) DeleteVM(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) DeleteGroup(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) DeleteDisc(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) Help(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) Config(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) CreateDiscs(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) CreateGroup(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) CreateVM(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ListAccounts(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ListDiscs(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ListGroups(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ListVMs(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ResetVM(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ResizeDisc(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) Restart(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))

}
func (cmds *Commands) Shutdown(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) Stop(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) Start(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) Show(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ShowAccount(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ShowGroup(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) ShowVM(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) Undelete(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) UndeleteVM(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) LockHWProfile(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) UnlockHWProfile(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) SetCores(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) SetHWProfile(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) SetMemory(args []string) util.ExitCode {
	r := cmds.Called(args)
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) EnsureAuth() error {
	r := cmds.Called()
	return r.Error(0)
}

func (cmds *Commands) HelpForConfig() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForCreate() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForDebug() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}
func (cmds *Commands) HelpForDelete() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForHelp() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForList() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForLocks() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForPower() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForResize() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForSet() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForShow() util.ExitCode {
	r := cmds.Called()
	return util.ExitCode(r.Int(0))
}
