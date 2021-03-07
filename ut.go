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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/epiclabs-io/diff3"
)

// Service is an interface to define a test service that needs
// cleanup on finish
type Service interface {
	Close() error
}

// SubTest interface allows you to define a custom subtest
type SubTest interface {
	fmt.Stringer
}

// StringSubTest is a simple subtest that is a string
type StringSubTest string

func (ss *StringSubTest) String() string {
	return string(*ss)
}

// GENERATE_RESULTS is a global flag that haves all tests
// generate test results rather than verify them
var GENERATE_RESULTS bool

type results map[string]json.RawMessage

// TestTools is a replacement for testing.T that adds
// useful testing methods
type TestTools struct {
	T
	W               sync.WaitGroup
	err             chan error
	SubTest         SubTest
	TestdataDir     string
	generateResults bool
	Results         results
	services        []Service
}

// ToolsBeginTest takes a *testing.T and returns a replacement
// TestTools. Don't forget to defer t.FinishTest() to ensure
// cleanup
func ToolsBeginTest(t T, generateResults bool) *TestTools {
	_, file, _, _ := runtime.Caller(2)
	tt := &TestTools{
		err:             make(chan error, 20),
		T:               t,
		TestdataDir:     filepath.Join(filepath.Dir(file), "testdata", t.Name()),
		generateResults: GENERATE_RESULTS || generateResults,
	}
	tt.loadResults()
	return tt
}

// loadResults loads key-value results from the corresponding
// test folder /results.json file
func (tt *TestTools) loadResults() {
	if tt.generateResults {
		tt.Results = make(results)
	} else {
		path := filepath.Join(tt.TestdataDir, "results.json")
		resultBytes, err := ioutil.ReadFile(path)
		if err != nil {
			tt.Results = make(results)
			return
		}
		err = json.Unmarshal(resultBytes, &tt.Results)
		if err != nil {
			tt.Results = make(results)
			return
		}
	}
}

// Error stops the test with the given error
func (tt *TestTools) Error(err error) {
	if tt.SubTest != nil {
		fmt.Printf("Failed subtest: %s\n", tt.SubTest.String())
	}
	select {
	case tt.err <- err:
		runtime.Goexit()
	default:
		tt.T.FailNow()
	}
}

// Assert verifies if the condition is true. If not, it fails the test
func (tt *TestTools) Assert(condition bool, msg string, v ...interface{}) {
	if Internal.NotAssert(0, condition, msg, v...) {
		tt.Error(fmt.Errorf("Assertion failed: %s", msg))
	}
}

// Ok checks if there is no error. Otherwise it fails the test
func (tt *TestTools) Ok(err error) {
	if Internal.NotOk(0, err) {
		tt.Error(err)
	}
}

// Equals tests if both objects are "deeply equal", otherwise it fails the test
func (tt *TestTools) Equals(expected, actual interface{}) {
	if Internal.NotEquals(0, expected, actual) {
		tt.Error(errors.New("Expressions don't match"))
	}
}

func (tt *TestTools) equalsJSONBytes(callDepth int, name string, actual interface{}, read func() []byte, write func(data []byte)) {
	if tt.generateResults {
		actualBytes, err := json.MarshalIndent(actual, "", "\t")
		if err != nil {
			tt.Fatalf("Cannot marshal actual value to json to store as result: %s", err)
		}
		write(json.RawMessage(actualBytes))
		return
	}

	actualValue := reflect.ValueOf(actual)
	if actualValue.Kind() == reflect.Ptr {
		actualValue = actualValue.Elem()
		actual = actualValue.Interface()
	}

	expectedValuePtr := reflect.New(actualValue.Type())
	err := json.Unmarshal(read(), expectedValuePtr.Interface())
	if err != nil {
		tt.Fatalf("Cannot unmarshal result value in '%s'", name)
	}

	expected := expectedValuePtr.Elem().Interface()
	if Internal.NotEquals(callDepth+1, expected, actual) {
		tt.Error(fmt.Errorf("Expressions don't match. Check file '%s' in testdata/%s or key '%s' in testdata/%s/results.json", name, tt.T.Name(), name, tt.T.Name()))
	}

	return
}

func (tt *TestTools) equalsString(callDepth int, name string, actual string, read func() string, write func(data string)) {
	if tt.generateResults {
		write(actual)
		return
	}

	expected := read()
	if Internal.NotEquals(callDepth+1, expected, actual) {
		r, err := diff3.Merge(strings.NewReader(expected), strings.NewReader(""), strings.NewReader(actual), true, "EXPECTED", "ACTUAL")
		if err == nil && r.Conflicts {
			diff, err := ioutil.ReadAll(r.Result)
			if err == nil {
				fmt.Printf("Diff:\n%s\n", string(diff))
			}
		}
		tt.Error(fmt.Errorf("Expressions don't match. Check file '%s' in testdata/%s", name, tt.T.Name()))
	}

	return
}

