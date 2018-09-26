// BuildVirtualMachineRequest builds an authenticated request for the given
// virtual machine. The virtual machine provides the base URL 
// (like /virtual_machines/15). suffix is formatted with
// fmt.Sprintf(suffix, suffixSubs) and appended to the base URL to create the
// full request path.
func (c *bytemarkClient) BuildVirtualMachineRequest(
	method string,
	virtualMachine brain.VirtualMachinePather,
	suffix string,
	suffixSubs ...string,
) (Request, error) {
	virtualMachine, err := c.CheckVirtualMachinePather(virtualMachine)
	if err != nil {
		return nil, err
	}
	baseUrl, err := virtualMachine.VirtualMachinePath()
	if err != nil {
		return nil, err
	}
	return BuildRequest(method, BrainEndpoint, baseUrl+suffix, suffixSubs...)
}
