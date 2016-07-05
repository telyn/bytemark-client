package lib

import (
	"bytes"
	"github.com/cheekybits/is"
	"testing"
)

func TestFormatVM(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)
	vm, _, _ := getFixtureVMWithManyIPs()

	err := FormatVirtualMachine(b, &vm, "serversummary")
	if err != nil {
		t.Error(err)
	}

	is.Equal(" ▸ valid-vm.default (powered on) in Default", b.String())
	b.Truncate(0)
	err = FormatVirtualMachine(b, &vm, "serverspec")
	if err != nil {
		t.Error(err)
	}
	is.Equal("   192.168.1.16 - 1 core, 1MiB, 25GiB on 1 disc", b.String())
}
