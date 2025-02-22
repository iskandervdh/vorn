package evaluator

import (
	"testing"
)

func TestType(t *testing.T) {
	input := `type(1)`

	testStringObject(t, testEval(input), "INTEGER")

	input = `type(1.0)`

	testStringObject(t, testEval(input), "FLOAT")

	input = `type("hello")`

	testStringObject(t, testEval(input), "STRING")

	input = `type([1, 2, 3])`

	testStringObject(t, testEval(input), "ARRAY")

	input = `type({"key": "value"})`

	testStringObject(t, testEval(input), "HASH")

	input = `type(true)`

	testStringObject(t, testEval(input), "BOOLEAN")

	input = `type(false)`

	testStringObject(t, testEval(input), "BOOLEAN")

	input = `type(null)`

	testStringObject(t, testEval(input), "NULL")

	input = `type()`

	testErrorObject(t, testEval(input), "[1:6] wrong number of arguments. got 0, want 1")
}

func TestRange(t *testing.T) {
	input := `range(5)`

	testArrayObject(t, testEval(input), []string{"0", "1", "2", "3", "4"})

	input = `range(0)`

	testArrayObject(t, testEval(input), []string{})

	input = `range(5, 10)`

	testArrayObject(t, testEval(input), []string{"5", "6", "7", "8", "9"})

	input = `range(4, -3)`

	testArrayObject(t, testEval(input), []string{"4", "3", "2", "1", "0", "-1", "-2"})

	input = `range(4, 4)`

	testArrayObject(t, testEval(input), []string{})

	input = `range(-2, 3)`

	testArrayObject(t, testEval(input), []string{"-2", "-1", "0", "1", "2"})

	input = `range(-1)`

	testErrorObject(t, testEval(input), "[1:7] argument to `range` must be non-negative, got -1")

	input = `range("hello")`

	testErrorObject(t, testEval(input), "[1:7] first argument to `range` must be INTEGER, got STRING")

	input = `range()`

	testErrorObject(t, testEval(input), "[1:7] wrong number of arguments. got 0, want 1 or 2")

	input = `range(1, "test")`

	testErrorObject(t, testEval(input), "[1:7] second argument to `range` must be INTEGER, got STRING")
}

func TestInt(t *testing.T) {
	input := `int(1)`

	testIntegerObject(t, testEval(input), 1)

	input = `int(1.0)`

	testIntegerObject(t, testEval(input), 1)

	input = `int("1")`

	testIntegerObject(t, testEval(input), 1)

	input = `int("1.0")`

	testErrorObject(t, testEval(input), "[1:5] could not parse \"1.0\" as INTEGER")

	input = `int("hello")`

	testErrorObject(t, testEval(input), "[1:5] could not parse \"hello\" as INTEGER")

	input = `int()`

	testErrorObject(t, testEval(input), "[1:5] wrong number of arguments. got 0, want 1")

	input = `int([1])`

	testErrorObject(t, testEval(input), "[1:5] argument to `int` not supported, got ARRAY")
}

func TestFloat(t *testing.T) {
	input := `float(1)`

	testFloatObject(t, testEval(input), 1.0)

	input = `float(1.0)`

	testFloatObject(t, testEval(input), 1.0)

	input = `float("1")`

	testFloatObject(t, testEval(input), 1.0)

	input = `float("1.0")`

	testFloatObject(t, testEval(input), 1.0)

	input = `float("hello")`

	testErrorObject(t, testEval(input), "[1:7] could not parse \"hello\" as FLOAT")

	input = `float()`

	testErrorObject(t, testEval(input), "[1:7] wrong number of arguments. got 0, want 1")

	input = `float([1])`

	testErrorObject(t, testEval(input), "[1:7] argument to `float` not supported, got ARRAY")
}

func TestString(t *testing.T) {
	input := `string(1)`

	testStringObject(t, testEval(input), "1")

	input = `string(1.0)`

	testStringObject(t, testEval(input), "1")

	input = `string("1")`

	testStringObject(t, testEval(input), "1")

	input = `string("hello")`

	testStringObject(t, testEval(input), "hello")

	input = `string([1])`

	testStringObject(t, testEval(input), "[1]")

	input = `string({1: 1})`

	testStringObject(t, testEval(input), "{1: 1}")

	input = `string()`

	testErrorObject(t, testEval(input), "[1:8] wrong number of arguments. got 0, want 1")
}

