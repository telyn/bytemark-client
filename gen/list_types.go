package main

import (
	"flag"
	"fmt"
	"go/importer"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type SliceFlag []string

func (sf *SliceFlag) Set(value string) error {
	*sf = append(*sf, value)
	return nil
}

func (sf SliceFlag) String() string {
	return strings.Join(sf, ", ")
}

func (sf SliceFlag) Get() interface{} {
	return sf
}

func main() {
	var packageFlag SliceFlag
	flag.Var(&packageFlag, "p", "package to list types of. can be used multiple times")
	outputFile := flag.String("o", "-", "File to output to. Blank or - for stdin")
	templateFile := flag.String("t", "", "File to use as template for sprintf. if blank, just list the types")
	fmtStr := flag.String("f", "%s", "Format string to use on each type before sending to the template")

	flag.Parse()
	var err error

	tmpl := "%s"
	if *templateFile != "" {
		var bytes []byte
		bytes, err = ioutil.ReadFile(*templateFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tmpl = string(bytes)
	}

	var wr io.WriteCloser = os.Stdout
	if *outputFile != "" && *outputFile != "-" {
		wr, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	importer := importer.For("gc", nil)

	// aaaallllrighty that's all the flag stuff outta the way
	// now we read all the packages and fmt.Fprintf(wr, tmpl, types)
	var types []string

	for _, p := range packageFlag {
		pkg, err := importer.Import(p)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pkgName := pkg.Name()
		pkgPath := pkg.Path()
		scope := pkg.Scope()
		names := scope.Names()
		for _, name := range names {
			obj := scope.Lookup(name)
			inScopeRef := fmt.Sprintf("%s.%s", pkgName, name)
			fullNameWithPath := fmt.Sprintf("%s.%s", pkgPath, name)
			if obj.Type().String() == fullNameWithPath {
				types = append(types, fmt.Sprintf(*fmtStr, inScopeRef))
			}
		}

	}

	fmt.Fprintf(wr, tmpl, strings.Join(types, "\n"))
	err = wr.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
