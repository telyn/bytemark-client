package brain

// VirtualMachineSpec represents the specification for a VM that is passed to the create_vm endpoint
type VirtualMachineSpec struct {
	VirtualMachine *VirtualMachine `json:"virtual_machine"`
	Discs          []Disc          `json:"discs"`
	Reimage        *ImageInstall   `json:"reimage"`
	IPs            *IPSpec         `json:"ips"`
}
