package mocks

import (
	"bigv.io/client/cmds/util"
	mock "github.com/maraino/go-mock"
)

// mock CommandSet

type Commands struct {
	mock.Mock
}

func (cmds *Commands) Console(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) Debug(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) Delete(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) DeleteVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) DeleteGroup(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) Help(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) Config(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) Create(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) CreateGroup(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) CreateVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) ResetVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) Restart(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))

}
func (cmds *Commands) Shutdown(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) Stop(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) Start(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) Show(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) ShowAccount(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) ShowGroup(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) ShowVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) Undelete(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) UndeleteVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) LockHWProfile(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) UnlockHWProfile(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) SetCores(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) SetHWProfile(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) SetMemory(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) EnsureAuth() error {
	r := cmds.Called()
	return r.Error(0)
}

func (cmds *Commands) HelpForConfig() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForCreate() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForDebug() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}
func (cmds *Commands) HelpForDelete() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForHelp() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForLocks() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForPower() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForSet() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *Commands) HelpForShow() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}
