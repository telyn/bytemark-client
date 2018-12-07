package main

import "flag"

// ok for this bad boy you gotta set the arguments as all the types you wanna
// generate slice flags for.
func main() {
	outputFile := flag.String("o", "-", "File to output to. Blank or - for stdin")
	templateFile := flag.String("t", "", "File to use as template for sprintf. if blank, just list the types and which ")
	fmtStr := flag.String("f", "%s", "Format string to use on each type before sending to the template")

}
