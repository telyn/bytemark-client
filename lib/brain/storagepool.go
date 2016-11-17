package brain

// StoragePool represents a Bytemark Cloud Servers disk storage pool, as returned by the admin API.
type StoragePool struct {
	ID            int
	Label         string
	Tail          *Tail
	Name          string
	Space         int
	FreeSpace     int
	StorageGrade  string
	IOPSLimit     int
	UsageStrategy string
	Note          string
}
