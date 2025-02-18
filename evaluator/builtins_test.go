package evaluator

import "testing"

func TestMap(t *testing.T) {
	input := `func timesTwo(x) {
	return x * 2;
}

map([1, 2, 3, 4], timesTwo);`

	testArrayObject(t, testEval(input), []string{"2", "4", "6", "8"})

	// Test with builtin function
	input = `map([1, 2, 3, 4], sqrt)`

	testArrayObject(t, testEval(input), []string{"1", "1.4142135623730951", "1.7320508075688772", "2"})
}
