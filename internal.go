package ut

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/epiclabs-io/diff3"
)

func suspend() {
	c := make(chan bool)
	<-c
}

func assert(callDepth int, condition bool, msg string, v ...interface{}) bool {
	if !condition {
		_, file, line, _ := runtime.Caller(2 + callDepth)
		fmt.Printf("%s:%d: Assertion failed: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		return true
	}
	return false
}

func errorString(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}

func ok(callDepth int, err error) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(2 + callDepth)
		fmt.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		return true
	}
	return false
}

func equals(callDepth int, expected, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(2 + callDepth)
		fmt.Printf("%s:%d:\n\n\texpected: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		return true
	}
	return false
}

func jsonPretty(jsonBytes []byte) []byte {
	var buf bytes.Buffer
	json.Indent(&buf, jsonBytes, "", "\t")
	return buf.Bytes()
}

func jsonEquals(callDepth int, expected, actual []byte) bool {
	//credit for the trick: turtlemonvh https://gist.github.com/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal(expected, &o1)
	if err != nil {
		_, file, line, _ := runtime.Caller(2 + callDepth)
		fmt.Printf("%s:%d:\n\n\tJSONEquals: Error decoding 'expected' JSON: %s.\n\t Can't decode this: `%s`\n\n",
			filepath.Base(file), line, err, string(expected))
		return true
	}
	err = json.Unmarshal(actual, &o2)
	if err != nil {
		_, file, line, _ := runtime.Caller(2 + callDepth)
		fmt.Printf("%s:%d:\n\n\tJSONEquals: Error decoding 'actual' JSON: %s.\n\tCan't decode this: `%s`\n\n",
			filepath.Base(file), line, err, string(actual))
		return true
	}

	if !reflect.DeepEqual(o1, o2) {
		_, file, line, _ := runtime.Caller(2 + callDepth)
		expectedPretty := jsonPretty(expected)
		actualPretty := jsonPretty(actual)
		fmt.Printf("%s:%d:\n\n\texpected JSON: %s\n\n\tgot JSON: %s\n\n", filepath.Base(file), line, expectedPretty, actualPretty)
		r, err := diff3.Merge(bytes.NewReader(expectedPretty), bytes.NewReader([]byte{}), bytes.NewReader(actualPretty), true, "EXPECTED", "ACTUAL")
		if err == nil && r.Conflicts {
			diff, err := ioutil.ReadAll(r.Result)
			if err == nil {
				fmt.Printf("Diff:\n%s\n", string(diff))
			}
		}
		return true
	}
	return false
}

func fatal(callDepth int, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2 + callDepth)
	fmt.Println(append([]interface{}{fmt.Sprintf("%s:%d: FATAL:", filepath.Base(file), line)}, args...)...)
}

func fatalf(callDepth int, formatString string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2 + callDepth)
	fmt.Printf("%s: FATAL: "+formatString+"\n",
		append([]interface{}{fmt.Sprintf("%s:%d", filepath.Base(file), line)}, args...)...)
}
