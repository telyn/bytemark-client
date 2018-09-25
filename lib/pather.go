package lib

// Pather is a type which returns a URL path to use with BuildRequestWithPather
type Pather interface {
	Path() (string, error)
}
