package lib

type DiskSpec struct {
	label         string
	storage_grade string
	size          int

	id                 int
	virtual_machine_id int
	storage_pool       string
}

type ImageSpec struct {
	distribution  string
	root_password string
}

type IPSpec struct {
	rdns string

	// this cannot be set.
	ip string
}

type NetworkInterfaceSpec struct {
	label string
	mac   string

	// the following can't be set (or at least, so I'm assuming..)

	id                 int
	vlan_num           int
	ips                []string
	extra_ips          []string
	virtual_machine_id int
}

type UserSpec struct {
	username        string
	email           string
	authorized_keys string
	password        string

	// "users can be created (using POST) without authentication. If the
	// request has no authentication, it will also accept an account_name
	// parameter and create an account at the same time."
	account_name string
}

type VirtualMachineSpec struct {
	autoreboot_on           bool
	cdrom_url               string
	cores                   int
	memory                  int
	name                    string
	power_on                string
	hardware_profile        string
	hardware_profile_locked bool
	group_id                int

	// zone name can be set during creation but not changed
	zone_name string

	// the following cannot be set
	id                 int
	management_address string
	deleted            bool
	hostname           string
	head               string
	network_interfaces []NetworkInterfaceSpec
}

type VirtualMachineExtendedSpec struct {
	virtual_machine VirtualMachineSpec
	discs           []DiskSpec
	reimage         ImageSpec
}

type GroupSpec struct {
	name string

	// the following cannot be set
	account_id       int
	id               int
	virtual_machiens []VirtualMachineSpec
}

type AccountSpec struct {
	name string

	// the following cannot be set
	id        int
	suspended bool
	groups    []GroupSpec
}
