package lib

type VirtualMachineName struct {
	VirtualMachine string
	Group          string
	Account        string
}
type GroupName struct {
	Group   string
	Account string
}

type Disk struct {
	Label        string `json:"label"`
	StorageGrade string `json:"storage_grade"`
	Size         int    `json:"size"`

	Id               int    `json:"id"`
	VirtualMachineId int    `json:"virtual_machine_id"`
	StoragePool      string `json:"storage_pool"`
}

type ImageInstall struct {
	Distribution string `json:"distribution"`
	RootPassword string `json:"root_password"`
}

type IP struct {
	RDns string `json:"rdns"`

	// this cannot be set.
	Ip string `json:"ip"`
}

type NetworkInterface struct {
	Label string `json:"label"`

	Mac string `json:"mac"`

	// the following can't be set (or at least, so I'm assuming..)

	Id               int               `json:"id"`
	VlanNum          int               `json:"vlan_num"`
	Ips              []string          `json:"ips"`
	ExtraIps         map[string]string `json:"extra_ips"`
	VirtualMachineId int               `json:"virtual_machine_id"`
}

type User struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	AuthorizedKeys string `json:"authorized_keys"`
	Password       string `json:"password"`

	// "users can be created (using POST) without authentication. If the
	// request has no authentication, it will also accept an account_name
	// parameter and create an account at the same time."
	AccountName string `json:"account_name"`
}

// TODO(telyn): new fields (last_imaged_with and
type VirtualMachine struct {
	Autoreboot            bool   `json:"autoreboot_on"`
	CdromUrl              string `json:"cdrom_url"`
	Cores                 int    `json:"cores"`
	Memory                int    `json:"memory"`
	Name                  string `json:"name"`
	PowerOn               bool   `json:"power_on"`
	HardwareProfile       string `json:"hardware_profile"`
	HardwareProfileLocked bool   `json:"hardware_profile_locked"`
	GroupId               int    `json:"group_id"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name"`

	// the following cannot be set
	Discs             []*Disk             `json:"discs"`
	Id                int                 `json:"id"`
	ManagementAddress string              `json:"management_address"`
	Deleted           bool                `json:"deleted"`
	Hostname          string              `json:"hostname"`
	Head              string              `json:"head"`
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces"`
}

type VirtualMachineSpec struct {
	VirtualMachine VirtualMachine `json:"virtual_machine"`
	Discs          []Disk         `json:"discs"`
	Reimage        ImageInstall   `json:"reimage"`
}

type Group struct {
	Name string `json:name"`

	// the following cannot be set
	AccountId       int              `json:"account_id"`
	Id              int              `json:"id"`
	VirtualMachines []VirtualMachine `json:"virtual_machines"`
}

type Account struct {
	Name string `json:"name"`

	// the following cannot be set
	Id        int      `json:"id"`
	Suspended bool     `json:"suspended"`
	Groups    []*Group `json:"groups"`
}
