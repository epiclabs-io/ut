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
	"errors"
	"sync"
	"testing"

	"github.com/epiclabs-io/ut"
)

type metaTest struct {
	name    string
	early   bool
	failed  bool
	panic   interface{}
	handler func(t *ut.TestTools)
}

func MetaTester(name string, pseudoTest func(tt *ut.TestTools)) (ft *fakeT, early bool, panicValue interface{}) {
	ft = new(fakeT)
	ft.name = "TestTools"
	var wg sync.WaitGroup

	early = true
	go func() {
		defer func() {
			panicValue = recover()
			wg.Done()
		}()

		func() {
			tt := ut.ToolsBeginTest(ft, false)
			defer tt.FinishTest()
			pseudoTest(tt)
		}()

		early = false
	}()

	wg.Add(1)
	wg.Wait()
	return
}

var batch = []metaTest{
	{
		name: "empty",
		handler: func(tt *ut.TestTools) {
		},
	},
	{
		name:   "Error",
		early:  true,
		failed: true,
		handler: func(t *ut.TestTools) {
			t.Error(errors.New("some error"))
		},
	},
	{
		name: "Assert-ok",
		handler: func(t *ut.TestTools) {
			t.Assert(true, "some message")
		},
	},
	{
		name:   "Assert-fail",
		failed: true,
		early:  true,
		handler: func(t *ut.TestTools) {
			t.Assert(false, "some message")
		},
	},
	{
		name: "Ok-ok",
		handler: func(t *ut.TestTools) {
			t.Ok(nil)
		},
	},
	{
		name:   "Ok-fail",
		failed: true,
		early:  true,
		handler: func(t *ut.TestTools) {
			t.Ok(errors.New("an error!"))
		},
	},
	{
		name: "Equals-ok",
		handler: func(t *ut.TestTools) {
			t.Equals("hello", "hello")
		},
	},
	{
		name:   "Equals-fail",
		failed: true,
		early:  true,
		handler: func(t *ut.TestTools) {
			t.Equals("hello", "world")
		},
	},
	{
		name: "MustFail-ok",
		handler: func(t *ut.TestTools) {
			t.MustFail(errors.New("an error"), "should have failed")
		},
	},
	{
		name:   "MustFail-fail",
		failed: true,
		early:  true,
		handler: func(t *ut.TestTools) {
			t.MustFail(nil, "should have failed")
		},
	},
	{
		name: "MustFailWith-ok",
		handler: func(t *ut.TestTools) {
			ErrSomeError := errors.New("an error")
			t.MustFailWith(ErrSomeError, ErrSomeError)
		},
	},
	{
		name:   "MustFail-fail",
		failed: true,
		early:  true,
		handler: func(t *ut.TestTools) {
			ErrSomeError := errors.New("an error")
			t.MustFailWith(nil, ErrSomeError)
		},
	},
	{
		name:   "MustFail-fail",
		failed: true,
		early:  true,
		handler: func(t *ut.TestTools) {
			ErrSomeError := errors.New("an error")
			ErrOtherError := errors.New("some other error")
			t.MustFailWith(ErrOtherError, ErrSomeError)
		},
	},
}

func TestTestTools(t *testing.T) {
	for i, mt := range batch {

		ft, early, panic := MetaTester(mt.name, mt.handler)

		if early != mt.early {
			t.Fatalf("#%d, '%s': Expected early=%v, got %v", i, mt.name, mt.early, early)
		}

		if panic != nil && mt.panic == nil {
			t.Fatalf("#%d, '%s': Expected no panic got %v", i, mt.name, panic)
		}

		if ft.fail != mt.failed {
			t.Fatalf("#%d, '%s': Expected failed=%v, got %v", i, mt.name, mt.failed, ft.fail)
		}
	}
}
