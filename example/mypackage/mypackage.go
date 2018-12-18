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
