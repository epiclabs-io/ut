package testutils

import (
	"testing"

	"github.com/epiclabs-io/ut"
)

type ExampleTestTools struct {
	*ut.TestTools
	Services *ExampleServices
}

func BeginTest(tb testing.TB) *ExampleTestTools {
	ett := new(ExampleTestTools)
	ett.TestTools = ut.BeginTest(tb)
	ett.Services = NewExampleServices(ett)
	return ett
}
