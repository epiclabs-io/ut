# &micro;t - microTest, a simple go testing framework

&micro;t is a minimalist golang testing framework designed so you can easily adapt it to your project, as opposed to forcing you to rewrite all your tests to conform to a particular structure. It has the following features:

* Toolbox of checking functions, such as Assert, Ok, MustFail, Equals...
* JSON marshalling out-of-the-box testing
* Clear error messages when values don't match, even showing a diff
* Test result storing and regeneration
* Multi-goroutine support
* Automated cleanup
* TestServices allow you to wrap third-party instances of services you need in your tests, such as temporary files, folders, database instances...

## Installing

To install &micro;t, simply grab it using `go get`:

```sh
$ go get github.com/epiclabs.io/ut
```

## Basic usage:

&micro;t supports two usage modes: basic and advanced. The basic one gives you access to the built-in checking functions only.

```go
package mypackage_test

import (
	"testing"
	. "github.com/epiclabs-io/ut"
	"github.com/epiclabs-io/ut/example/mypackage"
)

// Operation is a simple struct to demonstrate tests
type Operation struct {
	A   int     `json:"a"`
	B   int     `json:"b"`
	Sum int     `json:"sum"`
	Div float64 `json:"div"`
}

func TestBasic(t *testing.T) {
	// example struct:
	op := Operation{
		A: 10,
		B: 2,
	}
	var err error

	op.Sum = mypackage.Sum(op.A, op.B)
	Assert(t, op.Sum == 12, "Expected sum to equal 12, got %d", op.Sum)
	// another way:
	Equals(t, 12, op.Sum)

	op.Div, err = mypackage.Div(op.A, op.B)
	Ok(t, err)             // check there were no errors
	Equals(t, 5.0, op.Div) // check the correct value was returned

	_, err = mypackage.Div(op.A, 0)
	MustFail(t, err, "Expected div to fail since divisor is 0")
	MustFailWith(t, err, mypackage.ErrDivByZero) // you can also expect a specific error

	// Test JSON marshalling:
	expectedJSON := `{
		"a": 10,
		"b": 2,
		"sum": 12,
		"div": 5
	}`
	JSONEqualsString(t, expectedJSON, op)

	// note that the JSON keys do not have to be in the same order to be considered equal:
	expectedJSON = `{
		"div": 5,
		"b": 2,
		"sum": 12,
		"a": 10
	}`
	JSONEqualsString(t, expectedJSON, op)
}
```

## Advanced usage:

Advanced usage allows for automatically generating test results to files in the `testadata` directory. When `BeginTest()` is called with its second parameter set to `true`, all files will be automatically regenerated.

```go
package mypackage_test

import (
	"testing"
	"github.com/epiclabs-io/ut"
	"github.com/epiclabs-io/ut/example/mypackage"
)

// We can further reduce testing code and improve clarity by using advanced mode:
func TestAdvanced(tx *testing.T) {
	t := ut.BeginTest(tx, false) // set to true to regenerate test results
	defer t.FinishTest()

	// example struct:
	op := Operation{
		A: 10,
		B: 2,
	}
	var err error

	op.Sum = mypackage.Sum(op.A, op.B)
	t.EqualsKey("sum", op.Sum) // result will be stored in testdata/TestAdvanced/results.json, under the "sum" key.
	// that way it can be regenerated automatically if we change test cases or the behavior of the tested function

	op.Div, err = mypackage.Div(op.A, op.B)
	t.Ok(err) // check there were no errors
	t.EqualsKey("div", op.Div)

	_, err = mypackage.Div(op.A, 0)
	t.MustFail(err, "Expected div to fail since divisor is 0")
	t.MustFailWith(err, mypackage.ErrDivByZero) // you can also expect a specific error

	t.JSONEqualsFile("operation.json", op) // result will be stored in testdata/TestAdvanced/operation.json
	// that way it can be regenerated automatically if we change test cases or the behavior of the tested function
}
```

## Test Services

&micro;t includes the concept ot "test service". A Test service is a wrapper for some third-party functionality you need available during the test ,such as a throwaway database or a temporary folder that must be cleaned after the test ends. &micro;t comes with `FileServices` by default, which provides temporary files and folders that are automatically deleted once the test is finished.

```go
func TestServices(tx *testing.T) {
	t := ut.BeginTest(tx, false) // set to true to regenerate test results
	defer t.FinishTest()

	tempFileName := t.Services.NewTempFile() // request a temporary file
	err := ioutil.WriteFile(tempFileName, []byte("Some data"), 0666)
	t.Ok(err)

	tempDir := t.Services.NewTempDir() // request a temporary folder
	for i := 0; i < 5; i++ {
		err = ioutil.WriteFile(filepath.Join(tempDir, fmt.Sprintf("file-%d.txt", i)), []byte("some data!"), 0666)
		t.Ok(err)
	}
	// temp files and folders are cleaned up.
}
```

## Customizing for your project

One aspect that makes &micro;t powerful is how easy it is to customize and extend for your project. This enables you to add custom test functions and services that are unique to your project.

To extend &micro;t, do the following:

1. Create a `testutils` package. You probably already have in your project some package where you have put all sorts of helper functions to help you test. If you already have it, use it.
2. Add a two new files with this content:

**testutils.go:** Here we add custom testing functions
```go
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
```
**testservices.go:** Here we add custom test services
```go
type exampleServices struct {
	tt               *ExampleTestTools
	*ut.FileServices // embed file services so we get managed temp files and folders
}

func newExampleServices(tt *ExampleTestTools) *exampleServices {
	return &exampleServices{
		tt:           tt,
		FileServices: ut.NewFileServices(tt.TestTools),
	}
}
```
That's the basic. Now you can add custom testers like below:
```go
func (ett *ExampleTestTools) IsLongString(st string) {
	if len(st) < 100 {
		// use Internal methods to correctly print the line where the error took place,
		// otherwise the error message would always refer to this function, not very useful.
		ut.Internal.Fatalf(0, "Expected the string to be long, got string of length=%d", len(st))
		ett.FailNow()
	}
}
```

You can also add custom services:
```go
func (es *exampleServices) NewOpenFileForWriting() *os.File {
	file, err := os.Create(es.NewTempFile())
	es.tt.Ok(err)
	es.tt.AddService(file) // add service so .Close() is called when test ends.
	return file
}
```

Putting it all together:
```go
func TestCustom(tx *testing.T) {
	t := testutils.BeginTest(tx, false)
	defer t.FinishTest()

	longString := strings.Repeat("A", 101)
	t.IsLongString(longString) // assert it is a long string

	file := t.Services.NewOpenFileForWriting() // get a new file for writing stuff

	i, err := file.WriteString(longString)
	t.Ok(err)
	t.Assert(i == len(longString), "Expected to have written the entire string")
}
```
