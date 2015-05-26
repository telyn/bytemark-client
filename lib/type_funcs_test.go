package lib

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestParseVirtualMachineName(t *testing.T) {
	is := is.New(t)

	bigv, _ := New("endpoint.tld")

	is.Equal("a.b.c", bigv.ParseVirtualMachineName("a.b.c").String())
	is.Equal("a.b.c", bigv.ParseVirtualMachineName("a.b.c.").String())
	is.Equal("a.b.c", bigv.ParseVirtualMachineName("a.b.c.endpoint.tld").String())
	is.Equal("a.b.c", bigv.ParseVirtualMachineName("a.b.c.d.endpoint.tld").String())
	is.Equal("endpoint.tld.a", bigv.ParseVirtualMachineName("endpoint.tld.a.endpoint.tld").String())
	is.Equal("endpoint.tld.a", bigv.ParseVirtualMachineName("endpoint.tld.a.endpoint.tld.").String())
}

func TestParseGroupName(t *testing.T) {
	is := is.New(t)

	bigv, _ := New("endpoint.tld")

	is.Equal("a.b", bigv.ParseGroupName("a.b.c").String())
	is.Equal("a.b", bigv.ParseGroupName("a.b.c.").String())
	is.Equal("a.b", bigv.ParseGroupName("a.b.c.endpoint.tld").String())
	is.Equal("a.b", bigv.ParseGroupName("a.b.c.d.endpoint.tld").String())
	is.Equal("endpoint.tld", bigv.ParseGroupName("endpoint.tld.a.endpoint.tld").String())
	is.Equal("endpoint.tld", bigv.ParseGroupName("endpoint.tld.a.endpoint.tld.").String())
}

func TestParseAccountName(t *testing.T) {
	is := is.New(t)

	bigv, _ := New("endpoint.tld")

	is.Equal("a", bigv.ParseAccountName("a.b.c"))
	is.Equal("a", bigv.ParseAccountName("a.b.c."))
	is.Equal("a", bigv.ParseAccountName("a.b.c.endpoint.tld"))
	is.Equal("a", bigv.ParseAccountName("a.b.c.d.endpoint.tld"))
	is.Equal("endpoint", bigv.ParseAccountName("endpoint.tld.a.endpoint.tld"))
	is.Equal("endpoint", bigv.ParseAccountName("endpoint.tld.a.endpoint.tld."))
}
