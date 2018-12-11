package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type sliceFlag struct {
	TypeName     string
	Preprocess   bool
	ExampleInput string
}

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
	writeTemplate(*testOutputFile, *testTemplateFile, data)
	if err != nil {
		fmt.Printf("couldn't write %s: %s\n", *testOutputFile, err)
		os.Exit(1)
	}
}

func writeTemplate(outputFile, templateFile string, data sliceFlag) (err error) {
	var outputWriter io.WriteCloser = os.Stdout
	if outputFile != "" && outputFile != "-" {
		outputWriter, err = os.Create(outputFile)
		if err != nil {
			return
		}
	}
	defer outputWriter.Close()

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return
	}

	rc, wc, _ := pipe.Commands(
		exec.Command("gofmt"),
		exec.Command("goimports"),
	)
	if err != nil {
		return
	}
	defer rc.Close()

	err = tmpl.Execute(wc, data)
	if err != nil {
		return
	}
	wc.Close()

	io.Copy(outputWriter, rc)
	return
}
