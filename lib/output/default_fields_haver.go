package output

type DefaultFieldsHaver interface {
	DefaultFields(f Format) string
}