// EqualsKey verifies if the passed "actual" value is equal to the value in the
// given key of the current test's results.json
func (tt *TestTools) EqualsKey(key string, actual interface{}) {
	if tt.Results == nil {
		tt.Fatalf("To use EqualsKey(), call LoadResults() first")
	}
	tt.equalsJSONBytes(0, fmt.Sprintf("key:%s", key), actual, func() []byte {
		expectedValueBytes, ok := tt.Results[key]
		if !ok {
			tt.Fatalf("Cannot find result key '%s'", key)
		}
		return expectedValueBytes

	}, func(data []byte) {
		tt.Results[key] = data
	})
}

// EqualsTextFile checks if the passed "actual" value is equivalent
// to the text contained in the indicated file in the current test's
// testadata folder
func (tt *TestTools) EqualsTextFile(file string, actual string) {
	path := filepath.Join(tt.TestdataDir, file)
	tt.equalsString(0, file, actual, func() string {
		expectedValueBytes, err := ioutil.ReadFile(path)
		if err != nil {
			tt.Fatalf("Cannot open test result file %s: %s", err)
		}
		return string(expectedValueBytes)

	}, func(data string) {
		CreateDirectory(tt.TestdataDir)
		err := ioutil.WriteFile(path, []byte(data), 0660)
		if err != nil {
			tt.Fatalf("Cannot write test result file %s : %s", path, err)
		}
	})
}

//JSONEquals checks if the passed values are JSON-equal, comparing values
// taking into account keys can be in different order, etc.
func (tt *TestTools) JSONEquals(expected, actual []byte) {
	if Internal.NotJSONEquals(0, expected, actual) {
		tt.Error(errors.New("JSONs don't match"))
	}
}

func (tt *TestTools) jsonEqualsFile(callDepth int, file string, actual []byte) {
	path := filepath.Join(tt.TestdataDir, file)
	if tt.generateResults {
		CreateDirectory(tt.TestdataDir)
		err := ioutil.WriteFile(path, Internal.JSONPretty(actual), 0660)
		if err != nil {
			tt.Fatalf("Cannot write test result file %s : %s", path, err)
		}
	} else {
		f, err := os.Open(path)
		if err != nil {
			tt.Fatalf("Cannot open test result file %s : %s", path, err)
		}
		expected, err := ioutil.ReadAll(f)
		if err != nil {
			tt.Fatalf("Cannot read test result file %s : %s", path, err)
		}
		if Internal.NotJSONEquals(callDepth+1, expected, actual) {
			tt.Error(fmt.Errorf("JSONs don't match. Test result file: %s", path))
		}
	}
}

// JSONBytesEqualsFile performs a JSON comparison of the provided JSON bytes
// with the JSON contained in the referenced file
func (tt *TestTools) JSONBytesEqualsFile(file string, actual []byte) {
	tt.jsonEqualsFile(0, file, actual)
}

// JSONEqualsFile performs a JSON comparison of the given object
// with the JSON contained in the referenced file
func (tt *TestTools) JSONEqualsFile(file string, actual interface{}) {
	actualBytes, err := json.Marshal(actual)
	if err != nil {
		tt.Fatalf("Cannot marshal 'actual' to JSON: %s", err)
	}
	tt.jsonEqualsFile(0, file, actualBytes)
}

// TestJSONMarshaller is a convenient tool to test JSON marshalling/unmarshalling
// matching the marshalled data to the referenced file.
func (tt *TestTools) TestJSONMarshaller(filename string, sample interface{}) {
	actual, err := json.Marshal(sample)
	if Internal.NotOk(0, err) {
		tt.Error(err)
	}
	sampleType := reflect.TypeOf(sample)
	if sampleType.Kind() == reflect.Ptr {
		sampleType = sampleType.Elem()
		sample = reflect.ValueOf(sample).Elem().Interface()
	}
	tt.jsonEqualsFile(0, filename, actual)
	recoveredPtr := reflect.New(sampleType)
	err = json.Unmarshal(actual, recoveredPtr.Interface())
	if Internal.NotOk(0, err) {
		tt.Error(err)
	}
	if Internal.NotEquals(0, sample, recoveredPtr.Elem().Interface()) {
		tt.Error(errors.New("Expressions don't match"))
	}
}

