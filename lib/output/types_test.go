package output_test

// really telyn? is this how you're going to write this test?
// fmt.Sprintf ing this file with a list of interfaces?
// A: yes, yes it is :\

// if reflect ever gains support for listing all types in a package
// then this can be rewritten into a single file in lib/output

//go:generate go build "github.com/BytemarkHosting/bytemark-client/lib/brain" "github.com/BytemarkHosting/bytemark-client/lib/billing" "github.com/BytemarkHosting/bytemark-client/lib/spp"
//go:generate go run ../../gen/list_types.go -f "(*%s)(nil)," -o generated_types_test.go -t types_test.go.inc "github.com/BytemarkHosting/bytemark-client/lib/brain" "github.com/BytemarkHosting/bytemark-client/lib/billing" "github.com/BytemarkHosting/bytemark-client/lib/spp"
