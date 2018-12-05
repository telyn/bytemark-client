package flags

// A Preprocesser is a Flag that has a preprocess step that requires a Context
type Preprocesser interface {
	Preprocess(c *Context) error
}
