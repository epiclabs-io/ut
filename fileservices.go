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
