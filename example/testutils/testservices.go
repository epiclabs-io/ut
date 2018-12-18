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
