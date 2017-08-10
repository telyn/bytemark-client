package output

// Config is a simple struct to configure whichever OutputFormatFn is selected
// all its fields should be set regardless of Format.
type Config struct {
	Fields []string
	Format Format
}
