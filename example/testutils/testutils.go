package testutils

import (
	"testing"

	"github.com/epiclabs-io/ut"
)

// ExampleTestTools defines test helper functions
// for your project
type ExampleTestTools struct {
	*ut.TestTools
	Services *exampleServices
}

// BeginTest starts the test
// set generateResults to true to save test results to files.
func BeginTest(tb testing.TB, generateResults bool) *ExampleTestTools {
	ett := new(ExampleTestTools)
	ett.TestTools = ut.ToolsBeginTest(tb, false)
	ett.Services = newExampleServices(ett)
	return ett
}

func (ett *ExampleTestTools) IsLongString(st string) {
	if len(st) < 100 {
		// use Internal methods to correctly print the line where the error took place,
		// otherwise the error message would always refer to this function, not very useful.
		ut.Internal.Fatalf(0, "Expected the string to be long, got string of length=%d", len(st))
		ett.FailNow()
	}
}
