package mocks

import (
	"bigv.io/client/cmd"
	mock "github.com/maraino/go-mock"
)

// mock CommandSet

type mockCommandManager struct {
	mock.Mock
}

func (cmds *mockCommandManager) Console(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) Debug(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) Delete(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) DeleteVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) DeleteGroup(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) Help(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) Config(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) Create(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) CreateGroup(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) CreateVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) ResetVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) Restart(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))

}
func (cmds *mockCommandManager) Shutdown(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) Stop(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) Start(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) Show(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) ShowAccount(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) ShowGroup(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) ShowVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) Undelete(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) UndeleteVM(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) LockHWProfile(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) UnlockHWProfile(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) SetCores(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) SetHWProfile(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) SetMemory(args []string) cmd.ExitCode {
	r := cmds.Called(args)
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) EnsureAuth() error {
	r := cmds.Called()
	return r.Error(0)
}

func (cmds *mockCommandManager) HelpForConfig() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForCreate() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForDebug() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}
func (cmds *mockCommandManager) HelpForDelete() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForHelp() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForLocks() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForPower() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForSet() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}

func (cmds *mockCommandManager) HelpForShow() cmd.ExitCode {
	r := cmds.Called()
	return cmd.ExitCode(r.Int(0))
}
