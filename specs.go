package client

type UserSpec struct {
	username        string
	email           string
	authorized_keys string
	password        string
	account_name    string
}

type AccountSpec struct {
	name string
}

type GroupSpec struct {
	account_id int
	name       string
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
}

type DiskSpec struct {
	label         string
	storage_grade string
	size          int
}

type ImageSpec struct {
	distribution  string
	root_password string
}

type VirtualMachineExtendedSpec struct {
	virtual_machine VirtualMachineSpec
	discs           []DiskSpec
	reimage         ImageSpec
}
