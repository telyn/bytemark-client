package lib

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestParseVirtualMachineName(t *testing.T) {
	is := is.New(t)

	vm, err := ParseVirtualMachineName("a.b.c")
	is.Nil(err)
	is.Equal("a.b.c", vm.String())

	vm, err = ParseVirtualMachineName("a..c")
	is.Nil(err)
	is.Equal("a.default.c", vm.String())

	vm, err = ParseVirtualMachineName("a.b.c.")
	is.Equal("a.b.c", vm.String())
	is.Nil(err)

	vm, err = ParseVirtualMachineName("a.b.c.endpoint.tld")
	is.Equal("a.b.c", vm.String())
	is.Nil(err)

	vm, err = ParseVirtualMachineName("a.b.c.d.endpoint.tld")
	is.Equal("a.b.c", vm.String())
	is.Nil(err)

	vm, err = ParseVirtualMachineName("endpoint.tld.a.endpoint.tld")
	is.Equal("endpoint.tld.a", vm.String())
	is.Nil(err)

	vm, err = ParseVirtualMachineName("endpoint.tld.a.endpoint.tld.")
	is.Nil(err)
	is.Equal("endpoint.tld.a", vm.String())

	_, err = ParseVirtualMachineName(".b.c")
	is.NotNil(err)

	_, err = ParseVirtualMachineName(".")
	is.NotNil(err)

}

func TestParseGroupName(t *testing.T) {
	is := is.New(t)

	is.Equal("halloween-vms.spooky-steve", ParseGroupName("halloween-vms.spooky-steve").String())
	is.Equal("a.b", ParseGroupName("a.b").String())
	is.Equal("a.b", ParseGroupName("a.b.c").String())
	is.Equal("a.b", ParseGroupName("a.b.c.").String())
	is.Equal("a.b", ParseGroupName("a.b.c.endpoint.tld").String())
	is.Equal("a.b", ParseGroupName("a.b.c.d.endpoint.tld").String())
	is.Equal("endpoint.tld", ParseGroupName("endpoint.tld.a.endpoint.tld").String())
	is.Equal("endpoint.tld", ParseGroupName("endpoint.tld.a.endpoint.tld.").String())
}

func TestParseAccountName(t *testing.T) {
	is := is.New(t)

	is.Equal("a", ParseAccountName("a.b.c"))
	is.Equal("a", ParseAccountName("a.b.c."))
	is.Equal("a", ParseAccountName("a.b.c.endpoint.tld"))
	is.Equal("a", ParseAccountName("a.b.c.d.endpoint.tld"))
	is.Equal("endpoint", ParseAccountName("endpoint.tld.a.endpoint.tld"))
	is.Equal("endpoint", ParseAccountName("endpoint.tld.a.endpoint.tld."))
}

func TestParseAccountNameDefaulting(t *testing.T) {
	is := is.New(t)

	is.Equal("", ParseAccountName(""))
	is.Equal("hey", ParseAccountName("hey"))
	is.Equal("hey", ParseAccountName("hey.guys.it.me.telyn"))
	is.Equal("hey", ParseAccountName("hey", ""))
	is.Equal("hey", ParseAccountName("", "hey"))

}
