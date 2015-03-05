package lib

type Disk struct {
	label         string
	storage_grade string
	size          int

	id                 int
	virtual_machine_id int
	storage_pool       string
}

type Image struct {
	distribution  string
	root_password string
}

type IP struct {
	rdns string

	// this cannot be set.
	ip string
}

type NetworkInterface struct {
	label string
	mac   string

	// the following can't be set (or at least, so I'm assuming..)

	id                 int
	vlan_num           int
	ips                []string
	extra_ips          []string
	virtual_machine_id int
}

type User struct {
	username        string
	email           string
	authorized_keys string
	password        string

	// "users can be created (using POST) without authentication. If the
	// request has no authentication, it will also accept an account_name
	// parameter and create an account at the same time."
	account_name string
}

type VirtualMachine struct {
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
	network_interfaces []NetworkInterface
}

type VirtualMachineExtended struct {
	virtual_machine VirtualMachine
	discs           []Disk
	reimage         Image
}

type Group struct {
	name string

	// the following cannot be set
	account_id       int
	id               int
	virtual_machiens []VirtualMachine
}

type Account struct {
	name string

	// the following cannot be set
	id        int
	suspended bool
	groups    []Group
}
