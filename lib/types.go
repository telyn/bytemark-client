package lib

import (
	"net"
)

// Disc is a representation of a VM's disc.
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
	Distribution    string `json:"distribution"`
	FirstbootScript string `json:"firstboot_script"`
	RootPassword    string `json:"root_password"`
	PublicKeys      string `json:"ssh_public_key"`
}

// IP represents an IP for the purpose of setting RDNS
type IP struct {
	RDns string `json:"rdns"`

	// this cannot be set.
	IP *net.IP `json:"ip"`
}

// IPSpec represents one v4 and one v6 address to assign to a server during creation.
type IPSpec struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

// IPCreateRequest is used by the create_ip endpoint on a nic
type IPCreateRequest struct {
	Addresses  int    `json:"addresses"`
	Family     string `json:"family"`
	Reason     string `json:"reason"`
	Contiguous bool   `json:"contiguous"`
	// don't actually specify the IPs - this is for filling in from the response!
	IPs []*net.IP `json:"ips"`
}

// JSONUser is used as an intermediate type that gets processed into a User. It should not have been exported.
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

// CreditCard represents all the data for a credit card with SPP.
type CreditCard struct {
	Number   string `yaml:"account_number"`
	Name     string `yaml:"name"`
	Expiry   string `yaml:"expiry"`
	CVV      string `yaml:"cvv"`
	Street   string `yaml:"street,omitempty"`
	City     string `yaml:"city,omitempty"`
	Postcode string `yaml:"postcode,omitempty"`
	Country  string `yaml:"country,omitempty"`
}
