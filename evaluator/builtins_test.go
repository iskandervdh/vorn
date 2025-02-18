package evaluator

import (
	"testing"

	"github.com/iskandervdh/vorn/object"
)

func TestLen(t *testing.T) {
	input := `len([1, 2, 3, 4])`

	testIntegerObject(t, testEval(input), 4)

	input = `len("hello")`

	testIntegerObject(t, testEval(input), 5)

	input = `len("hello" + " world")`

	testIntegerObject(t, testEval(input), 11)

	input = `len([])`

	testIntegerObject(t, testEval(input), 0)

	input = `len("")`

	testIntegerObject(t, testEval(input), 0)

	input = `len(1)`

	result := testEval(input)

	if result.Type() != object.ERROR_OBJ {
		t.Fatalf("expected ERROR_OBJ. got=%T", result)
	}

	if result.(*object.Error).Message != "argument to `len` not supported, got INTEGER" {
		t.Fatalf("expected 'argument to `len` not supported, got INTEGER'. got=%q", result.(*object.Error).Message)
	}

	input = `len([1, 2, 3], [4, 5, 6])`

	result = testEval(input)

	if result.Type() != object.ERROR_OBJ {
		t.Fatalf("expected ERROR_OBJ. got=%T", result)
	}

	if result.(*object.Error).Message != "wrong number of arguments. got 2, want 1" {
		t.Fatalf("expected 'wrong number of arguments. got 2, want 1'. got=%q", result.(*object.Error).Message)
	}
}

