package lib

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestParseVirtualMachineName(t *testing.T) {
	is := is.New(t)

	client, _ := New("endpoint.tld", "billing.endpoint.tld", "spp.endpoint.tld")
	vm, err := client.ParseVirtualMachineName("a.b.c")
	is.Nil(err)
	is.Equal("a.b.c", vm.String())

	vm, err = client.ParseVirtualMachineName("a..c")
	is.Nil(err)
	is.Equal("a.default.c", vm.String())
	vm, err = client.ParseVirtualMachineName("a.b.c.")
	is.Equal("a.b.c", vm.String())
	is.Nil(err)
	vm, err = client.ParseVirtualMachineName("a.b.c.endpoint.tld")
	is.Equal("a.b.c", vm.String())
	is.Nil(err)
	vm, err = client.ParseVirtualMachineName("a.b.c.d.endpoint.tld")
	is.Equal("a.b.c", vm.String())
	is.Nil(err)
	vm, err = client.ParseVirtualMachineName("endpoint.tld.a.endpoint.tld")
	is.Equal("endpoint.tld.a", vm.String())
	is.Nil(err)
	vm, err = client.ParseVirtualMachineName("endpoint.tld.a.endpoint.tld.")
	is.Equal("endpoint.tld.a", vm.String())

	_, err = client.ParseVirtualMachineName(".b.c")
	is.NotNil(err)
	_, err = client.ParseVirtualMachineName(".")
	is.NotNil(err)

}

func TestParseGroupName(t *testing.T) {
	is := is.New(t)

	client, _ := New("endpoint.tld", "billing.endpoint.tld", "spp.endpoint.tld")

	is.Equal("halloween-vms.spooky-steve", client.ParseGroupName("halloween-vms.spooky-steve").String())
	is.Equal("a.b", client.ParseGroupName("a.b").String())
	is.Equal("a.b", client.ParseGroupName("a.b.c").String())
	is.Equal("a.b", client.ParseGroupName("a.b.c.").String())
	is.Equal("a.b", client.ParseGroupName("a.b.c.endpoint.tld").String())
	is.Equal("a.b", client.ParseGroupName("a.b.c.d.endpoint.tld").String())
	is.Equal("endpoint.tld", client.ParseGroupName("endpoint.tld.a.endpoint.tld").String())
	is.Equal("endpoint.tld", client.ParseGroupName("endpoint.tld.a.endpoint.tld.").String())
}

func TestParseAccountName(t *testing.T) {
	is := is.New(t)

	client, _ := New("endpoint.tld", "billing.endpoint.tld", "spp.endpoint.tld")

	is.Equal("a", client.ParseAccountName("a.b.c"))
	is.Equal("a", client.ParseAccountName("a.b.c."))
	is.Equal("a", client.ParseAccountName("a.b.c.endpoint.tld"))
	is.Equal("a", client.ParseAccountName("a.b.c.d.endpoint.tld"))
	is.Equal("endpoint", client.ParseAccountName("endpoint.tld.a.endpoint.tld"))
	is.Equal("endpoint", client.ParseAccountName("endpoint.tld.a.endpoint.tld."))
}

func TestParseAccountNameDefaulting(t *testing.T) {
	is := is.New(t)

	client, auth, brain, billing, err := mkTestClientAndServers(mkNilHandler(t), mkNilHandler(t))
	if err != nil {
		t.Fatal(err)
	}
	defer auth.Close()
	defer brain.Close()
	defer billing.Close()

	is.Equal("", client.ParseAccountName(""))
	is.Equal("hey", client.ParseAccountName("hey"))
	is.Equal("hey", client.ParseAccountName("hey.guys.it.me.telyn"))
	is.Equal("hey", client.ParseAccountName("hey", ""))
	is.Equal("hey", client.ParseAccountName("", "hey"))

}
