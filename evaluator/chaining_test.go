package evaluator

import (
	"testing"

	"github.com/iskandervdh/vorn/object"
)

func checkChainingExpression(t *testing.T, input string, expected interface{}) {
	evaluated := testEval(input)

	switch evaluated.Type() {
	case object.STRING_OBJ:
		testStringObject(t, evaluated, expected.(string))
	case object.INTEGER_OBJ:
		testIntegerObject(t, evaluated, int64(expected.(int)))
	case object.ERROR_OBJ:
		testErrorObject(t, evaluated, expected.(string))
	case object.ARRAY_OBJ:
		testArrayObject(t, evaluated, expected.([]string))
	case object.BOOLEAN_OBJ:
		testBooleanObject(t, evaluated, expected.(bool))
	default:
		t.Errorf("Expected STRING_OBJ, INTEGER_OBJ, ERROR_OBJ, ARRAY_OBJ or BOOLEAN_OBJ, got %s", evaluated.Type())
	}
}

func TestStringChainingExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// String
		{`"hello".upper()`, "HELLO"},
		{`"hello".upper().lower()`, "hello"},
		{`"hello".upper().lower().upper()`, "HELLO"},
		{`("hello" + "world").upper()`, "HELLOWORLD"},
		{`"hElLo".lower() + "world".lower().upper()`, "helloWORLD"},
		{`"hello".length()`, 5},
		{`"hello".split()`, []string{"hello"}},
		{`"hello".split("e")`, []string{"h", "llo"}},
		{`"hello".split("l")`, []string{"he", "", "o"}},
		{`"hello world".split(" ")`, []string{"hello", "world"}},
		{`"hello world".split("")`, []string{"h", "e", "l", "l", "o", " ", "w", "o", "r", "l", "d"}},
		{`"hello world".split("o")`, []string{"hell", " w", "rld"}},
		{`"hello world".split()`, []string{"hello", "world"}},
		{`"hello world".contains("world")`, true},
		{`"hello world".contains("worlds")`, false},
		{`"hello world".replace("world", "you")`, "hello you"},
		{`"hello world".replace("world", "you").replace("you", "world")`, "hello world"},
		{`"hello  ".trim()`, "hello"},
		{`"  hello".trim()`, "hello"},
		{`"  hello	".trim()`, "hello"},
		{`"hello".trim()`, "hello"},
		{`"  hello".trimStart()`, "hello"},
		{`"  hello	".trimStart()`, "hello	"},
		{`"hello".trimStart()`, "hello"},
		{`"hello  ".trimEnd()`, "hello"},
		{`"hello	".trimEnd()`, "hello"},
		{`"hello".trimEnd()`, "hello"},
		{`"hello".repeat(3)`, "hellohellohello"},
		{`"hello".repeat(0)`, ""},
		{`"hello".repeat(-1)`, ""},
		{`"hello".repeat(1)`, "hello"},
		{`"hello".reverse()`, "olleh"},
		{`"hello".reverse().reverse()`, "hello"},
		{`"hello".slice(1)`, "ello"},
		{`"hello".slice(1, 3)`, "el"},
		{`"hello".slice(1, 1)`, ""},
		{`"hello".slice(1, -1)`, "ell"},
		{`"hello".startsWith("he")`, true},
		{`"hello".startsWith("lo")`, false},
		{`"hello".endsWith("lo")`, true},
		{`"hello".endsWith("he")`, false},

		// Errors
		{`"hello".length(1)`, "[1:2] String.length() takes no arguments"},
		{`"hello".upper(1)`, "[1:2] String.upper() takes no arguments"},
		{`"hello".lower(1)`, "[1:2] String.lower() takes no arguments"},
		{`"hello".append(1)`, "[1:10] String has no method append"},
		{`"hello".split(1)`, "[1:2] argument to `split` must be STRING, got INTEGER"},
		{`"hello".split("e", "l")`, "[1:2] String.split() takes at most 1 argument, got 2"},
		{`"hello world".contains("world", 6)`, "[1:2] String.contains() takes exactly 1 argument"},
		{`"hello world".contains(6)`, "[1:2] argument to `contains` must be STRING, got INTEGER"},
		{`"hello world".replace("world", "you", "me")`, "[1:2] String.replace() takes exactly 2 arguments"},
		{`"hello world".replace(1, 2)`, "[1:2] first argument to `replace` must be STRING, got INTEGER"},
		{`"hello world".replace("world", 2)`, "[1:2] second argument to `replace` must be STRING, got INTEGER"},
		{`"hello".trim(1)`, "[1:2] String.trim() takes no arguments"},
		{`"hello".trimStart(1)`, "[1:2] String.trimStart() takes no arguments"},
		{`"hello".trimEnd(1)`, "[1:2] String.trimEnd() takes no arguments"},
		{`"hello".repeat()`, "[1:2] String.repeat() takes exactly 1 argument"},
		{`"hello".repeat("1")`, "[1:2] argument to `repeat` must be INTEGER, got STRING"},
		{`"hello".reverse(1)`, "[1:2] String.reverse() takes no arguments"},
		{`"hello".slice()`, "[1:2] String.slice() takes 1 or 2 arguments"},
		{`"hello".slice(1, 2, 3)`, "[1:2] String.slice() takes 1 or 2 arguments"},
		{`"hello".slice("1", 2)`, "[1:2] first argument to `slice` must be INTEGER, got STRING"},
		{`"hello".slice(1, "2")`, "[1:2] second argument to `slice` must be INTEGER, got STRING"},
		{`"hello".slice(10, 10)`, "[1:2] first argument to `slice` out of range"},
		{`"hello".slice(0, 10)`, "[1:2] second argument to `slice` out of range"},
		{`"hello".startsWith(1)`, "[1:2] argument to `startsWith` must be STRING, got INTEGER"},
		{`"hello".startsWith("1", 2)`, "[1:2] String.startsWith() takes exactly 1 argument"},
		{`"hello".endsWith(1)`, "[1:2] argument to `endsWith` must be STRING, got INTEGER"},
		{`"hello".endsWith("1", 2)`, "[1:2] String.endsWith() takes exactly 1 argument"},
	}

	for _, test := range tests {
		checkChainingExpression(t, test.input, test.expected)
	}
}

func TestArrayChainingExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1,2,3].length()`, 3},
		{`[1,2,3].append(4).length()`, 4},
		{`let a = [1,2,3]; a.append(4); a.length()`, 4},
		{`[1,2,3].pop()`, 3},
		{`let a = [1,2,3]; a.pop(); a.pop(); a`, []string{"1"}},
		{`let a = [1,2,3]; a.pop(1)`, 2},
		{`func timesTwo(x) { x * 2 }; [1, 2, 3, 4].map(timesTwo)`, []string{"2", "4", "6", "8"}},
		{`[1, 2, 3, 4].map(sqrt)`, []string{"1", "1.4142135623730951", "1.7320508075688772", "2"}},
		{`[1, 2, 3, 4].map(func(x, i) { return x + i; })`, []string{"1", "3", "5", "7"}},
		{`[1, 2, 3, 4].filter(func(x) { return x > 2; })`, []string{"3", "4"}},
		{`[1, 2, 3, 4].filter(func(x) { return 10; })`, []string{"1", "2", "3", "4"}},
		{`[1, 2, 3, 4].reduce(func(x, y) { return x + y; }, 0)`, 10},
		{`[1, 2, 3, 4].reduce(func(x, y, i) { return x + y + i; }, 0)`, 16},

		// Errors
		{`[1,2,3].length(1)`, "[1:2] Array.length() takes no arguments"},
		{`[1,2,3].append()`, "[1:2] Array.append() takes exactly 1 argument"},
		{`[1,2,3].upper()`, "[1:10] Array has no method upper"},
		{`[].pop()`, "[1:2] Array.pop() called on empty array"},
		{`[1,2,3].pop(1,2)`, "[1:2] Array.pop() 0 or 1 argument"},
		{`[1,2,3].pop("1")`, "[1:2] Array.pop() argument must be an integer"},
		{`[1,2,3].pop(3)`, "[1:2] Array.pop() index out of range"},
		{`[1,2,3].vorn`, "[1:10] chaining operator not supported: ARRAY.vorn"},
		{`[1,2,3].map()`, "[1:2] Array.map() takes exactly 1 argument"},
		{`[1, 2, 3, 4].map(2)`, "[1:19] Array.map() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].map(sqrt, sqrt)`, "[1:2] Array.map() takes exactly 1 argument"},
		{`[1, 2, 3, 4].map()`, "[1:2] Array.map() takes exactly 1 argument"},
		{`[1, 2, 3, 4].map(func() { return true; })`, "[1:19] Array.map() callback must take at least 1 argument"},
		{`[1, 2, 3, 4].filter(2)`, "[1:22] Array.filter() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].filter(func() { return true; })`, "[1:22] Array.filter() callback must take at least 1 argument"},
		{`[1, 2, 3, 4].filter(func(x) { return x > 2; }, func(x) { return x < 2; })`, "[1:2] Array.filter() takes exactly 1 argument"},
		{`[1, 2, 3, 4].filter()`, "[1:2] Array.filter() takes exactly 1 argument"},
		{`[1, 2, 3, 4].reduce(2, 0)`, "[1:22] Array.reduce() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].reduce(func(x, y) { return x + y; }, 0, 0)`, "[1:2] Array.reduce() takes exactly 2 arguments, got 3"},
		{`[1, 2, 3, 4].reduce(func (x) { return x; }, 0)`, "[1:22] Array.reduce() callback must take at least 2 arguments"},
	}

	for _, test := range tests {
		checkChainingExpression(t, test.input, test.expected)
	}
}

func TestObjectChainingExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"a": 1, "b": 2}.keys()`, []string{"a", "b"}},
		{`{"a": 1, "b": 2}.keys(1)`, "[1:2] Object.keys() takes no arguments"},
		{`{"a": 1, "b": 2}.values()`, []string{"1", "2"}},
		{`{"a": 1, "b": 2}.values(1)`, "[1:2] Object.values() takes no arguments"},
		{`{"a": 1, "b": 2}.items()`, []string{"a:1", "b:2"}},
		{`{"a": 1, "b": 2}.items(1)`, "[1:2] Object.items() takes no arguments"},
		{`{}.upper()`, "[1:5] Object has no method upper"},
	}

	for _, test := range tests {
		checkChainingExpression(t, test.input, test.expected)
	}
}

func TestChainingExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{}.upper("2" - "1")`, "[1:15] unknown operator: STRING - STRING"},
	}

	for _, test := range tests {
		checkChainingExpression(t, test.input, test.expected)
	}
}
