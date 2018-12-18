// Copyright 2018 The ut/microtest Authors
// This file is part of ut/microtest library.
//
// This library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this library. If not, see <http://www.gnu.org/licenses/>.

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
