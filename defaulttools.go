package ut

import (
	"testing"
)

type DefaultServices struct {
	*FileServices
}

type DefaultTestTools struct {
	*TestTools
	Services *DefaultServices
}

func BeginTest(tb testing.TB, generateResults bool) *DefaultTestTools {
	dtt := new(DefaultTestTools)
	dtt.TestTools = ToolsBeginTest(tb, generateResults)
	dtt.Services = new(DefaultServices)
	dtt.Services.FileServices = new(FileServices)
	dtt.Services.FileServices.t = dtt.TestTools
	return dtt
}
