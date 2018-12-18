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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/epiclabs-io/ut"
)

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
