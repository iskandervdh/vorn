package evaluator

import (
	"testing"

	"github.com/iskandervdh/vorn/lexer"
	"github.com/iskandervdh/vorn/object"
	"github.com/iskandervdh/vorn/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l, false)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	e := New()

	return e.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %d, want %d", result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)

	if !ok {
		t.Errorf("object is not Float. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %g, want %g", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %T, want %t", result.Value, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)

	if !ok {
		t.Errorf("object is not String. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %q, want %q", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got %T (%+v)", obj, obj)
		return false
	}

	return true
}

func testArrayObject(t *testing.T, obj object.Object, expected []string) bool {
	result, ok := obj.(*object.Array)

	if !ok {
		t.Errorf("object is not Array. got %T (%+v)", obj, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("object has the wrong length. got %d, want %d", len(result.Elements), len(expected))
		return false
	}

	for i := 0; i < len(expected); i++ {
		if result.Elements[i].Inspect() != expected[i] {
			return false
		}
	}

	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	errObj, ok := obj.(*object.Error)

	if !ok {
		t.Errorf("object is not Error. got %T (%+v)", obj, obj)
		return false
	}

	if errObj.Message != expected {
		t.Errorf("wrong error message. expected %q, got %q", expected, errObj.Message)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		// {"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		// {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.0", 5.0},
		{"10.0", 10.0},
		{"-5.0", -5.0},
		{"-10.0", -10.0},
		{"5.5 + 5 + 5 + 5 - 10", 10.5},
		{"2.5 * 2 * 2 * 2 * 2.225", 44.5},
		{"-50.0 + 100 + -50", 0.0},
		{"5.0 * 2 + 10", 20.0},
		{"5 + 2 * 10.0", 25.0},
		{"20 + 2 * -10.0", 0.0},
		{"50 / 2 * 2 + 10.0", 60.0},
		{"2 * (5 + 10.2)", 30.4},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testFloatObject(t, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestBooleanOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", nil},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", nil},
		{"if (1 < 2) { 10 }", nil},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", nil},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let i = 0; while (i < 10) { i = i + 1; }; i;", 10},
		{"let i = 0; while (i < 10) { i = i + 1; if (i == 4) { break; } }; i;", 4},
		{"let x = 0; let i = 0; while (i < 4) { i = i + 1; if (i != 3) { continue; } x = i; }; x;", 3},
		{`while (1 + "") { 1 }`, "[1:11] type mismatch: INTEGER + STRING"},
		{`while (1) { 1 + "" }`, "[1:16] type mismatch: INTEGER + STRING"},
		{`func(x) { while (x) { return x; } }(1)`, 1},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
			continue
		}

		err, ok := test.expected.(string)

		if ok {
			testErrorObject(t, evaluated, err)
			continue
		}

		testNullObject(t, evaluated)
	}
}

func TestFor(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let x = 0; for (let i = 0; i < 10; i = i + 1) { x = x + 1; }; x;", 10},
		{"let x = 0; for (let i = 0; i < 10; i = i + 1) { if (i == 4) { x = i; break; } }; x;", 4},
		{"let x = 0; for (let i = 0; i < 4; i = i + 1) { if (i != 3) { continue; } x = i; }; x;", 3},
		{"let i = 0; for (; i < 10; i = i + 1) { }; i;", 10},
		{"let i = 0; for (; i < 10;) { i = i + 1; }; i;", 10},
		{`for (let i = 0; i < 1 + ""; i = i + 1) { 1 }`, "[1:24] type mismatch: INTEGER + STRING"},
		{`func(x) { for (let i = 0; i < x; i = i + 1) { return x; } }(1)`, 1},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
			continue
		}

		err, ok := test.expected.(string)

		if ok {
			testErrorObject(t, evaluated, err)
			continue
		}

		testNullObject(t, evaluated)
	}
}

func TestExclamationOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!first([])", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestBitwiseNotOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"~0", -1},
		{"~1", -2},
		{"~2", -3},
		{"~-1", 0},
		{"~-2", 1},
		{"~-3", 2},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestIntegerReturningInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1 + 1", 2},
		{"1 - 1", 0},
		{"1 * 1", 1},
		{`1 % 1`, 0.0},
		{"1 & 1", 1},
		{"1 & 2", 0},
		{"1 | 1", 1},
		{"1 | 2", 3},
		{"1 ^ 1", 0},
		{"1 ^ 2", 3},
		{"1 << 1", 2},
		{"1 << 2", 4},
		{"4 >> 1", 2},
		{"4 >> 2", 1},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestFloatReturningInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1.0 + 1.0", 2.0},
		{"1.0 - 1.0", 0.0},
		{"1.0 * 1.0", 1.0},
		{"1.0 / 1.0", 1.0},
		{`1 / 1`, 1.0},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testFloatObject(t, evaluated, test.expected)
	}
}

func TestBooleanReturningInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true != true", false},
		{"true == false", false},
		{"true != false", true},
		{"false == false", true},
		{"false != false", false},
		{"false == true", false},
		{"false != true", true},
		{`1.0 == 1.0`, true},
		{`1.0 != 1.0`, false},
		{`1.0 == 2.0`, false},
		{`1.0 != 2.0`, true},
		{`1 < 2`, true},
		{`1 > 2`, false},
		{`1 <= 2`, true},
		{`1 >= 2`, false},
		{`1 <= 1`, true},
		{`1 >= 1`, true},
		{`1 < 1`, false},
		{`1 > 1`, false},
		{`1.0 < 2.0`, true},
		{`1.0 > 2.0`, false},
		{`1.0 <= 2.0`, true},
		{`1.0 >= 2.0`, false},
		{`1.0 <= 1.0`, true},
		{`1.0 >= 1.0`, true},
		{`1.0 < 1.0`, false},
		{`1.0 > 1.0`, false},
		{`"hello" == "hello"`, true},
		{`"hello" != "hello"`, false},
		{`"hello" == "world"`, false},
		{`"hello" != "world"`, true},
		{`1 || 0`, true},
		{`1 || 1`, true},
		{`0 || 1`, true},
		{`0 || 0`, false},
		{`1 && 0`, false},
		{`1 && 1`, true},
		{`0 && 1`, false},
		{`0 && 0`, false},
		{`true || false`, true},
		{`true && false`, false},
		{`false || true`, true},
		{`false && true`, false},
		{`true || true`, true},
		{`true && true`, true},
		{`false || false`, false},
		{`null || false`, false},
		{`null && false`, false},
		{`null || true`, true},
		{`null && true`, false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
	return 10;
} else {
	return 1;
}`,
			10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"[1:4] type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"[1:4] type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"[1:1] unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"[1:7] unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"[1:10] unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"[1:21] unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return 1;
}`,
			"[4:16] unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"[1:2] identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"[1:10] unknown operator: STRING - STRING",
		},
		{
			`{"name": "Vorn"}[func(x) { x }];`,
			"[1:19] unusable as object key: FUNCTION",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error object returned. got %T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != test.expectedMessage {
			t.Errorf("wrong error message. expected %q, got %q", test.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"
	evaluated := testEval(input)
	function, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. got %T (%+v)", evaluated, evaluated)
	}

	if len(function.Arguments) != 1 {
		t.Fatalf("function has wrong arguments. Arguments=%+v", function.Arguments)
	}

	if function.Arguments[0].String() != "x" {
		t.Fatalf("argument is not 'x'. got %q", function.Arguments[0])
	}

	expectedBody := "{\n  (x + 2)\n}"

	if function.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, function.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { return x * 2; }; double(5);", 10},
		{"let add = func(x, y) { return x + y; }; add(5, 5);", 10},
		{"let add = func(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { return x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = func(x) {
	return func(y) { return x + y; };
};

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got %q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got %T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got %q", str.Value)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("object is not Array. got %T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got %d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let arr = [1, 2, 3]; arr[2];",
			3,
		},
		{
			"let arr = [1, 2, 3]; arr[0] + arr[1] + arr[2];",
			6,
		},
		{
			"let arr = [1, 2, 3]; let i = arr[0]; arr[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3, 4, 5, 6][-1]",
			6,
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{
"one": 10 - 9,
two: 1 + 1,
"thr" + "ee": -1 * -3,
4: 4,
true: 5,
false: 6
}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)

	if !ok {
		t.Fatalf("Eval didn't return Hash. got %T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(object.NewString(nil, "one")).HashKey():   1,
		(object.NewString(nil, "two")).HashKey():   2,
		(object.NewString(nil, "three")).HashKey(): 3,
		(object.NewInteger(nil, 4)).HashKey():      4,
		object.TRUE.HashKey():                      5,
		object.FALSE.HashKey():                     6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got %d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]

		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestVariableReassignment(t *testing.T) {
	input := `let x = 1;
x = 4;
x;`

	testIntegerObject(t, testEval(input), 4)
}

func TestIsError(t *testing.T) {
	input := `5 + "a"`
	evaluated := testEval(input)
	if evaluated.Type() != object.ERROR_OBJ {
		t.Errorf("no error object returned. got %T (%+v)", evaluated, evaluated)
	}

	if isError(nil) {
		t.Errorf("nil is error")
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"first([])", false},
		{"null", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		if isTruthy(evaluated) != test.expected {
			t.Errorf("wrong truth value. got %t for input %s", isTruthy(evaluated), test.input)
		}
	}
}

func TestReassignmentExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let x = 1; x = 2; x;", 2},
		{"let x = 1; x += 4; x;", 5},
		{"let x = 1; x -= 4; x;", -3},
		{"let x = 1; x *= 4; x;", 4},
		{"let x = 1; x /= 4; x;", 0.25},
		{"let x = 1; x %= 4; x;", 1},
		{"let x = 1; x &= 4; x;", 0},
		{"let x = 1; x |= 4; x;", 5},
		{"let x = 1; x ^= 4; x;", 5},
		{"let x = 1; x <<= 4; x;", 16},
		{"let x = 1; x >>= 4; x;", 0},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
			continue
		}

		float, ok := test.expected.(float64)

		if ok {
			testFloatObject(t, evaluated, float)
			continue
		}

		err, ok := test.expected.(string)

		if ok {
			testErrorObject(t, evaluated, err)
			continue
		}

		testNullObject(t, evaluated)
	}
}

func TestIncrementDecrementExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let x = 1; x++; x;", 2},
		{"let x = 1; x--; x;", 0},
		{"let x = 1; ++x; x;", 2},
		{"let x = 1; --x; x;", 0},
		{"let x = 1.5; x++; x;", 2.5},
		{"let x = 1.5; x--; x;", 0.5},
		{"let x = 1.5; ++x; x;", 2.5},
		{"let x = 1.5; --x; x;", 0.5},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
			continue
		}

		float, ok := test.expected.(float64)

		if ok {
			testFloatObject(t, evaluated, float)
			continue
		}

		testNullObject(t, evaluated)
	}
}
