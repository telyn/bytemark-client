// +build quality

package lib_test

import (
	"go/importer"
	"go/types"
	"sort"
	"testing"
)

// TestQualityInterfaceHasntGrown will fail if the Client interface has gained any new
// functions since the 16th Mar 2018 (Impersonate was added)
// See billing/updatedefinitions.go for an example of how ordinary request functions
// should be written from now.
func TestQualityInterfaceHasntGrown(t *testing.T) {
	allowed := sort.StringSlice{
		"AddIP",
		"AdminCreateGroup",
		"AllowInsecureRequests",
		"ApproveVM",
		"AuthWithCredentials",
		"AuthWithToken",
		"BuildRequest",
		"BuildRequestNoAuth",
		"CancelDiscMigration",
		"CancelVMMigration",
		"CreateBackup",
		"CreateBackupSchedule",
		"CreateCreditCard",
		"CreateCreditCardWithToken",
		"CreateDisc",
		"CreateGroup",
		"CreateIPRange",
		"CreateUser",
		"CreateVirtualMachine",
		"DeleteBackup",
		"DeleteBackupSchedule",
		"DeleteDisc",
		"DeleteGroup",
		"DeleteVLAN",
		"DeleteVirtualMachine",
		"EmptyHead",
		"EmptyStoragePool",
		"EnsureAccountName",
		"EnsureGroupName",
		"EnsureVirtualMachineName",
		"GetAccount",
		"GetAccounts",
		"GetBackups",
		"GetDefaultAccount",
		"GetDisc",
		"GetDiscByID",
		"GetEndpoint",
		"GetGroup",
		"GetHead",
		"GetHeads",
		"GetIPRange",
		"GetIPRanges",
		"GetMigratingDiscs",
		"GetMigratingVMs",
		"GetPrivileges",
		"GetPrivilegesForAccount",
		"GetPrivilegesForGroup",
		"GetPrivilegesForVirtualMachine",
		"GetRecentVMs",
		"GetSPPToken",
		"GetSessionFactors",
		"GetSessionToken",
		"GetSessionUser",
		"GetStoppedEligibleVMs",
		"GetStoragePool",
		"GetStoragePools",
		"GetTail",
		"GetTails",
		"GetUser",
		"GetVLAN",
		"GetVLANs",
		"GetVirtualMachine",
		"GrantPrivilege",
		"Impersonate",
		"MigrateDisc",
		"MigrateVirtualMachine",
		"MoveVirtualMachine",
		"ReadDefinitions",
		"ReapVMs",
		"RegisterNewAccount",
		"RegradeDisc",
		"ReifyDisc",
		"ReimageVirtualMachine",
		"RejectVM",
		"ResetVirtualMachine",
		"ResizeDisc",
		"RestartVirtualMachine",
		"RestoreBackup",
		"RevokePrivilege",
		"SetDebugLevel",
		"SetDiscIopsLimit",
		"SetVirtualMachineCDROM",
		"SetVirtualMachineCores",
		"SetVirtualMachineHardwareProfile",
		"SetVirtualMachineHardwareProfileLock",
		"SetVirtualMachineMemory",
		"ShutdownVirtualMachine",
		"StartVirtualMachine",
		"StopVirtualMachine",
		"UndeleteVirtualMachine",
		"UpdateHead",
		"UpdateStoragePool",
		"UpdateTail",
		"UpdateVMMigration",
	}
	pkg, err := importer.Default().Import("github.com/BytemarkHosting/bytemark-client/lib")
	if err != nil {
		t.Fatal(err)
	}
	obj := pkg.Scope().Lookup("Client")
	// I don't know why we need the underlying type in order to cast it
	// to a *types.Interface, but we do... soo...

	// #GoTypesIsBlackMagic
	// i guess that makes me a witch. i'm ok with that
	if iface, ok := obj.Type().Underlying().(*types.Interface); ok {
		for i := 0; i < iface.NumMethods(); i++ {
			name := iface.Method(i).Name()
			if allowed[allowed.Search(name)] != name {
				t.Errorf("New method on the Client interface called %s. The Client interface is not allowed to get any bigger - instead, define a function that receives a Client as an argument, and place it in a relevant package. See billing.UpdateDefinitions as an example.", name)
			}
		}
	} else {
		t.Fatalf("Couldn't cast %s to types.Interface", obj)
	}
}
