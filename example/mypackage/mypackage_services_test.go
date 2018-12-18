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
