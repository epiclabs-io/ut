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
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// T is a interface to *testing.T with the minimal
// set of methods needed
type T interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
}

// Assert fails the test if the condition is false.
func Assert(tb T, condition bool, msg string, v ...interface{}) {
	if Internal.NotAssert(0, condition, msg, v...) {
		tb.FailNow()
	}
}

// MustFail checks if err == nil. If so, it fails the test
func MustFail(tb T, err error, msg string, v ...interface{}) {
	if Internal.NotAssert(0, err != nil, msg, v...) {
		tb.FailNow()
	}
}

// MustFailWith checks if err equals an expected error. If not, it will fail the test.
func MustFailWith(tb T, err error, expectedError error) {
	if Internal.NotAssert(0, err == expectedError, fmt.Sprintf("Expected error to be '%s'. Got '%s'",
		Internal.ErrorString(expectedError), Internal.ErrorString(err))) {
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb T, err error) {
	if Internal.NotOk(0, err) {
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb T, expected, actual interface{}) {
	if Internal.NotEquals(0, expected, actual) {
		tb.FailNow()
	}
}

// JSONEquals fails if provided JSONs are not equivalent
func JSONEquals(tb T, expected, actual []byte) {
	if Internal.NotJSONEquals(0, expected, actual) {
		tb.FailNow()
	}
}

// JSONEqualsString performs a JSON comparison of the given object
// with the JSON contained in the referenced string
func JSONEqualsString(tb T, expected string, actual interface{}) {
	actualBytes, err := json.Marshal(actual)
	if err != nil {
		//tt.Fatalf("Cannot marshal 'actual' to JSON: %s", err)
		tb.FailNow()
	}
	if Internal.NotJSONEquals(0, []byte(expected), actualBytes) {
		tb.FailNow()
	}
}

// RandomArray returns a deterministically generated random array
// so values are the same across tests.
func RandomArray(i, length int) []byte {
	source := rand.NewSource(int64(i))
	r := rand.New(source)
	b := make([]byte, length)
	for n := 0; n < length; n++ {
		b[n] = byte(r.Intn(256))
	}
	return b
}

// CreateDirectory creates a directory and all necessary parents.
func CreateDirectory(path string) {

	_ = os.MkdirAll(path, 0750|os.ModeDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

}
