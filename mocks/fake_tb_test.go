package mocks_test

import "fmt"

type fakeTestingT struct {
	failed  bool
	skipped bool
}

func (fakeT fakeTestingT) Error(args ...interface{}) {
	fakeT.failed = true
	fmt.Println(args...)
}

func (fakeT fakeTestingT) Errorf(format string, args ...interface{}) {
	fakeT.failed = true
	fmt.Printf(format, args...)
}

func (fakeT fakeTestingT) Fail() {
	fakeT.failed = true
}

func (fakeT fakeTestingT) FailNow() {
	fakeT.failed = true
}

func (fakeT fakeTestingT) Failed() bool {
	return fakeT.failed
}

func (fakeT fakeTestingT) Fatal(args ...interface{}) {
	fakeT.failed = true
	fmt.Println(args...)
}

func (fakeT fakeTestingT) Fatalf(format string, args ...interface{}) {
	fakeT.failed = true
	fmt.Printf(format, args...)
}

func (fakeT fakeTestingT) Log(args ...interface{}) {
	fmt.Println(args...)
}

func (fakeT fakeTestingT) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)

}

func (fakeT fakeTestingT) Name() string {
	return ""
}

func (fakeT fakeTestingT) private() {
	// THANKS GO
	return
}

func (fakeT fakeTestingT) Skip(args ...interface{}) {
	fakeT.skipped = true
}

func (fakeT fakeTestingT) SkipNow() {
	fakeT.skipped = true
}

func (fakeT fakeTestingT) Skipf(format string, args ...interface{}) {
	fakeT.skipped = true
	fmt.Printf(format, args...)
	return
}

func (fakeT fakeTestingT) Skipped() bool {
	return fakeT.skipped
}

func (fakeT fakeTestingT) Helper() {
	return
}
