package morestrings_test

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib/output/morestrings"
)

func ExampleJoinWithSpecialLast_oneString() {
	strs := []string{
		"hi",
	}
	fmt.Print(morestrings.JoinWithSpecialLast(", ", " and ", strs))

	// Output: hi
}
func ExampleJoinWithSpecialLast_twoStrings() {
	strs := []string{
		"hi", "hello",
	}
	fmt.Print(morestrings.JoinWithSpecialLast(", ", " and ", strs))

	// Output: hi and hello
}
func ExampleJoinWithSpecialLast_threeStrings() {
	strs := []string{
		"hi", "hello", "welcome",
	}
	fmt.Print(morestrings.JoinWithSpecialLast(", ", " and ", strs))

	// Output: hi, hello and welcome
}
func ExampleJoinWithSpecialLast_long() {
	strs := []string{
		"hi", "hello", "welcome", "good evening",
	}
	fmt.Print(morestrings.JoinWithSpecialLast(", ", " and ", strs))

	// Output: hi, hello, welcome and good evening
}
