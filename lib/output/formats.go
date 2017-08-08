package output

// Format is a canonical name of output formats
type Format string

const (
	// List is the canonical name of the List output format
	List Format = "list"
	// Table is the canonical name of the Table output format
	Table = "table"
	// JSON is the canonical name of the JSON output format
	JSON = "json"
	// Human is the canonical name of the Human output format
	Human = "human"
	// Debug is the canonical name of the Debug output format
	Debug = "debug"
)
