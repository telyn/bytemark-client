package bigv

import (
	"net"
)

type Tail struct {
	ID    int
	UUID  string
	Label string

	CCAddress *net.IP
	ZoneName  string

	IsOnline     bool
	StoragePools []*StoragePool
}
