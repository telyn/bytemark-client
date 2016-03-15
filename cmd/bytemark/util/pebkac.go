package util

// TODO(telyn): Delete this file before 0.6

type PEBKACError struct{}

func (err PEBKACError) Error() string {
	return "Yo you did a thing wrong"
}
