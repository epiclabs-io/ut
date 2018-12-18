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

package ut_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/epiclabs-io/ut"
)

type fakeT struct {
	fail    bool
	stopped bool
	name    string
}

func (t *fakeT) Error(args ...interface{}) {
	t.Log(args...)
	t.Fail()
}
func (t *fakeT) Errorf(format string, args ...interface{}) {
	t.Logf(format, args...)
	t.Fail()
}
func (t *fakeT) Fail() {
	t.fail = true
}
func (t *fakeT) FailNow() {
	t.Fail()
	t.stopped = true
}
func (t *fakeT) Failed() bool {
	return t.fail
}
func (t *fakeT) Fatal(args ...interface{}) {
	t.Log(args...)
	t.FailNow()
}
func (t *fakeT) Fatalf(format string, args ...interface{}) {
	t.Log(fmt.Sprintln(args...))
	t.FailNow()
}
func (t *fakeT) Log(args ...interface{})                 {}
func (t *fakeT) Logf(format string, args ...interface{}) {}
func (t *fakeT) Name() string {
	if t.name == "" {
		return "testname"
	} else {
		return t.name
	}
}

func TestAssert(t *testing.T) {
	ft := new(fakeT)

	ut.Assert(ft, false, "")

	if ft.fail != true {
		t.Fatalf("Expected failed Assert to mark test as failed")
	}

	ft = new(fakeT)
	ut.Assert(ft, true, "")
	if ft.fail != false {
		t.Fatalf("Expected valid Assert to not mark the test as failed")
	}

}

func TestMustFail(t *testing.T) {
	ft := new(fakeT)

	ut.MustFail(ft, nil, "")

	if ft.fail != true {
		t.Fatalf("Expected MustFail to mark test as failed")
	}

	ft = new(fakeT)
	ut.MustFail(ft, errors.New("Some error"), "")
	if ft.fail != false {
		t.Fatalf("Expected MustFail with error to not mark the test as failed")
	}

}

func TestMustFailWith(t *testing.T) {
	ft := new(fakeT)
	SomeError := errors.New("SomeError")
	OtherError := errors.New("OtherError")

	ut.MustFailWith(ft, nil, SomeError)

	if ft.fail != true {
		t.Fatalf("Expected MustFailWith to mark test as failed")
	}

	ft = new(fakeT)
	ut.MustFailWith(ft, SomeError, SomeError)
	if ft.fail != false {
		t.Fatalf("Expected MustFail with error to not mark the test as failed")
	}

	ft = new(fakeT)
	ut.MustFailWith(ft, SomeError, OtherError)
	if ft.fail != true {
		t.Fatalf("Expected MustFail to mark the test as failed, since errors don't match")
	}

}

func TestOk(t *testing.T) {
	ft := new(fakeT)

	ut.Ok(ft, errors.New("error!!"))

	if ft.fail != true {
		t.Fatalf("Expected Ok to mark test as failed")
	}

	ft = new(fakeT)
	ut.Ok(ft, nil)
	if ft.fail != false {
		t.Fatalf("Expected no error to not mark the test as failed")
	}

}

func TestEquals(t *testing.T) {
	ft := new(fakeT)

	ut.Equals(ft, 5, 6)

	if ft.fail != true {
		t.Fatalf("Expected Equals to mark test as failed")
	}

	ft = new(fakeT)
	ut.Equals(ft, 5, 5)
	if ft.fail != false {
		t.Fatalf("Expected equal expressions to not mark the test as failed")
	}
}

func TestJSONEquals(t *testing.T) {
	ft := new(fakeT)

	json1 := []byte(`{
		"a": 6,
		"b": 5
	}
	`)

	json2 := []byte(`{
		"b": 5,
		"a": 6
	}
	`)

	json3 := []byte(`{
		"b": 5
	}
	`)

	ut.JSONEquals(ft, json1, json3)

	if ft.fail != true {
		t.Fatalf("Expected JSONEquals to mark test as failed")
	}

	ft = new(fakeT)
	ut.JSONEquals(ft, json1, json2)
	if ft.fail != false {
		t.Fatalf("Expected equal json expressions to not mark the test as failed")
	}
}

func TestRandomArray(t *testing.T) {
	var data [][]byte
	for i := 0; i < 100; i++ {
		data = append(data, ut.RandomArray(i, i))
	}

	for i := 0; i < 100; i++ {
		if !bytes.Equal(data[i], ut.RandomArray(i, i)) {
			t.Fatalf("Expected deterministic array to be the same")
		}
	}

}
