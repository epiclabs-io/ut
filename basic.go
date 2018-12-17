package ut

import (
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
	if assert(0, condition, msg, v...) {
		tb.FailNow()
	}
}

// MustFail checks if err == nil. If so, it fails the test
func MustFail(tb T, err error, msg string, v ...interface{}) {
	if assert(0, err != nil, msg, v...) {
		tb.FailNow()
	}
}

// MustFailWith checks if err equals an expected error. If not, it will fail the test.
func MustFailWith(tb T, err error, expectedError error) {
	if assert(0, err == expectedError, fmt.Sprintf("Expected error to be '%s'. Got '%s'", errorString(expectedError), errorString(err))) {
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb T, err error) {
	if ok(0, err) {
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb T, expected, actual interface{}) {
	if equals(0, expected, actual) {
		tb.FailNow()
	}
}

// JSONEquals fails if provided JSONs are not equivalent
func JSONEquals(tb T, expected, actual []byte) {
	if jsonEquals(0, expected, actual) {
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
