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
	"os"
	"time"

	"github.com/epiclabs-io/ut"
	"github.com/epiclabs-io/ut/example/somepackage"
)

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

func (es *exampleServices) NewInterestingService(ID int) *somepackage.InterestingService {
	is := somepackage.NewInterestingService(ID, 1000*time.Second, es.NewTempFile())
	es.tt.AddService(is)
	return is
}

func (es *exampleServices) NewOpenFileForWriting() *os.File {
	file, err := os.Create(es.NewTempFile())
	es.tt.Ok(err)
	es.tt.AddService(file) // add service so .Close() is called when test ends.
	return file
}
