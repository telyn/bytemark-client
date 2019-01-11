package lib
import "github.com/BytemarkHosting/bytemark-client/lib/brain"

// BuildAccountRequest builds an authenticated request for the given
// account. The account provides the base URL 
// (like /accounts/15). suffix is formatted with
// fmt.Sprintf(suffix, suffixSubs) and appended to the base URL to create the
// full request path.
func (c *bytemarkClient) BuildAccountRequest(
	method string,
	account brain.AccountPather,
	suffix string,
	suffixSubs ...string,
) (Request, error) {
	account, err := c.checkAccountPather(account)
	if err != nil {
		return nil, err
	}
	baseUrl, err := account.AccountPath()
	if err != nil {
		return nil, err
	}
	return BuildRequest(method, BrainEndpoint, baseUrl+suffix, suffixSubs...)
}


// BuildDiscRequest builds an authenticated request for the given
// disc. The disc provides the base URL 
// (like /discs/15). suffix is formatted with
// fmt.Sprintf(suffix, suffixSubs) and appended to the base URL to create the
// full request path.
func (c *bytemarkClient) BuildDiscRequest(
	method string,
	disc brain.DiscPather,
	suffix string,
	suffixSubs ...string,
) (Request, error) {
	disc, err := c.checkDiscPather(disc)
	if err != nil {
		return nil, err
	}
	baseUrl, err := disc.DiscPath()
	if err != nil {
		return nil, err
	}
	return BuildRequest(method, BrainEndpoint, baseUrl+suffix, suffixSubs...)
}


// BuildGroupRequest builds an authenticated request for the given
// group. The group provides the base URL 
// (like /groups/15). suffix is formatted with
// fmt.Sprintf(suffix, suffixSubs) and appended to the base URL to create the
// full request path.
func (c *bytemarkClient) BuildGroupRequest(
	method string,
	group brain.GroupPather,
	suffix string,
	suffixSubs ...string,
) (Request, error) {
	group, err := c.checkGroupPather(group)
	if err != nil {
		return nil, err
	}
	baseUrl, err := group.GroupPath()
	if err != nil {
		return nil, err
	}
	return BuildRequest(method, BrainEndpoint, baseUrl+suffix, suffixSubs...)
}


// BuildVirtualMachineRequest builds an authenticated request for the given
// virtual_machine. The virtual_machine provides the base URL 
// (like /virtual_machines/15). suffix is formatted with
// fmt.Sprintf(suffix, suffixSubs) and appended to the base URL to create the
// full request path.
func (c *bytemarkClient) BuildVirtualMachineRequest(
	method string,
	virtualMachine brain.VirtualMachinePather,
	suffix string,
	suffixSubs ...string,
) (Request, error) {
	virtualMachine, err := c.checkVirtualMachinePather(virtualMachine)
	if err != nil {
		return nil, err
	}
	baseUrl, err := virtualMachine.VirtualMachinePath()
	if err != nil {
		return nil, err
	}
	return BuildRequest(method, BrainEndpoint, baseUrl+suffix, suffixSubs...)
}

