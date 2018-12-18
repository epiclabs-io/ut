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