func TestBool(t *testing.T) {
	input := `bool(true)`

	testBooleanObject(t, testEval(input), true)

	input = `bool(false)`

	testBooleanObject(t, testEval(input), false)

	input = `bool(1)`

	testBooleanObject(t, testEval(input), true)

	input = `bool(0)`

	testBooleanObject(t, testEval(input), false)

	input = `bool(1.0)`

	testBooleanObject(t, testEval(input), true)

	input = `bool(0.0)`

	testBooleanObject(t, testEval(input), false)

	input = `bool("true")`

	testBooleanObject(t, testEval(input), true)

	input = `bool("false")`

	testBooleanObject(t, testEval(input), true)

	input = `bool("hello")`

	testBooleanObject(t, testEval(input), true)

	input = `bool([1])`

	testBooleanObject(t, testEval(input), true)

	input = `bool([])`

	testBooleanObject(t, testEval(input), false)

	input = `bool({1: 1})`

	testBooleanObject(t, testEval(input), true)

	input = `bool({})`

	testBooleanObject(t, testEval(input), false)

	input = `bool(null)`

	testBooleanObject(t, testEval(input), false)

	input = `bool()`

	testErrorObject(t, testEval(input), "[1:6] wrong number of arguments. got 0, want 1")

	input = `bool(continue)`

	testErrorObject(t, testEval(input), "[1:6] argument to `bool` not supported, got CONTINUE")
}

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

	testErrorObject(t, result, "[1:5] argument to `len` not supported, got INTEGER")

	input = `len([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "[1:5] wrong number of arguments. got 2, want 1")
}

func TestFirst(t *testing.T) {
	input := `first([1, 2, 3, 4])`

	testIntegerObject(t, testEval(input), 1)

	input = `first([])`

	testNullObject(t, testEval(input))

	input = `first("hello")`

	testStringObject(t, testEval(input), "h")

	input = `first(1234)`

	testErrorObject(t, testEval(input), "[1:7] argument to `first` must be ARRAY or STRING, got INTEGER")

	input = `first([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "[1:7] wrong number of arguments. got 2, want 1")
}

func TestLast(t *testing.T) {
	input := `last([1, 2, 3, 4])`

	testIntegerObject(t, testEval(input), 4)

	input = `last([])`

	testNullObject(t, testEval(input))

	input = `last("hello")`

	testStringObject(t, testEval(input), "o")

	input = `last(1234)`

	testErrorObject(t, testEval(input), "[1:6] argument to `last` must be ARRAY or STRING, got INTEGER")

	input = `last([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "[1:6] wrong number of arguments. got 2, want 1")
}

func TestRest(t *testing.T) {
	input := `rest([1, 2, 3, 4])`

	testArrayObject(t, testEval(input), []string{"2", "3", "4"})

	input = `rest([])`

	testNullObject(t, testEval(input))

	input = `rest(1234)`

	testErrorObject(t, testEval(input), "[1:6] argument to `rest` must be ARRAY, got INTEGER")

	input = `rest([1, 2, 3], [4, 5, 6])`

	testErrorObject(t, testEval(input), "[1:6] wrong number of arguments. got 2, want 1")
}

func TestPrint(t *testing.T) {
	input := `print("hello", "world")`

	testNullObject(t, testEval(input))

	input = `print("hello", "world", 1)`

	testNullObject(t, testEval(input))
}

func TestAbs(t *testing.T) {
	input := `abs(1)`

	testIntegerObject(t, testEval(input), 1)

	input = `abs(-1)`

	testIntegerObject(t, testEval(input), 1)

	input = `abs(1.0)`

	testFloatObject(t, testEval(input), 1.0)

	input = `abs(-1.0)`

	testFloatObject(t, testEval(input), 1.0)

	input = `abs("test")`

	testErrorObject(t, testEval(input), "[1:5] argument to `abs` must be INTEGER or FLOAT, got STRING")

	input = `abs()`

	testErrorObject(t, testEval(input), "[1:5] wrong number of arguments. got 0, want 1")
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

	testErrorObject(t, testEval(input), "[1:5] wrong number of arguments. got 1, want 2")

	input = `pow("test", 3)`

	testErrorObject(t, testEval(input), "[1:5] arguments to `pow` must be INTEGER or FLOAT, got STRING and INTEGER")

	input = `pow(2, "test")`

	testErrorObject(t, testEval(input), "[1:5] arguments to `pow` must be INTEGER or FLOAT, got INTEGER and STRING")
}

func TestSqrt(t *testing.T) {
	input := `sqrt(4)`

	testFloatObject(t, testEval(input), 2.0)

	input = `sqrt(0)`

	testFloatObject(t, testEval(input), 0.0)

	input = `sqrt(-1)`

	testErrorObject(t, testEval(input), "[1:6] argument to `sqrt` must be non-negative, got -1")

	input = `sqrt(4.0)`

	testFloatObject(t, testEval(input), 2.0)

	input = `sqrt("test")`

	testErrorObject(t, testEval(input), "[1:6] argument to `sqrt` must be INTEGER or FLOAT, got STRING")

	input = `sqrt(4, 5)`

	testErrorObject(t, testEval(input), "[1:6] wrong number of arguments. got 2, want 1")
}
