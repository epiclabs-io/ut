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
