package brain

// VirtualMachineSpec represents the specification for a VM that is passed to the create_vm endpoint
type VirtualMachineSpec struct {
	VirtualMachine *VirtualMachine `json:"virtual_machine"`
	Discs          []Disc          `json:"discs,omitempty"`
	Reimage        *ImageInstall   `json:"reimage,omitempty"`
	IPs            *IPSpec         `json:"ips,omitempty"`
}
