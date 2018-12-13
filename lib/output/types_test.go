// +build quality

package output_test

// This file is used by go_generate to create a test from types_test.go.inc
// using the generator code in gen/list_types/list_types.go

// The test is to make sure that all types in billing, brain and spp, and
// the Account type (and Overview type if I ever make one) all implement
// PrettyPrinter and DefaultFieldsHaver

// if reflect ever gains support for listing all types in a package
// then this and types_test.go.inc can be rewritten into a single file in lib/output

// the go install step is needed because gen/list_types/list_types.go uses the compiled libraries
// and is not able to use the source.

//go:generate go install "github.com/BytemarkHosting/bytemark-client/lib" "github.com/BytemarkHosting/bytemark-client/lib/..."
//go:generate go run ../../gen/list_types/list_types.go -f "(*%s)(nil)," -o generated_types_test.go -t types_test.go.inc "github.com/BytemarkHosting/bytemark-client/lib/brain" "github.com/BytemarkHosting/bytemark-client/lib/billing" "github.com/BytemarkHosting/bytemark-client/lib/spp"
