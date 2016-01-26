package lib

import ()

// VirtualMachineName is the triplet-form of the name of a VirtualMachine, which should be enough to find the VM.
type VirtualMachineName struct {
	VirtualMachine string
	Group          string
	Account        string
}

// GroupName is the double-form of the name of a Group, which should be enough to find the group.
type GroupName struct {
	Group   string
	Account string
}

// Disc is a representation of a VM's disc, which can be unmarshalled from JSON output by the BigV server.
type Disc struct {
	Label        string `json:"label"`
	StorageGrade string `json:"storage_grade"`
	Size         int    `json:"size"`

	ID               int    `json:"id"`
	VirtualMachineID int    `json:"virtual_machine_id"`
	StoragePool      string `json:"storage_pool"`
}

// ImageInstall represents what image was most recently installed on a VM along with its root password.
// This might only be returned when creating a VM.
type ImageInstall struct {
	Distribution string `json:"distribution"`
	RootPassword string `json:"root_password"`
	PublicKeys   string `json:"ssh_public_key"`
}

// IP represents an IP for the purpose of setting RDNS with BigV
type IP struct {
	RDns string `json:"rdns"`

	// this cannot be set.
	IP string `json:"ip"`
}

// IPSpec represents one v4 and one v6 address to assign to a bigv during creation.
type IPSpec struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

// NetworkInterface represents a BigV virtual NIC and what IPs it has routed.
type NetworkInterface struct {
	Label string `json:"label"`

	Mac string `json:"mac"`

	// the following can't be set (or at least, so I'm assuming..)

	ID               int               `json:"id"`
	VlanNum          int               `json:"vlan_num"`
	IPs              []string          `json:"ips"`
	ExtraIPs         map[string]string `json:"extra_ips"`
	VirtualMachineID int               `json:"virtual_machine_id"`
}

type JSONUser struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	AuthorizedKeys string `json:"authorized_keys"`

	// passwords are handled by auth these days
	//Password       string `json:"password"`

	// "users can be created (using POST) without authentication. If the
	// request has no authentication, it will also accept an account_name
	// parameter and create an account at the same time."
	// this is almost certainly never going to be useful
	//AccountName string `json:"account_name"`
}

// User represents a Bytemark user.
type User struct {
	Username       string
	Email          string
	AuthorizedKeys []string
}

// VirtualMachine represents a VirtualMachine, as passed around from the BigV virtual_machines endpoint
type VirtualMachine struct {
	Autoreboot            bool   `json:"autoreboot_on"`
	CdromURL              string `json:"cdrom_url"`
	Cores                 int    `json:"cores"`
	Memory                int    `json:"memory"`
	Name                  string `json:"name"`
	PowerOn               bool   `json:"power_on"`
	HardwareProfile       string `json:"hardware_profile"`
	HardwareProfileLocked bool   `json:"hardware_profile_locked"`
	GroupID               int    `json:"group_id"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name"`

	// the following cannot be set
	Discs             []*Disc             `json:"discs"`
	ID                int                 `json:"id"`
	ManagementAddress string              `json:"management_address"`
	Deleted           bool                `json:"deleted"`
	Hostname          string              `json:"hostname"`
	Head              string              `json:"head"`
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces"`

	// TODO(telyn): new fields (last_imaged_with and there is another but I forgot)
}

// VirtualMachineSpec represents the specification for a VM that is passed to the create_vm endpoint
type VirtualMachineSpec struct {
	VirtualMachine *VirtualMachine `json:"virtual_machine"`
	Discs          []Disc          `json:"discs"`
	Reimage        *ImageInstall   `json:"reimage"`
	IPs            *IPSpec         `json:"ips"`
}

// Group represents a group
type Group struct {
	Name string `json:"name"`

	// the following cannot be set
	AccountID       int               `json:"account_id"`
	ID              int               `json:"id"`
	VirtualMachines []*VirtualMachine `json:"virtual_machines"`
}

// Account represents a BigV account
type Account struct {
	Name string `json:"name"`

	// the following cannot be set
	ID        int      `json:"id"`
	Suspended bool     `json:"suspended"`
	Groups    []*Group `json:"groups"`
}
