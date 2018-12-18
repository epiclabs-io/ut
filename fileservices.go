package ut

import (
	"fmt"
	"os"

	"io/ioutil"
	"path/filepath"
)

type FileServices struct {
	t             *TestTools
	tempFileCount int
	tempFilesDir  string
}

type TempDir struct {
	tempFileDir string
}

func (td *TempDir) Close() error {
	os.RemoveAll(td.tempFileDir)
	return nil
}

func NewFileServices(t *TestTools) *FileServices {
	return &FileServices{t: t}
}

func (fs *FileServices) NewTempDir() string {
	dir, err := ioutil.TempDir("", fs.t.Name())
	fs.t.Ok(err)
	fs.t.AddService(&TempDir{
		tempFileDir: dir,
	})
	return dir
}

func (fs *FileServices) NewTempFile() string {
	if fs.tempFilesDir == "" {
		fs.tempFilesDir = fs.NewTempDir()
	}
	defer func() { fs.tempFileCount++ }()
	return filepath.Join(fs.tempFilesDir, fmt.Sprintf("tempfile-%d", fs.tempFileCount))
}
