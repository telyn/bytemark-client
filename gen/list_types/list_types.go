package main

import (
	"flag"
	"fmt"
	"go/importer"
	gotypes "go/types"
	"os"
	"strings"
)

func main() {
	fmtStr := flag.String("f", "%s", "Format string to use on each type before sending to the template")
	tmpl, wr := setupGenerator()
	importer := importer.Default()

	// aaaallllrighty that's all the flag stuff outta the way
	// now we read all the packages and fmt.Fprintf(wr, tmpl, types)
	var types []string
	hasFailed := false

	for _, p := range flag.Args() {
		var pkg *gotypes.Package
		pkg, err := importer.Import(p)
		if err != nil {
			fmt.Println(err)
			hasFailed = true
			continue
		}
		pkgName := pkg.Name()
		pkgPath := pkg.Path()
		scope := pkg.Scope()
		names := scope.Names()
		for _, name := range names {
			obj := scope.Lookup(name)
			inScopeRef := fmt.Sprintf("%s.%s", pkgName, name)
			fullNameWithPath := fmt.Sprintf("%s.%s", pkgPath, name)
			if obj.Exported() && obj.Type().String() == fullNameWithPath {
				types = append(types, fmt.Sprintf(*fmtStr, inScopeRef))
			}
		}

	}

	_, err := fmt.Fprintf(wr, tmpl, strings.Join(types, "\n"))
	if err != nil {
		fmt.Println(err)
		hasFailed = true
	}
	err = wr.Close()
	if err != nil {
		fmt.Println(err)
		hasFailed = true
	}
	if hasFailed {
		os.Exit(1)
	}
}
