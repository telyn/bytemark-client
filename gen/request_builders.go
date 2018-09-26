package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/BytemarkHosting/bytemark-client/gen/internal"
)

func findClientInterface(f *ast.File) (clientInterface *ast.InterfaceType) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "Client" {
				if iface, ok := x.Type.(*ast.InterfaceType); ok {
					clientInterface = iface
				} else {
					fmt.Printf("ARGH! was hoping for an *ast.InterfaceType but got a %T????\nprobably check syntax of interface.go\n", x.Type)
				}
				return false
			}
		}
		return true
	})
	return
}

func findRequestTypes(clientInterface *ast.InterfaceType) []string {
	nameRegex := regexp.MustCompile("^Build(.+)Request")
	requestTypes := make([]string, 0)
	for _, method := range clientInterface.Methods.List {
		if len(method.Names) < 1 {
			continue
		}
		name := method.Names[0].Name
		if nameRegex.MatchString(name) {
			name = nameRegex.ReplaceAllString(name, "$1")
			requestTypes = append(requestTypes, name)
		}
	}
	return requestTypes
}

// lowerFirst lowercases the first character in str
func lowerFirst(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// takes a template and a CamelCase name and replaces the various cases of
// "virtual machine" in it with corresponding cases of camelName
func formatTmpl(tmpl string, camelName string) (out string) {
	snakeName := internal.ToSnake(camelName)
	englishName := strings.Replace(snakeName, "_", " ", 0)
	lowerCamelName := lowerFirst(camelName)

	englishRegexp := regexp.MustCompile("virtual machine")
	snakeRegexp := regexp.MustCompile("virtual_machine")
	camelRegexp := regexp.MustCompile("VirtualMachine")
	lowerCamelRegexp := regexp.MustCompile("virtualMachine")

	out = tmpl
	out = englishRegexp.ReplaceAllString(out, englishName)
	out = snakeRegexp.ReplaceAllString(out, snakeName)
	out = camelRegexp.ReplaceAllString(out, camelName)
	out = lowerCamelRegexp.ReplaceAllString(out, lowerCamelName)
	return
}

func main() {
	tmpl, wr := internal.SetupGenerator()

	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "interface.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	clientInterface := findClientInterface(f)

	names := findRequestTypes(clientInterface)

	out := []string{
		"package lib",
		"import \"github.com/BytemarkHosting/bytemark-client/lib/brain\"",
		"",
	}

	for _, name := range names {
		out = append(out, formatTmpl(tmpl, name))
		out = append(out, "")
	}

	hasFailed := false

	_, err = fmt.Fprint(wr, strings.Join(out, "\n"))
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
