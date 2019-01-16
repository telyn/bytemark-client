package flags

// This file is used by go generate to create SliceFlag variants of other flag
// types, and tests for said SliceFlags

// -preprocesser to generate a func Preprocess(ctx *app.Context) error
//go:generate go run ./gen/slice_flags -preprocesser -o virtual_machine_name_slice.go --example-input "staples.stapler.photocopier" VirtualMachineName
//go:generate go run ./gen/slice_flags -preprocesser -o group_name_slice.go --example-input "staples.stapler" GroupName
//go:generate go run ./gen/slice_flags -preprocesser -o account_name_slice.go --example-input "photocopier" AccountName
