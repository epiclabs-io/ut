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

package mypackage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/epiclabs-io/ut"
)

func TestConcurrent(tx *testing.T) {
	t := ut.BeginTest(tx, false)
	defer t.FinishTest()

	t.Go(func() {
		fmt.Println("lengthy process started...")
		time.Sleep(200 * time.Millisecond)
		// test some work that has to run in parallel
		t.Assert(1 == 1, "one should be equal to one!")
		fmt.Println("lengthy process finished...")

	})

	// you can also launch goroutines yourself, but you'll need to increment
	// the counter with t.RoutineStart() and call
	// t.RoutineEnd() when your routine ends, so the main routine can wait for it to finish

	t.RoutineStart()
	go func() {
		defer t.RoutineEnd()
		fmt.Println("second lengthy process started...")
		time.Sleep(300 * time.Millisecond)
		// test some work that has to run in parallel
		t.Assert(7 == 7, "seven should be equal to one!")
		//t.Fatal("crashed!") // you can call any test function inside a goroutine
		fmt.Println("second lengthy process finished...")
	}()

	fmt.Println("Some quick tests here...")
	// test some other things
	t.Assert(5 > 3, "5 should be greater than 3.")
	fmt.Println("finished quick part...")

}
