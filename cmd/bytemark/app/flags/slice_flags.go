package flags

// This file is used by go generate to create SliceFlag variants of other flag
// types, and tests for said SliceFlags

//go:generate go run gen/slice_flags.go -o slice_flags_generated.go -t ./slice_template.go.tmpl VirtualMachineName GroupName AccountName