// Fatal will fail the test immediately with an error message
func (tt *TestTools) Fatal(args ...interface{}) {
	Internal.Fatal(0, args...)
	tt.Error(errors.New("Fatal error"))
}

// Fatalf will fail the test immediately with a formatted error message
func (tt *TestTools) Fatalf(formatString string, args ...interface{}) {
	Internal.Fatalf(0, formatString, args...)
	tt.Error(errors.New("Fatal error"))
}

// MustFail checks if err == nil. If so, it fails the test
func (tt *TestTools) MustFail(err error, msg string, v ...interface{}) {
	if Internal.NotAssert(0, err != nil, msg, v...) {
		tt.Error(fmt.Errorf("Should have failed: %s", msg))
	}
}

// MustFailWith checks if err equals an expected error. If not, it will fail the test.
func (tt *TestTools) MustFailWith(err error, expectedError error) {
	msg := fmt.Sprintf("Expected error to be '%s'. Got '%s'",
		Internal.ErrorString(expectedError),
		Internal.ErrorString(err))
	if Internal.NotAssert(0, err == expectedError, msg) {
		tt.Error(fmt.Errorf("Should have failed: %s", msg))
	}
}

func testPanic(f func()) (bool, interface{}) {
	didPanic := false
	var message interface{}
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()
		f()

	}()
	return didPanic, message
}

// MustPanic runs a function and checks that it panics
func (tt *TestTools) MustPanic(f func()) {
	didPanic, _ := testPanic(f)
	msg := "Expected function to panic"
	if Internal.NotAssert(0, didPanic, msg) {
		tt.Error(errors.New("should have panicked"))
	}
}

// MustPanicWith runs a function and checks that it panics throwing a specific value
func (tt *TestTools) MustPanicWith(expectedMessage interface{}, f func()) {
	didPanic, recoveredMessage := testPanic(f)
	msg := "Expected function to panic"
	if Internal.NotAssert(0, didPanic, msg) {
		tt.Error(errors.New("should have panicked"))
	}
	if Internal.NotEquals(0, expectedMessage, recoveredMessage) {
		tt.Error(fmt.Errorf("Should have panicked with message: %v", expectedMessage))
	}
}

// AddService adds a service that will be cleaned up when the test ends for any reason.
func (tt *TestTools) AddService(s Service) {
	tt.services = append(tt.services, s)
}

// Go will launch a test goroutine
func (tt *TestTools) Go(subroutine func()) {
	tt.RoutineStart()
	go func() {
		defer tt.RoutineEnd()
		subroutine()
	}()
}

// RoutineStart adds one to the goroutine waiting counter
func (tt *TestTools) RoutineStart() {
	tt.W.Add(1)
}

// RoutineEnd notifies a child goroutine is finished
func (tt *TestTools) RoutineEnd() {
	tt.W.Done()
}

//StartSubTest marks the beginning of a new subtest
func (tt *TestTools) StartSubTest(fmtString string, args ...interface{}) {
	s := StringSubTest(fmt.Sprintf(fmtString, args...))
	tt.SubTest = &s
}

// EndSubTest indicates a subtest has finished
func (tt *TestTools) EndSubTest() {
	tt.SubTest = nil
}

// FinishTest waits for all test goroutines and cleans up
func (tt *TestTools) FinishTest() {
	tt.W.Wait()
	close(tt.err)

	for i := len(tt.services) - 1; i >= 0; i-- {
		err := tt.services[i].Close()
		if err != nil {
			fmt.Printf("Error closing service: %s\n", err)
		}
	}
	tt.services = nil

	e := recover()
	if e != nil {
		panic(e)
	}

	var errorCount int
	for err := range tt.err {
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			errorCount++
		}
	}
	if errorCount > 0 {
		if tt.SubTest != nil {
			fmt.Printf("Failed subtest: %s\n", tt.SubTest.String())
		}
		fmt.Printf("%d errors\n", errorCount)
		tt.T.FailNow()
	}
	if tt.generateResults {
		if tt.Results != nil && len(tt.Results) > 0 {
			resultsBytes, err := json.MarshalIndent(tt.Results, "", "\t")
			if err != nil {
				tt.T.Fatal("Cannot marshal results to JSON")
			}
			CreateDirectory(tt.TestdataDir)
			err = ioutil.WriteFile(filepath.Join(tt.TestdataDir, "results.json"), resultsBytes, 0660)
			if err != nil {
				tt.T.Fatal("Cannot write results.json")
			}
		}
		tt.T.Fatal("\n!!!!!\nTest actually passed :-), but GENERATE_RESULTS is activated. Set to false before committing!\n!!!!!\n")
	}
}
