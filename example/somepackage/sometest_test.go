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

package somepackage_test

import (
	"strings"
	"testing"

	"github.com/epiclabs-io/ut/example/testutils"
)

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
