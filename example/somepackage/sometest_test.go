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
