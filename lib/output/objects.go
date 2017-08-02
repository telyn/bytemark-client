package output

// Object is a kind of thing that can be output. The Objects defined in this package are the complete set of types that can be passed to Output
type Object int

const (
	Account Object = iota
	Backup
	BackupSchedule
	Disc
	Group
	Privilege
	Server
	Head
	Tail
	StoragePool
	IPRange
	VLAN
	Definition
)
