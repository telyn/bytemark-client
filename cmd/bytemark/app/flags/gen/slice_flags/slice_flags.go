package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ok for this bad boy you gotta set the arguments as all the types you wanna
// generate slice flags for.
func main() {
	outputFile := flag.String("o", "-", "File to output to. Blank or - for stdin")
	testOutputFile := flag.String("ot", "", "Test file to output to. - for stdin. If unset, defaults to same base same path without extension as outputFile, with _test.go appended")
	templateFile := flag.String("t", "./gen/slice_flags/template.go.tmpl", "File to use as template")
	testTemplateFile := flag.String("tt", "./gen/slice_flags/template_test.go.tmpl", "File to use as template")
	exampleInput := flag.String("example-input", "", "Example input")
	preprocess := flag.Bool("preprocesser", false, "Whether Preprocessor is implemented on the single-version of these flags")

	flag.Parse()
	var err error

	if *testOutputFile == "" {
		name := strings.TrimSuffix(filepath.Base(*outputFile), ".go")
		*testOutputFile = filepath.Join(filepath.Dir(*outputFile), name+"_test.go")
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("usage: go run slice_flags -o <output file> -t <template file> [-p] <type name>")
		fmt.Println("")
		fmt.Println("-p: implement Preprocesser by calling Preprocess on each value in the sliceflag slice.")
		fmt.Println("")
		fmt.Println("To generate multiple slice flags, call this multiple times")
		os.Exit(1)
	}

	data := struct {
		TypeName     string
		Preprocess   bool
		ExampleInput string
	}{
		TypeName:     args[0],
		Preprocess:   *preprocess,
		ExampleInput: *exampleInput,
	}

	var wr io.WriteCloser = os.Stdout
	if *outputFile != "" && *outputFile != "-" {
		wr, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	tmpl, err := template.ParseFiles(*templateFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = tmpl.Execute(wr, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wr = os.Stdout
	if *testOutputFile != "" && *testOutputFile != "-" {
		wr, err = os.Create(*testOutputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	tmpl, err = template.ParseFiles(*testTemplateFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = tmpl.Execute(wr, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
