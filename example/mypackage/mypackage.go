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

package mypackage

import "errors"

// Sum adds to integers and returns the result.
func Sum(a, b int) int {
	return a + b
}

var ErrDivByZero = errors.New("Division by zero!")

func Div(a, b int) (float64, error) {
	if b == 0 {
		return 0, ErrDivByZero
	}
	return float64(a / b), nil
}
