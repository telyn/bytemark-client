package output

// DefaultFieldsHaver is an interface that must be implemented in order to pass objects to output.Write
type DefaultFieldsHaver interface {
	// DefaultFields must return a string of valid field names (according to github.com/BytemarkHosting/row - i.e. they must be in the list output by row.FieldsFrom) for the type it is implemented on.
	// It is used to discover what fields should be output by output.Write when there's no list of user-specified fields.
	DefaultFields(f Format) string
}
