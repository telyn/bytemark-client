// slice_flags is a tool to instantly implement a flag.Value called
// <Something>SliceFlag, which contains a
// slice of <Something>Flags, an accessor called <Something>Slice, and a test to
// ensure that the SliceFlag parses and Preprocesses correctly if necessary.
//
// It's intended for use in cmd/bytemark/app/flags - see slice_flags.go in that
// package for the go:generate directives.
//
// example usage: to add a SliceFlag for some NewFlag which uses Preprocess and
// for which "blah" is a value value:
// go run ./gen/slice_flags/slice_flags.go -preprocesser -o new_flag.go
// --example-input "blah" New
//
// For the generated test it is assumed that your original Flag, when Set with
// some string will return the same string when String() is called. If this is
// not the case, we need to modify this code and template_test.go.tmpl to
// accept an example output as well as the example input.
//
// If, for whatever reason, it is useful to reuse this code with different
// templates, the -t and -tt options specify the template for the SliceFlag and
// the test for the SliceFlag, respectively.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type sliceFlag struct {
	TypeName     string
	Preprocess   bool
	ExampleInput string
}

func log(text string) {
	_, _ = fmt.Fprintln(os.Stderr, text)
}

// ok for this bad boy you gotta set the arguments as all the types you wanna
// generate slice flags for.
func main() {
	outputFile := flag.String("o", "-", "File to output to. Blank or - for stdin")
	testOutputFile := flag.String("ot", "", "Test file to output to. - for stdin. If unset, defaults to same base same path without extension as outputFile, with _test.go appended")
	templateFile := flag.String("t", "./gen/slice_flags/template.go.tmpl", "File to use as template")
	testTemplateFile := flag.String("tt", "./gen/slice_flags/template_test.go.tmpl", "File to use as template")
	exampleInput := flag.String("example-input", "", "Example input, used to ensure that the test works")
	preprocess := flag.Bool("preprocesser", false, "Whether Preprocessor is implemented on the single-version of these flags")

	flag.Parse()
	var err error

	if *testOutputFile == "" {
		name := strings.TrimSuffix(filepath.Base(*outputFile), ".go")
		*testOutputFile = filepath.Join(filepath.Dir(*outputFile), name+"_test.go")
	}

	args := flag.Args()
	if len(args) != 1 {
		log("usage: go run slice_flags -o <output file> -t <template file> [-preprocesser] [-example-input <example>] <type name>")
		log("")
		log("creates a new SomeTypeSliceFlag for a flag type names SomeTypeFlag.")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("To generate multiple slice flags, call this multiple times")
		os.Exit(1)
	}

	data := sliceFlag{
		TypeName:     args[0],
		Preprocess:   *preprocess,
		ExampleInput: *exampleInput,
	}
	err = writeTemplate(*outputFile, *templateFile, data)
	if err != nil {
		fmt.Printf("couldn't write %s: %s\n", *outputFile, err)
		os.Exit(1)
	}
	err = writeTemplate(*testOutputFile, *testTemplateFile, data)
	if err != nil {
		fmt.Printf("couldn't write %s: %s\n", *testOutputFile, err)
		os.Exit(1)
	}
}

func writeTemplate(outputFile, templateFile string, data sliceFlag) (err error) {
	fmt.Println("writeTemplate started")
	var outputWriter io.WriteCloser = os.Stdout
	if outputFile != "" && outputFile != "-" {
		outputWriter, err = os.Create(outputFile)
		if err != nil {
			return
		}
	}
	defer func() { _ = outputWriter.Close() }()

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return
	}

	inputReader, inputWriter := io.Pipe()
	betweenReader, betweenWriter := io.Pipe()
	gofmtPath, err := exec.LookPath("gofmt")
	if err != nil {
		return
	}
	importsPath, err := exec.LookPath("goimports")
	if err != nil {
		return
	}
	gofmt := exec.Cmd{
		Path:   gofmtPath,
		Args:   []string{"gofmt"},
		Stdin:  inputReader,
		Stdout: betweenWriter,
		Stderr: os.Stderr,
	}
	imports := exec.Cmd{
		Path:   importsPath,
		Args:   []string{"goimports"},
		Stdin:  betweenReader,
		Stdout: outputWriter,
		Stderr: os.Stderr,
	}
	err = gofmt.Start()
	if err != nil {
		return fmt.Errorf("Couldn't start gofmt: %s", err)
	}
	err = imports.Start()
	if err != nil {
		return fmt.Errorf("Couldn't start goimports: %s", err)
	}

	// errors are echoed rather than returned from here on out in order to
	// ensure resources get cleaned up properly, otherwise we get a deadlock
	fmt.Println("executing template")
	err = tmpl.Execute(inputWriter, data)
	if err != nil {
		fmt.Printf("template errored: %s\n", err)
	}

	fmt.Println("closing inputWriter")
	_ = inputWriter.Close()

	fmt.Println("waiting for gofmt to finish")
	fmtErr := gofmt.Wait()
	if fmtErr != nil {
		fmt.Printf("gofmt errored: %s\n", fmtErr)
	}
	_ = betweenWriter.Close()

	fmt.Println("waiting for goimports to finish")
	importsErr := imports.Wait()
	if importsErr != nil {
		fmt.Printf("goimports errored: %s\n", importsErr)
	}
	if err != nil || fmtErr != nil || importsErr != nil {
		err = errors.New("template/goimports/gofmt errored")
	}
	fmt.Println("returning from writeTemplate")
	return
}
