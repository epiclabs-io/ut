package marshaller_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/epiclabs-io/ut"
)

type TestStruct struct {
	A string
}

func (ts *TestStruct) MarshalJSON() ([]byte, error) {
	if ts.A == "" {
		return nil, errors.New("A cannot be empty")
	}
	return json.Marshal([]byte(ts.A))
}

func (ts *TestStruct) UnmarshalJSON(data []byte) error {
	var b []byte
	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	ts.A = string(b)
	return nil
}

func TestCustomMarshaller(tx *testing.T) {
	t := ut.BeginTest(tx, false)
	defer t.FinishTest()

	x := &TestStruct{
		A: "hello",
	}

	t.TestJSONMarshaller("x.json", x)

}
