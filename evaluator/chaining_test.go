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
	case object.NULL_OBJ:
		if expected != "null" {
			t.Errorf("Expected %s, got NULL", expected)
		}
	default:
		t.Errorf("Expected STRING_OBJ, INTEGER_OBJ, ERROR_OBJ, ARRAY_OBJ, BOOLEAN_OBJ or NULL_OBJ, got %s", evaluated.Type())
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
		{`"hello".split(1)`, "[1:2] argument to `String.split()` must be STRING, got INTEGER"},
		{`"hello".split("e", "l")`, "[1:2] String.split() takes at most 1 argument, got 2"},
		{`"hello world".contains("world", 6)`, "[1:2] String.contains() takes exactly 1 argument"},
		{`"hello world".contains(6)`, "[1:2] argument to `String.contains()` must be STRING, got INTEGER"},
		{`"hello world".replace("world", "you", "me")`, "[1:2] String.replace() takes exactly 2 arguments"},
		{`"hello world".replace(1, 2)`, "[1:2] first argument to `String.replace()` must be STRING, got INTEGER"},
		{`"hello world".replace("world", 2)`, "[1:2] second argument to `String.replace()` must be STRING, got INTEGER"},
		{`"hello".trim(1)`, "[1:2] String.trim() takes no arguments"},
		{`"hello".trimStart(1)`, "[1:2] String.trimStart() takes no arguments"},
		{`"hello".trimEnd(1)`, "[1:2] String.trimEnd() takes no arguments"},
		{`"hello".repeat()`, "[1:2] String.repeat() takes exactly 1 argument"},
		{`"hello".repeat("1")`, "[1:2] argument to `String.repeat()` must be INTEGER, got STRING"},
		{`"hello".reverse(1)`, "[1:2] String.reverse() takes no arguments"},
		{`"hello".slice()`, "[1:2] String.slice() takes 1 or 2 arguments"},
		{`"hello".slice(1, 2, 3)`, "[1:2] String.slice() takes 1 or 2 arguments"},
		{`"hello".slice("1", 2)`, "[1:2] first argument to `String.slice()` must be INTEGER, got STRING"},
		{`"hello".slice(1, "2")`, "[1:2] second argument to `String.slice()` must be INTEGER, got STRING"},
		{`"hello".slice(10, 10)`, "[1:2] first argument to `String.slice()` out of range"},
		{`"hello".slice(0, 10)`, "[1:2] second argument to `String.slice()` out of range"},
		{`"hello".startsWith(1)`, "[1:2] argument to `String.startsWith()` must be STRING, got INTEGER"},
		{`"hello".startsWith("1", 2)`, "[1:2] String.startsWith() takes exactly 1 argument"},
		{`"hello".endsWith(1)`, "[1:2] argument to `String.endsWith()` must be STRING, got INTEGER"},
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
		{`[1,2,3].prepend(0).length()`, 4},
		{`let a = [1,2,3]; a.prepend(0); a.length()`, 4},
		{`[1,2,3].shift()`, 1},
		{`let a = [1,2,3]; a.shift(); a.shift(); a`, []string{"3"}},
		{`[1,2,3].pop()`, 3},
		{`let a = [1,2,3]; a.pop(); a.pop(); a`, []string{"1"}},
		{`let a = [1,2,3]; a.pop(1)`, 2},
		{`[1,2,3].concat([4,5,6])`, []string{"1", "2", "3", "4", "5", "6"}},
		{`[1,2,3].concat([4,5,6], [7,8,9])`, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}},

		{`func timesTwo(x) { x * 2 }; [1, 2, 3, 4].map(timesTwo)`, []string{"2", "4", "6", "8"}},
		{`[1, 2, 3, 4].map(sqrt)`, []string{"1", "1.4142135623730951", "1.7320508075688772", "2"}},
		{`[1, 2, 3, 4].map(func(x, i) { return x + i; })`, []string{"1", "3", "5", "7"}},
		{`[1, 2, 3, 4].filter(func(x) { return x > 2; })`, []string{"3", "4"}},
		{`[1, 2, 3, 4].filter(func(x) { return 10; })`, []string{"1", "2", "3", "4"}},
		{`[1, 2, 3, 4].reduce(func(x, y) { return x + y; }, 0)`, 10},
		{`[1, 2, 3, 4].reduce(func(x, y, i) { return x + y + i; }, 0)`, 16},

		{`[1, 2, 3, 4].contains(2)`, true},
		{`[1, 2, 3, 4].contains(5)`, false},
		{`[1, 2, 3, 4].indexOf(2)`, 1},
		{`[1, 2, 3, 4].indexOf(5)`, -1},
		{`[1, 2, 3, 4].find(func(x) { return x > 2; })`, 3},
		{`[1, 2, 3, 4].find(func(x) { return x > 5; })`, "null"},
		{`[1, 2, 3, 4].find(func(x) { return 1; })`, 1},
		{`["a", "b", "c", "d"].find(func(x) { return x == "c"; })`, "c"},

		{`["a", "b", "c", "d"].join()`, "a,b,c,d"},
		{`["a", "b", "c", "d"].join("")`, "abcd"},
		{`[1, 2, 3, 4].reverse()`, []string{"4", "3", "2", "1"}},
		{`[1, 2, 3, 4].reverse().reverse()`, []string{"1", "2", "3", "4"}},
		{`[1, 2, 3, 4].slice(1)`, []string{"2", "3", "4"}},
		{`[1, 2, 3, 4].slice(1, 3)`, []string{"2", "3"}},
		{`[1, 2, 3, 4].slice(1, 1)`, []string{}},
		{`[1, 2, 3, 4].slice(1, -1)`, []string{"2", "3"}},

		{`[3,6,8,3,1,2,4,6,3].sort()`, []string{"1", "2", "3", "3", "3", "4", "6", "6", "8"}},
		{`[3,6,8,3,1,2,4,6,3].sort(false)`, []string{"1", "2", "3", "3", "3", "4", "6", "6", "8"}},
		{`[3,6,8,3,1,2,4,6,3].sort(true)`, []string{"8", "6", "6", "4", "3", "3", "3", "2", "1"}},
		{`[3,6,8,3,1,2,4,6,3].sort(func(a, b) { return a - b; })`, []string{"1", "2", "3", "3", "3", "4", "6", "6", "8"}},
		{`[3,6,8,3,1,2,4,6,3].sort(func(a, b) { return b - a; })`, []string{"8", "6", "6", "4", "3", "3", "3", "2", "1"}},
		{`[].sort()`, []string{}},

		{`[1,5,2,3].any(func(x) {return x > 4;})`, true},
		{`[1,5,2,3].any(func(x) {return x > 10;})`, false},
		{`[1,5,2,3].any(func(x) { return 1; })`, true},
		{`[1,2,3,4].every(func(x) {return x != 0;})`, true},
		{`[1,2,0,4].every(func(x) {return x == 0;})`, false},
		{`[1,2,0,4].every(func(x) { return 0; })`, false},

		// Errors
		{`[1,2,3].vorn`, "[1:10] chaining operator not supported: ARRAY.vorn"},
		{`[1,2,3].upper()`, "[1:10] Array has no method upper"},

		{`[1,2,3].length(1)`, "[1:2] Array.length() takes no arguments"},
		{`[1,2,3].append()`, "[1:2] Array.append() takes exactly 1 argument"},
		{`[1,2,3].prepend()`, "[1:2] Array.prepend() takes exactly 1 argument"},
		{`[].shift()`, "[1:2] Array.shift() called on empty array"},
		{`[1,2,3].shift("1")`, "[1:2] Array.shift() takes no arguments"},
		{`[].pop()`, "[1:2] Array.pop() called on empty array"},
		{`[1,2,3].pop(1,2)`, "[1:2] Array.pop() 0 or 1 argument"},
		{`[1,2,3].pop("1")`, "[1:2] Array.pop() argument must be an integer"},
		{`[1,2,3].pop(3)`, "[1:2] Array.pop() index out of range"},
		{`[1,2,3].concat()`, "[1:2] Array.concat() takes at least 1 argument"},
		{`[1,2,3].concat(1)`, "[1:2] argument to `Array.concat()` must be ARRAY, got INTEGER"},

		{`[1,2,3].map()`, "[1:2] Array.map() takes exactly 1 argument"},
		{`[1, 2, 3, 4].map(2)`, "[1:19] Array.map() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].map(sqrt, sqrt)`, "[1:2] Array.map() takes exactly 1 argument"},
		{`[1, 2, 3, 4].map()`, "[1:2] Array.map() takes exactly 1 argument"},
		{`[1, 2, 3, 4].map(func() { return true; })`, "[1:19] Array.map() callback must take at least 1 argument"},
		{`[1, 2, 3, 4].map(func(x) { if (x == 2) { return x + ""; } return x; })`, "[1:52] type mismatch: INTEGER + STRING"},

		{`[1, 2, 3, 4].filter(2)`, "[1:22] Array.filter() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].filter(func() { return true; })`, "[1:22] Array.filter() callback must take at least 1 argument"},
		{`[1, 2, 3, 4].filter(func(x) { return x > 2; }, func(x) { return x < 2; })`, "[1:2] Array.filter() takes exactly 1 argument"},
		{`[1, 2, 3, 4].filter()`, "[1:2] Array.filter() takes exactly 1 argument"},
		{`[1, 2, 3, 4].filter(func(x) { if (x == 2) { return x + ""; } return x > 1; })`, "[1:55] type mismatch: INTEGER + STRING"},

		{`[1, 2, 3, 4].reduce(2, 0)`, "[1:22] Array.reduce() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].reduce(func(x, y) { return x + y; }, 0, 0)`, "[1:2] Array.reduce() takes exactly 2 arguments, got 3"},
		{`[1, 2, 3, 4].reduce(func (x) { return x; }, 0)`, "[1:22] Array.reduce() callback must take at least 2 arguments"},
		{`[1, 2, 3, 4].reduce(func(x, y) { if (y == 2) { return x + y + ""; } return x + y; }, 0)`, "[1:62] type mismatch: INTEGER + STRING"},

		{`[1, 2, 3, 4].contains()`, "[1:2] Array.contains() takes exactly 1 argument"},
		{`[1, 2, 3, 4].contains(2, 3)`, "[1:2] Array.contains() takes exactly 1 argument"},
		{`[1, 2, 3, 4].indexOf()`, "[1:2] Array.indexOf() takes exactly 1 argument"},
		{`[1, 2, 3, 4].indexOf(2, 3)`, "[1:2] Array.indexOf() takes exactly 1 argument"},
		{`[1, 2, 3, 4].find(1, 2)`, "[1:2] Array.find() takes exactly 1 argument"},
		{`[1, 2, 3, 4].find(1)`, "[1:20] Array.find() callback must be a function, got INTEGER"},
		{`[1, 2, 3, 4].find(func(){})`, "[1:20] Array.find() callback must take at least 1 argument"},
		{`[1, 2, 3, 4].find(func(x){ return x + ""; })`, "[1:38] type mismatch: INTEGER + STRING"},

		{`[1,2,3,4].join(1)`, "[1:2] argument to `Array.join()` must be STRING, got INTEGER"},
		{`[1,2,3,4].join("", " ")`, "[1:2] Array.join() takes at most 1 argument, got 2"},
		{`[1,2,3,4].reverse(1)`, "[1:2] Array.reverse() takes no arguments"},
		{`[1,2,3,4].slice()`, "[1:2] Array.slice() takes 1 or 2 arguments"},
		{`[1,2,3,4].slice("")`, "[1:2] first argument to `Array.slice()` must be INTEGER, got STRING"},
		{`[1,2,3,4].slice(1, "2")`, "[1:2] second argument to `Array.slice()` must be INTEGER, got STRING"},
		{`[1,2,4,5].slice(-1)`, "[1:2] first argument to `Array.slice()` out of range"},
		{`[1,2,4,5].slice(0, 20)`, "[1:2] second argument to `Array.slice()` out of range"},

		{`[1,2,3,4].sort(1)`, "[1:2] argument to `Array.sort()` must be BOOLEAN, FUNCTION or BUILTIN, got INTEGER"},
		{`[1,2,3,4].sort(func(a, b) { return a - b; }, func(a, b) { return b - a; })`, "[1:2] Array.sort() takes at most 1 argument, got 2"},
		{`[1,2,3,4].sort(func(){})`, "[1:17] Array.sort() callback must take at least 2 arguments"},
		{`[1,2,3,4].sort(func(a, b) { return b + ""; })`, "[1:39] type mismatch: INTEGER + STRING"},

		{`[1,2,3,4].any()`, "[1:2] Array.any() takes exactly 1 argument"},
		{`[1,2,3,4].any(1)`, "[1:16] Array.any() callback must be a function, got INTEGER"},
		{`[1,2,3,4].any(func(){})`, "[1:16] Array.any() callback must take at least 1 argument"},
		{`[1,2,3,4].any(func(x){return x + "";})`, "[1:33] type mismatch: INTEGER + STRING"},
		{`[1,2,3,4].every()`, "[1:2] Array.every() takes exactly 1 argument"},
		{`[1,2,3,4].every(1)`, "[1:18] Array.every() callback must be a function, got INTEGER"},
		{`[1,2,3,4].every(func(){})`, "[1:18] Array.every() callback must take at least 1 argument"},
		{`[1,2,3,4].every(func(x){return x + "";})`, "[1:35] type mismatch: INTEGER + STRING"},
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