func TestFirst(t *testing.T) {
	input := `first([1, 2, 3, 4])`

	testIntegerObject(t, testEval(input), 1)

	input = `first([])`

	testNullObject(t, testEval(input))

	input = `first("hello")`

	testStringObject(t, testEval(input), "h")

	input = `first(1234)`

	testErrorObject(t, testEval(input), "argument to `first` must be ARRAY or STRING, got INTEGER")

	input = `first([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 2, want 1")
}

func TestLast(t *testing.T) {
	input := `last([1, 2, 3, 4])`

	testIntegerObject(t, testEval(input), 4)

	input = `last([])`

	testNullObject(t, testEval(input))

	input = `last("hello")`

	testStringObject(t, testEval(input), "o")

	input = `last(1234)`

	testErrorObject(t, testEval(input), "argument to `last` must be ARRAY or STRING, got INTEGER")

	input = `last([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 2, want 1")
}

func TestRest(t *testing.T) {
	input := `rest([1, 2, 3, 4])`

	testArrayObject(t, testEval(input), []string{"2", "3", "4"})

	input = `rest([])`

	testNullObject(t, testEval(input))

	input = `rest(1234)`

	testErrorObject(t, testEval(input), "argument to `rest` must be ARRAY, got INTEGER")

	input = `rest([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 2, want 1")
}

func TestPush(t *testing.T) {
	input := `push([1, 2, 3, 4], 5)`

	testArrayObject(t, testEval(input), []string{"1", "2", "3", "4", "5"})

	input = `push([1, 2, 3, 4], 5, 6)`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 3, want 2")

	input = `push(1, 2)`

	testErrorObject(t, testEval(input), "first argument to `push` must be ARRAY, got INTEGER")
}

func TestPop(t *testing.T) {
	input := `pop([1, 2, 3, 4])`

	testArrayObject(t, testEval(input), []string{"1", "2", "3"})

	input = `pop([])`

	testNullObject(t, testEval(input))

	input = `pop(1234)`

	testErrorObject(t, testEval(input), "first argument to `pop` must be ARRAY, got INTEGER")

	input = `pop([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 2, want 1")
}

func TestMap(t *testing.T) {
	input := `func timesTwo(x) {
	return x * 2;
}

map([1, 2, 3, 4], timesTwo);`

	testArrayObject(t, testEval(input), []string{"2", "4", "6", "8"})

	// Test with builtin function
	input = `map([1, 2, 3, 4], sqrt)`

	testArrayObject(t, testEval(input), []string{"1", "1.4142135623730951", "1.7320508075688772", "2"})

	input = `map([1, 2, 3, 4], 2)`

	testErrorObject(t, testEval(input), "second argument to `map` must be FUNCTION or BUILTIN, got INTEGER")

	input = `map([1, 2, 3, 4], sqrt, sqrt)`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 3, want 2")

	input = `map([1, 2, 3, 4])`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 1, want 2")

	input = `map(1, 2)`

	testErrorObject(t, testEval(input), "first argument to `map` must be ARRAY, got INTEGER")

	// TestIterMap
	e := New()

	testErrorObject(t, e.builtinIterMap(), "wrong number of arguments. got 0, want 2")
}

func TestReduce(t *testing.T) {
	input := `func add(x, y) {
	return x + y;
}

reduce([1, 2, 3, 4], 0, add);`

	testIntegerObject(t, testEval(input), 10)

	input = `reduce([1, 2, 3], 2, pow)`

	testIntegerObject(t, testEval(input), 64)

	input = `reduce([1, 2, 3, 4], 0, 2)`

	testErrorObject(t, testEval(input), "third argument to `reduce` must be FUNCTION or BUILTIN, got INTEGER")

	input = `reduce([1, 2, 3, 4], 0, sqrt, sqrt)`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 4, want 3")

	input = `reduce([1, 2, 3, 4], 0)`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 2, want 3")

	input = `reduce(1, 2, 3)`

	testErrorObject(t, testEval(input), "first argument to `reduce` must be ARRAY, got INTEGER")

	// TestIterReduce
	e := New()

	testErrorObject(t, e.builtinIterReduce(), "wrong number of arguments. got 0, want 3")
}

func TestPrint(t *testing.T) {
	input := `print("hello", "world")`

	testNullObject(t, testEval(input))

	input = `print("hello", "world", 1)`

	testNullObject(t, testEval(input))
}

func TestPow(t *testing.T) {
	input := `pow(2, 3)`

	testIntegerObject(t, testEval(input), 8)

	input = `pow(2, 0)`

	testIntegerObject(t, testEval(input), 1)

	input = `pow(2, -1)`

	testFloatObject(t, testEval(input), 0.5)

	input = `pow(-2, 3)`

	testIntegerObject(t, testEval(input), -8)

	input = `pow(-2, 0)`

	testIntegerObject(t, testEval(input), 1)

	input = `pow(-2, -1)`

	testFloatObject(t, testEval(input), -0.5)

	input = `pow(2, 3.0)`

	testFloatObject(t, testEval(input), 8.0)

	input = `pow(2.0, 3)`

	testFloatObject(t, testEval(input), 8.0)

	input = `pow(2.0, 3.0)`

	testFloatObject(t, testEval(input), 8.0)

	input = `pow(2)`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 1, want 2")

	input = `pow("test", 3)`

	testErrorObject(t, testEval(input), "arguments to `pow` must be INTEGER or FLOAT, got STRING and INTEGER")

	input = `pow(2, "test")`

	testErrorObject(t, testEval(input), "arguments to `pow` must be INTEGER or FLOAT, got INTEGER and STRING")
}

func TestSqrt(t *testing.T) {
	input := `sqrt(4)`

	testFloatObject(t, testEval(input), 2.0)

	input = `sqrt(0)`

	testFloatObject(t, testEval(input), 0.0)

	input = `sqrt(-1)`

	testErrorObject(t, testEval(input), "argument to `sqrt` must be non-negative, got -1")

	input = `sqrt(4.0)`

	testFloatObject(t, testEval(input), 2.0)

	input = `sqrt("test")`

	testErrorObject(t, testEval(input), "argument to `sqrt` must be INTEGER or FLOAT, got STRING")

	input = `sqrt(4, 5)`

	testErrorObject(t, testEval(input), "wrong number of arguments. got 2, want 1")
}
