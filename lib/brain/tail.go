package brain

import (
	"net"
)

// Tail represents a Bytemark Cloud Servers tail (disk storage machine), as returned by the admin API.
type Tail struct {
	ID    int
	UUID  string
	Label string

	CCAddress *net.IP
	ZoneName  string

	IsOnline     bool
	StoragePools []*StoragePool
}
