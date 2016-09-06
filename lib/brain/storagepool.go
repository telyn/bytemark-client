package brain

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
