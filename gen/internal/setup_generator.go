package internal

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// does all the tedious setup to allow writing templated stuff to a writer
func SetupGenerator() (tmpl string, wr io.WriteCloser) {
	outputFile := flag.String("o", "-", "File to output to. Blank or - for stdin")
	templateFile := flag.String("t", "", "File to use as template for sprintf. if blank, just list the types")

	flag.Parse()
	var err error

	tmpl = "%s"
	if *templateFile != "" {
		var bytes []byte
		bytes, err = ioutil.ReadFile(*templateFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tmpl = string(bytes)
	}

	wr = os.Stdout
	if *outputFile != "" && *outputFile != "-" {
		wr, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return
}
