package testutils

import (
	"time"

	"github.com/epiclabs-io/ut"
	"github.com/epiclabs-io/ut/example/somepackage"
)

type ExampleServices struct {
	tt *ExampleTestTools
	*ut.FileServices
}

func NewExampleServices(tt *ExampleTestTools) *ExampleServices {
	return &ExampleServices{
		tt:           tt,
		FileServices: ut.NewFileServices(tt.TestTools),
	}
}

func (es *ExampleServices) NewInterestingService(ID int) *somepackage.InterestingService {
	is := somepackage.NewInterestingService(ID, 1000*time.Second, es.NewTempFile())
	es.tt.AddService(is)
	return is
}
